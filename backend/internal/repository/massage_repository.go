package repository

import (
    "database/sql"
    "time"

    "github.com/google/uuid"
    "github.com/halizadz/chat-app-backend/internal/models"
)

type MessageRepository struct {
    db *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
    return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(message *models.Message) error {
    query := `
        INSERT INTO messages (id, room_id, sender_id, content, type, file_url, file_name, file_size, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING id, created_at
    `
    
    message.ID = uuid.New()
    now := time.Now()
    
    return r.db.QueryRow(
        query,
        message.ID,
        message.RoomID,
        message.SenderID,
        message.Content,
        message.Type,
        message.FileURL,
        message.FileName,
        message.FileSize,
        now,
        now,
    ).Scan(&message.ID, &message.CreatedAt)
}

func (r *MessageRepository) GetByRoomID(roomID uuid.UUID, limit, offset int) ([]*models.Message, error) {
    query := `
        SELECT m.id, m.room_id, m.sender_id, m.content, m.type, m.file_url, 
               m.file_name, m.file_size, m.created_at,
               u.id, u.username, u.email, u.avatar_url
        FROM messages m
        JOIN users u ON m.sender_id = u.id
        WHERE m.room_id = $1
        ORDER BY m.created_at DESC
        LIMIT $2 OFFSET $3
    `
    
    rows, err := r.db.Query(query, roomID, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var messages []*models.Message
    for rows.Next() {
        msg := &models.Message{
            Sender: &models.User{},
        }
        
        err := rows.Scan(
            &msg.ID,
            &msg.RoomID,
            &msg.SenderID,
            &msg.Content,
            &msg.Type,
            &msg.FileURL,
            &msg.FileName,
            &msg.FileSize,
            &msg.CreatedAt,
            &msg.Sender.ID,
            &msg.Sender.Username,
            &msg.Sender.Email,
            &msg.Sender.AvatarURL,
        )
        if err != nil {
            return nil, err
        }
        
        messages = append(messages, msg)
    }
    
    // Reverse to get chronological order
    for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
        messages[i], messages[j] = messages[j], messages[i]
    }
    
    return messages, nil
}

func (r *MessageRepository) MarkAsRead(messageID, userID uuid.UUID) error {
    query := `
        INSERT INTO message_read_status (id, message_id, user_id, read_at)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (message_id, user_id) DO NOTHING
    `
    
    _, err := r.db.Exec(query, uuid.New(), messageID, userID, time.Now())
    return err
}