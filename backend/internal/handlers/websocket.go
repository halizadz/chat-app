package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/halizadz/chat-app-backend/internal/middleware"
	"github.com/halizadz/chat-app-backend/internal/models"
	"github.com/halizadz/chat-app-backend/internal/repository"
	"github.com/halizadz/chat-app-backend/internal/utils"
	ws "github.com/halizadz/chat-app-backend/internal/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // In production, check origin properly
	},
}

type WebSocketHandler struct {
	hub         *ws.Hub
	roomRepo    *repository.RoomRepository
	messageRepo *repository.MessageRepository
	jwtSecret   string
}

func NewWebSocketHandler(hub *ws.Hub, roomRepo *repository.RoomRepository, messageRepo *repository.MessageRepository, jwtSecret string) *WebSocketHandler {
	return &WebSocketHandler{
		hub:         hub,
		roomRepo:    roomRepo,
		messageRepo: messageRepo,
		jwtSecret:   jwtSecret,
	}
}

func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	log.Printf("WebSocket connection attempt from %s", r.RemoteAddr)

	// Try to get token from query parameter
	tokenString := r.URL.Query().Get("token")

	var claims *utils.Claims
	var ok bool

	if tokenString != "" {
		log.Printf("Validating token from query parameter")
		var err error
		claims, err = utils.ValidateToken(tokenString, h.jwtSecret)
		if err != nil {
			log.Printf("Token validation failed: %v", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		log.Printf("Token validated for user: %s", claims.Username)
	} else {
		// Try to get from context (middleware)
		claims, ok = middleware.GetUserFromContext(r.Context())
		if !ok {
			log.Printf("No token found")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	vars := mux.Vars(r)
	roomID, err := uuid.Parse(vars["roomId"])
	if err != nil {
		log.Printf("Invalid room ID: %v", err)
		http.Error(w, "Invalid room ID", http.StatusBadRequest)
		return
	}

	log.Printf("Checking room membership for user %s in room %s", claims.Username, roomID)
	// Check if user is member of room
	isMember, err := h.roomRepo.IsMember(roomID, claims.UserID)
	if err != nil {
		log.Printf("Error checking membership: %v", err)
		http.Error(w, "Error checking membership", http.StatusInternalServerError)
		return
	}

	if !isMember {
		log.Printf("User %s is not a member of room %s", claims.Username, roomID)
		http.Error(w, "Not a member of this room", http.StatusForbidden)
		return
	}

	log.Printf("Upgrading connection to WebSocket for user %s", claims.Username)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}

	log.Printf("WebSocket connected: user=%s room=%s", claims.Username, roomID)
	client := ws.NewClient(h.hub, conn, claims.UserID, claims.Username)

	h.hub.Register <- client
	h.hub.JoinRoom(client, roomID)

	// Start goroutines for reading and writing
	go client.WritePump()
	go h.readPump(client, roomID)
}

func (h *WebSocketHandler) readPump(client *ws.Client, roomID uuid.UUID) {
	defer func() {
		h.hub.Unregister <- client
		client.Conn.Close()
	}()

	client.Conn.SetReadDeadline(time.Now().Add(ws.PongWait))
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(ws.PongWait))
		return nil
	})

	client.Conn.SetReadLimit(ws.MaxMessageSize)

	for {
		_, messageBytes, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var msg ws.Message
		if err := json.Unmarshal(messageBytes, &msg); err != nil {
			log.Printf("error unmarshaling message: %v", err)
			continue
		}

		msg.SenderID = client.ID
		msg.Username = client.Username
		msg.RoomID = roomID
		msg.Timestamp = time.Now()

		// Handle different message types
		switch msg.Type {
		case "message", "file":
			// Validate message content
			if msg.Type == "message" && len(msg.Content) == 0 {
				log.Printf("Empty message rejected from user %s", client.Username)
				continue
			}
			if len(msg.Content) > 10000 { // Max 10KB content
				log.Printf("Message too long rejected from user %s", client.Username)
				continue
			}

			// Save message to database
			dbMessage := &models.Message{
				RoomID:   msg.RoomID,
				SenderID: msg.SenderID,
				Content:  msg.Content,
				Type:     msg.Type,
			}

			if msg.Type == "file" {
				dbMessage.FileURL = &msg.FileURL
				dbMessage.FileName = &msg.FileName
				dbMessage.FileSize = &msg.FileSize
			}

			if err := h.messageRepo.Create(dbMessage); err != nil {
				log.Printf("error saving message: %v", err)
				continue
			}

			msg.Timestamp = dbMessage.CreatedAt

			// Broadcast to all clients in room
			h.hub.Broadcast <- &msg

			// Mark message as read for sender (they sent it, so they've seen it)
			h.messageRepo.MarkAsRead(dbMessage.ID, client.ID)

		case "typing":
			// Handle typing indicator
			typingIndicator := &ws.TypingIndicator{
				RoomID:   roomID,
				UserID:   client.ID,
				Username: client.Username,
				IsTyping: true,
			}

			// Check if content is "stop" to indicate stop typing
			if msg.Content == "stop" {
				typingIndicator.IsTyping = false
			}

			h.hub.Typing <- typingIndicator
		}
	}
}
