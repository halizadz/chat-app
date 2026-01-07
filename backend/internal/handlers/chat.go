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

	// Check if user exists and get username
	otherUser, err := h.userRepo.FindByID(req.UserID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	room, err := h.roomRepo.FindOrCreatePrivateRoom(claims.UserID, req.UserID, otherUser.Username)
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

// MarkRoomAsRead marks all messages in a room as read
func (h *ChatHandler) MarkRoomAsRead(w http.ResponseWriter, r *http.Request) {
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

	// Mark all messages as read
	if err := h.messageRepo.MarkRoomMessagesAsRead(roomID, claims.UserID); err != nil {
		http.Error(w, "Error marking messages as read: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Messages marked as read"})
}

// SearchMessages searches messages in a room
func (h *ChatHandler) SearchMessages(w http.ResponseWriter, r *http.Request) {
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

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Search query is required", http.StatusBadRequest)
		return
	}

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

	messages, err := h.messageRepo.SearchMessages(roomID, query, limit, offset)
	if err != nil {
		http.Error(w, "Error searching messages: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// UpdateMessage updates a message
func (h *ChatHandler) UpdateMessage(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	messageID, err := uuid.Parse(vars["messageId"])
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	// Get message
	message, err := h.messageRepo.FindByID(messageID)
	if err != nil {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	// Check if user is the sender
	if message.SenderID != claims.UserID {
		http.Error(w, "You can only edit your own messages", http.StatusForbidden)
		return
	}

	var req struct {
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Content == "" {
		http.Error(w, "Content is required", http.StatusBadRequest)
		return
	}

	// Update message
	if err := h.messageRepo.Update(messageID, req.Content); err != nil {
		http.Error(w, "Error updating message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get updated message
	updatedMessage, _ := h.messageRepo.FindByID(messageID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedMessage)
}

// DeleteMessage deletes a message
func (h *ChatHandler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	messageID, err := uuid.Parse(vars["messageId"])
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	// Get message
	message, err := h.messageRepo.FindByID(messageID)
	if err != nil {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	// Check if user is the sender or admin
	isSender := message.SenderID == claims.UserID
	isAdmin := false

	if !isSender {
		// Check if user is admin of the room
		members, err := h.roomRepo.GetMembers(message.RoomID)
		if err == nil {
			for _, member := range members {
				if member.ID == claims.UserID {
					// Check role in room_members
					// For simplicity, we'll check if user is creator
					room, _ := h.roomRepo.FindByID(message.RoomID)
					isAdmin = room.CreatedBy == claims.UserID
					break
				}
			}
		}
	}

	if !isSender && !isAdmin {
		http.Error(w, "You can only delete your own messages", http.StatusForbidden)
		return
	}

	// Delete message
	if err := h.messageRepo.Delete(messageID); err != nil {
		http.Error(w, "Error deleting message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Message deleted successfully"})
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

	// Prevent user from adding themselves
	if req.UserID == claims.UserID {
		http.Error(w, "Cannot add yourself as a member", http.StatusBadRequest)
		return
	}

	// Check if user is already a member
	isAlreadyMember, err := h.roomRepo.IsMember(roomID, req.UserID)
	if err != nil {
		http.Error(w, "Error checking membership", http.StatusInternalServerError)
		return
	}
	if isAlreadyMember {
		http.Error(w, "User is already a member of this room", http.StatusConflict)
		return
	}

	// Check if room is private and already has 2 members
	room, err := h.roomRepo.FindByID(roomID)
	if err != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	if room.Type == "private" {
		members, err := h.roomRepo.GetMembers(roomID)
		if err != nil {
			http.Error(w, "Error fetching members", http.StatusInternalServerError)
			return
		}
		if len(members) >= 2 {
			http.Error(w, "Private room can only have 2 members", http.StatusBadRequest)
			return
		}
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

// UpdateRoom updates room information
func (h *ChatHandler) UpdateRoom(w http.ResponseWriter, r *http.Request) {
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

	// Get room
	room, err := h.roomRepo.FindByID(roomID)
	if err != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	// Check if user is creator or admin
	if room.CreatedBy != claims.UserID {
		http.Error(w, "Only room creator can update room", http.StatusForbidden)
		return
	}

	var req struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name != nil {
		room.Name = *req.Name
	}
	if req.Description != nil {
		room.Description = req.Description
	}

	if err := h.roomRepo.Update(room); err != nil {
		http.Error(w, "Error updating room: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(room)
}

// DeleteRoom deletes a room
func (h *ChatHandler) DeleteRoom(w http.ResponseWriter, r *http.Request) {
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

	// Get room
	room, err := h.roomRepo.FindByID(roomID)
	if err != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	// Check if user is creator
	if room.CreatedBy != claims.UserID {
		http.Error(w, "Only room creator can delete room", http.StatusForbidden)
		return
	}

	if err := h.roomRepo.Delete(roomID); err != nil {
		http.Error(w, "Error deleting room: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Room deleted successfully"})
}

// RemoveRoomMember removes a member from room
func (h *ChatHandler) RemoveRoomMember(w http.ResponseWriter, r *http.Request) {
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

	userID, err := uuid.Parse(vars["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Check if requester is creator or admin
	room, err := h.roomRepo.FindByID(roomID)
	if err != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	isCreator := room.CreatedBy == claims.UserID
	if !isCreator {
		http.Error(w, "Only room creator can remove members", http.StatusForbidden)
		return
	}

	// Cannot remove creator
	if userID == room.CreatedBy {
		http.Error(w, "Cannot remove room creator", http.StatusBadRequest)
		return
	}

	if err := h.roomRepo.RemoveMember(roomID, userID); err != nil {
		http.Error(w, "Error removing member: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Member removed successfully"})
}
