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
               m.file_name, m.file_size, m.created_at, m.updated_at,
               COALESCE(m.updated_at > m.created_at, false) as is_edited,
               COALESCE(m.content = '[DELETED]', false) as is_deleted,
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

		var isEdited, isDeleted bool

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
			&msg.UpdatedAt,
			&isEdited,
			&isDeleted,
			&msg.Sender.ID,
			&msg.Sender.Username,
			&msg.Sender.Email,
			&msg.Sender.AvatarURL,
		)
		if err != nil {
			return nil, err
		}

		msg.IsEdited = isEdited
		msg.IsDeleted = isDeleted

		// Get read status
		readBy, _ := r.GetReadBy(msg.ID)
		msg.ReadBy = readBy

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

// MarkRoomMessagesAsRead marks all messages in a room as read for a user
func (r *MessageRepository) MarkRoomMessagesAsRead(roomID, userID uuid.UUID) error {
	query := `
        INSERT INTO message_read_status (id, message_id, user_id, read_at)
        SELECT uuid_generate_v4(), m.id, $2, NOW()
        FROM messages m
        WHERE m.room_id = $1
        AND NOT EXISTS (
            SELECT 1 FROM message_read_status mrs 
            WHERE mrs.message_id = m.id AND mrs.user_id = $2
        )
    `

	_, err := r.db.Exec(query, roomID, userID)
	return err
}

// GetReadBy returns list of user IDs who read the message
func (r *MessageRepository) GetReadBy(messageID uuid.UUID) ([]uuid.UUID, error) {
	query := `
        SELECT user_id FROM message_read_status
        WHERE message_id = $1
    `

	rows, err := r.db.Query(query, messageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var readBy []uuid.UUID
	for rows.Next() {
		var userID uuid.UUID
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		readBy = append(readBy, userID)
	}

	return readBy, nil
}

// FindByID finds a message by ID
func (r *MessageRepository) FindByID(messageID uuid.UUID) (*models.Message, error) {
	msg := &models.Message{
		Sender: &models.User{},
	}

	query := `
        SELECT m.id, m.room_id, m.sender_id, m.content, m.type, m.file_url, 
               m.file_name, m.file_size, m.created_at, m.updated_at,
               COALESCE(m.updated_at > m.created_at, false) as is_edited,
               COALESCE(m.content = '[DELETED]', false) as is_deleted,
               u.id, u.username, u.email, u.avatar_url
        FROM messages m
        JOIN users u ON m.sender_id = u.id
        WHERE m.id = $1
    `

	var isEdited, isDeleted bool
	err := r.db.QueryRow(query, messageID).Scan(
		&msg.ID,
		&msg.RoomID,
		&msg.SenderID,
		&msg.Content,
		&msg.Type,
		&msg.FileURL,
		&msg.FileName,
		&msg.FileSize,
		&msg.CreatedAt,
		&msg.UpdatedAt,
		&isEdited,
		&isDeleted,
		&msg.Sender.ID,
		&msg.Sender.Username,
		&msg.Sender.Email,
		&msg.Sender.AvatarURL,
	)

	if err != nil {
		return nil, err
	}

	msg.IsEdited = isEdited
	msg.IsDeleted = isDeleted

	readBy, _ := r.GetReadBy(msg.ID)
	msg.ReadBy = readBy

	return msg, nil
}

// Update updates a message
func (r *MessageRepository) Update(messageID uuid.UUID, content string) error {
	query := `
        UPDATE messages 
        SET content = $1, updated_at = NOW()
        WHERE id = $2
    `

	_, err := r.db.Exec(query, content, messageID)
	return err
}

// Delete soft deletes a message
func (r *MessageRepository) Delete(messageID uuid.UUID) error {
	query := `
        UPDATE messages 
        SET content = '[DELETED]', updated_at = NOW()
        WHERE id = $1
    `

	_, err := r.db.Exec(query, messageID)
	return err
}

// SearchMessages searches messages in a room
func (r *MessageRepository) SearchMessages(roomID uuid.UUID, query string, limit, offset int) ([]*models.Message, error) {
	searchQuery := `
        SELECT m.id, m.room_id, m.sender_id, m.content, m.type, m.file_url, 
               m.file_name, m.file_size, m.created_at, m.updated_at,
               COALESCE(m.updated_at > m.created_at, false) as is_edited,
               COALESCE(m.content = '[DELETED]', false) as is_deleted,
               u.id, u.username, u.email, u.avatar_url
        FROM messages m
        JOIN users u ON m.sender_id = u.id
        WHERE m.room_id = $1 
        AND m.content ILIKE $2
        AND m.content != '[DELETED]'
        ORDER BY m.created_at DESC
        LIMIT $3 OFFSET $4
    `

	searchPattern := "%" + query + "%"
	rows, err := r.db.Query(searchQuery, roomID, searchPattern, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		msg := &models.Message{
			Sender: &models.User{},
		}

		var isEdited, isDeleted bool

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
			&msg.UpdatedAt,
			&isEdited,
			&isDeleted,
			&msg.Sender.ID,
			&msg.Sender.Username,
			&msg.Sender.Email,
			&msg.Sender.AvatarURL,
		)
		if err != nil {
			return nil, err
		}

		msg.IsEdited = isEdited
		msg.IsDeleted = isDeleted

		messages = append(messages, msg)
	}

	return messages, nil
}
