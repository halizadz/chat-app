package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID  `json:"id"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	AvatarURL    *string    `json:"avatar_url"`
	Status       string     `json:"status"`
	LastSeen     *time.Time `json:"last_seen"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type Room struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Type        string    `json:"type"` // private, group
	CreatedBy   uuid.UUID `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Message struct {
	ID        uuid.UUID   `json:"id"`
	RoomID    uuid.UUID   `json:"room_id"`
	SenderID  uuid.UUID   `json:"sender_id"`
	Content   string      `json:"content"`
	Type      string      `json:"type"`
	FileURL   *string     `json:"file_url,omitempty"`
	FileName  *string     `json:"file_name,omitempty"`
	FileSize  *int64      `json:"file_size,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	IsEdited  bool        `json:"is_edited"`
	IsDeleted bool        `json:"is_deleted"`
	Sender    *User       `json:"sender,omitempty"`
	ReadBy    []uuid.UUID `json:"read_by,omitempty"` // Users who read this message
}

type RoomMember struct {
	ID       uuid.UUID `json:"id"`
	RoomID   uuid.UUID `json:"room_id"`
	UserID   uuid.UUID `json:"user_id"`
	Role     string    `json:"role"` // admin, member
	JoinedAt time.Time `json:"joined_at"`
}

// WebSocket Message Types
type WSMessage struct {
	Type    string      `json:"type"` // message, typing, join, leave
	Payload interface{} `json:"payload"`
}

type TypingIndicator struct {
	RoomID   uuid.UUID `json:"room_id"`
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	IsTyping bool      `json:"is_typing"`
}
