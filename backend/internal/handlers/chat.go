package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/google/uuid"
    "github.com/gorilla/mux"
    "github.com/halizadz/chat-app-backend/internal/middleware"
    "github.com/halizadz/chat-app-backend/internal/models"
    "github.com/halizadz/chat-app-backend/internal/repository"
)

type ChatHandler struct {
    roomRepo    *repository.RoomRepository
    messageRepo *repository.MessageRepository
    userRepo    *repository.UserRepository
}

func NewChatHandler(roomRepo *repository.RoomRepository, messageRepo *repository.MessageRepository, userRepo *repository.UserRepository) *ChatHandler {
    return &ChatHandler{
        roomRepo:    roomRepo,
        messageRepo: messageRepo,
        userRepo:    userRepo,
    }
}

type CreateRoomRequest struct {
    Name        string  `json:"name"`
    Description *string `json:"description"`
    Type        string  `json:"type"` // private or group
}

type CreatePrivateRoomRequest struct {
    UserID uuid.UUID `json:"user_id"`
}

type SendMessageRequest struct {
    Content string `json:"content"`
    Type    string `json:"type"` // text, file, image
}

func (h *ChatHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
    claims, ok := middleware.GetUserFromContext(r.Context())
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    var req CreateRoomRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if req.Name == "" || req.Type == "" {
        http.Error(w, "Name and type are required", http.StatusBadRequest)
        return
    }

    if req.Type != "private" && req.Type != "group" {
        http.Error(w, "Type must be 'private' or 'group'", http.StatusBadRequest)
        return
    }

    room := &models.Room{
        Name:        req.Name,
        Description: req.Description,
        Type:        req.Type,
        CreatedBy:   claims.UserID,
    }

    if err := h.roomRepo.Create(room); err != nil {
        http.Error(w, "Error creating room: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Add creator as admin
    if err := h.roomRepo.AddMember(room.ID, claims.UserID, "admin"); err != nil {
        http.Error(w, "Error adding member: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(room)
}

func (h *ChatHandler) CreateOrGetPrivateRoom(w http.ResponseWriter, r *http.Request) {
    claims, ok := middleware.GetUserFromContext(r.Context())
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    var req CreatePrivateRoomRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Check if user exists
    _, err := h.userRepo.FindByID(req.UserID)
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    room, err := h.roomRepo.FindOrCreatePrivateRoom(claims.UserID, req.UserID)
    if err != nil {
        http.Error(w, "Error creating private room: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(room)
}

func (h *ChatHandler) GetUserRooms(w http.ResponseWriter, r *http.Request) {
    claims, ok := middleware.GetUserFromContext(r.Context())
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    rooms, err := h.roomRepo.GetUserRooms(claims.UserID)
    if err != nil {
        http.Error(w, "Error fetching rooms: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(rooms)
}

func (h *ChatHandler) GetRoom(w http.ResponseWriter, r *http.Request) {
    claims, ok := middleware.GetUserFromContext(r.Context())
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    vars := mux.Vars(r)
    roomID, err := uuid.Parse(vars["roomId"])
    if err != nil {
        http.Error(w, "Invalid room ID", http.StatusBadRequest)
        return
    }

    // Check if user is member
    isMember, err := h.roomRepo.IsMember(roomID, claims.UserID)
    if err != nil {
        http.Error(w, "Error checking membership", http.StatusInternalServerError)
        return
    }

    if !isMember {
        http.Error(w, "Not a member of this room", http.StatusForbidden)
        return
    }

    room, err := h.roomRepo.FindByID(roomID)
    if err != nil {
        http.Error(w, "Room not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(room)
}

func (h *ChatHandler) GetRoomMessages(w http.ResponseWriter, r *http.Request) {
    claims, ok := middleware.GetUserFromContext(r.Context())
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    vars := mux.Vars(r)
    roomID, err := uuid.Parse(vars["roomId"])
    if err != nil {
        http.Error(w, "Invalid room ID", http.StatusBadRequest)
        return
    }

    // Check if user is member
    isMember, err := h.roomRepo.IsMember(roomID, claims.UserID)
    if err != nil {
        http.Error(w, "Error checking membership", http.StatusInternalServerError)
        return
    }

    if !isMember {
        http.Error(w, "Not a member of this room", http.StatusForbidden)
        return
    }

    // Get pagination params
    limitStr := r.URL.Query().Get("limit")
    offsetStr := r.URL.Query().Get("offset")

    limit := 50
    offset := 0

    if limitStr != "" {
        if l, err := strconv.Atoi(limitStr); err == nil {
            limit = l
        }
    }

    if offsetStr != "" {
        if o, err := strconv.Atoi(offsetStr); err == nil {
            offset = o
        }
    }

    messages, err := h.messageRepo.GetByRoomID(roomID, limit, offset)
    if err != nil {
        http.Error(w, "Error fetching messages: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(messages)
}

func (h *ChatHandler) GetRoomMembers(w http.ResponseWriter, r *http.Request) {
    claims, ok := middleware.GetUserFromContext(r.Context())
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    vars := mux.Vars(r)
    roomID, err := uuid.Parse(vars["roomId"])
    if err != nil {
        http.Error(w, "Invalid room ID", http.StatusBadRequest)
        return
    }

    // Check if user is member
    isMember, err := h.roomRepo.IsMember(roomID, claims.UserID)
    if err != nil {
        http.Error(w, "Error checking membership", http.StatusInternalServerError)
        return
    }

    if !isMember {
        http.Error(w, "Not a member of this room", http.StatusForbidden)
        return
    }

    members, err := h.roomRepo.GetMembers(roomID)
    if err != nil {
        http.Error(w, "Error fetching members: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(members)
}

func (h *ChatHandler) AddRoomMember(w http.ResponseWriter, r *http.Request) {
    claims, ok := middleware.GetUserFromContext(r.Context())
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    vars := mux.Vars(r)
    roomID, err := uuid.Parse(vars["roomId"])
    if err != nil {
        http.Error(w, "Invalid room ID", http.StatusBadRequest)
        return
    }

    var req struct {
        UserID uuid.UUID `json:"user_id"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Check if requester is member
    isMember, err := h.roomRepo.IsMember(roomID, claims.UserID)
    if err != nil {
        http.Error(w, "Error checking membership", http.StatusInternalServerError)
        return
    }

    if !isMember {
        http.Error(w, "Not a member of this room", http.StatusForbidden)
        return
    }

    // Add new member
    if err := h.roomRepo.AddMember(roomID, req.UserID, "member"); err != nil {
        http.Error(w, "Error adding member: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Member added successfully"})
}

func (h *ChatHandler) LeaveRoom(w http.ResponseWriter, r *http.Request) {
    claims, ok := middleware.GetUserFromContext(r.Context())
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    vars := mux.Vars(r)
    roomID, err := uuid.Parse(vars["roomId"])
    if err != nil {
        http.Error(w, "Invalid room ID", http.StatusBadRequest)
        return
    }

    if err := h.roomRepo.RemoveMember(roomID, claims.UserID); err != nil {
        http.Error(w, "Error leaving room: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Left room successfully"})
}