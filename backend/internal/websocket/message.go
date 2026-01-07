package websocket

import (
    "time"
    "github.com/google/uuid"
)

// Constants for WebSocket configuration
const (
    WriteWait      = 10 * time.Second
    PongWait       = 60 * time.Second
    PingPeriod     = (PongWait * 9) / 10
    MaxMessageSize = 512 * 1024 // 512 KB
)

// Message represents a chat message
type Message struct {
    Type      string    `json:"type"` // message, typing, join, leave, file
    RoomID    uuid.UUID `json:"room_id"`
    SenderID  uuid.UUID `json:"sender_id"`
    Username  string    `json:"username"`
    Content   string    `json:"content"`
    FileURL   string    `json:"file_url,omitempty"`
    FileName  string    `json:"file_name,omitempty"`
    FileSize  int64     `json:"file_size,omitempty"`
    Timestamp time.Time `json:"timestamp"`
}

// TypingIndicator represents typing status
type TypingIndicator struct {
    RoomID   uuid.UUID `json:"room_id"`
    UserID   uuid.UUID `json:"user_id"`
    Username string    `json:"username"`
    IsTyping bool      `json:"is_typing"`
}