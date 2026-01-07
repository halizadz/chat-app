package repository

import (
    "database/sql"
    "fmt"
    "hash/fnv"
    "time"

    "github.com/google/uuid"
    "github.com/halizadz/chat-app-backend/internal/models"
)

type RoomRepository struct {
    db *sql.DB
}

func NewRoomRepository(db *sql.DB) *RoomRepository {
    return &RoomRepository{db: db}
}

func (r *RoomRepository) Create(room *models.Room) error {
    query := `
        INSERT INTO rooms (id, name, description, type, created_by, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, created_at, updated_at
    `
    
    room.ID = uuid.New()
    now := time.Now()
    
    return r.db.QueryRow(
        query,
        room.ID,
        room.Name,
        room.Description,
        room.Type,
        room.CreatedBy,
        now,
        now,
    ).Scan(&room.ID, &room.CreatedAt, &room.UpdatedAt)
}

func (r *RoomRepository) FindByID(id uuid.UUID) (*models.Room, error) {
    room := &models.Room{}
    query := `
        SELECT id, name, description, type, created_by, created_at, updated_at
        FROM rooms WHERE id = $1
    `
    
    err := r.db.QueryRow(query, id).Scan(
        &room.ID,
        &room.Name,
        &room.Description,
        &room.Type,
        &room.CreatedBy,
        &room.CreatedAt,
        &room.UpdatedAt,
    )
    
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("room not found")
    }
    
    return room, err
}

func (r *RoomRepository) GetUserRooms(userID uuid.UUID) ([]*models.Room, error) {
    query := `
        SELECT r.id, r.name, r.description, r.type, r.created_by, r.created_at, r.updated_at
        FROM rooms r
        JOIN room_members rm ON r.id = rm.room_id
        WHERE rm.user_id = $1
        ORDER BY r.updated_at DESC
    `
    
    rows, err := r.db.Query(query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var rooms []*models.Room
    for rows.Next() {
        room := &models.Room{}
        err := rows.Scan(
            &room.ID,
            &room.Name,
            &room.Description,
            &room.Type,
            &room.CreatedBy,
            &room.CreatedAt,
            &room.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        rooms = append(rooms, room)
    }
    
    return rooms, nil
}

func (r *RoomRepository) AddMember(roomID, userID uuid.UUID, role string) error {
    query := `
        INSERT INTO room_members (id, room_id, user_id, role, joined_at)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (room_id, user_id) DO NOTHING
    `
    
    _, err := r.db.Exec(query, uuid.New(), roomID, userID, role, time.Now())
    return err
}

func (r *RoomRepository) RemoveMember(roomID, userID uuid.UUID) error {
    query := `DELETE FROM room_members WHERE room_id = $1 AND user_id = $2`
    _, err := r.db.Exec(query, roomID, userID)
    return err
}

func (r *RoomRepository) Update(room *models.Room) error {
    query := `
        UPDATE rooms 
        SET name = $1, description = $2, updated_at = $3
        WHERE id = $4
    `
    _, err := r.db.Exec(query, room.Name, room.Description, time.Now(), room.ID)
    return err
}

func (r *RoomRepository) Delete(roomID uuid.UUID) error {
    query := `DELETE FROM rooms WHERE id = $1`
    _, err := r.db.Exec(query, roomID)
    return err
}

func (r *RoomRepository) IsMember(roomID, userID uuid.UUID) (bool, error) {
    var exists bool
    query := `SELECT EXISTS(SELECT 1 FROM room_members WHERE room_id = $1 AND user_id = $2)`
    err := r.db.QueryRow(query, roomID, userID).Scan(&exists)
    return exists, err
}

func (r *RoomRepository) GetMembers(roomID uuid.UUID) ([]*models.User, error) {
    query := `
        SELECT u.id, u.username, u.email, u.avatar_url, u.status, u.last_seen
        FROM users u
        JOIN room_members rm ON u.id = rm.user_id
        WHERE rm.room_id = $1
    `
    
    rows, err := r.db.Query(query, roomID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var users []*models.User
    for rows.Next() {
        user := &models.User{}
        err := rows.Scan(
            &user.ID,
            &user.Username,
            &user.Email,
            &user.AvatarURL,
            &user.Status,
            &user.LastSeen,
        )
        if err != nil {
            return nil, err
        }
        users = append(users, user)
    }
    
    return users, nil
}

// Find or create a private room between two users
// Uses advisory lock to prevent race condition when creating private rooms
// otherUsername: username of the other user (user2) for room naming
func (r *RoomRepository) FindOrCreatePrivateRoom(user1ID, user2ID uuid.UUID, otherUsername string) (*models.Room, error) {
    // Ensure consistent ordering of user IDs for lock key
    // This ensures same lock is used regardless of parameter order
    var lockKey int64
    id1Str := user1ID.String()
    id2Str := user2ID.String()
    
    // Create consistent key regardless of parameter order
    combinedStr := id1Str + "|" + id2Str
    if id1Str > id2Str {
        combinedStr = id2Str + "|" + id1Str
    }
    
    // Hash the combined string to get lock key
    h := fnv.New64a()
    h.Write([]byte(combinedStr))
    lockKey = int64(h.Sum64())
    
    tx, err := r.db.Begin()
    if err != nil {
        return nil, err
    }
    defer tx.Rollback()
    
    // Acquire advisory lock to prevent race condition
    // Lock will be released when transaction commits/rolls back
    _, err = tx.Exec("SELECT pg_advisory_xact_lock($1)", lockKey)
    if err != nil {
        return nil, err
    }
    
    // Now check again if room exists (double-check pattern)
    query := `
        SELECT r.id, r.name, r.description, r.type, r.created_by, r.created_at, r.updated_at
        FROM rooms r
        WHERE r.type = 'private'
        AND r.id IN (
            SELECT room_id FROM room_members WHERE user_id = $1
        )
        AND r.id IN (
            SELECT room_id FROM room_members WHERE user_id = $2
        )
        LIMIT 1
    `
    
    room := &models.Room{}
    err = tx.QueryRow(query, user1ID, user2ID).Scan(
        &room.ID,
        &room.Name,
        &room.Description,
        &room.Type,
        &room.CreatedBy,
        &room.CreatedAt,
        &room.UpdatedAt,
    )
    
    if err == nil {
        // Room found, commit transaction (releases lock)
        tx.Commit()
        return room, nil
    }
    
    if err != sql.ErrNoRows {
        return nil, err
    }
    
    // Create new private room with other user's username as name
    roomName := otherUsername
    if roomName == "" {
        roomName = "Private Chat" // Fallback jika username tidak tersedia
    }
    
    room = &models.Room{
        ID:        uuid.New(),
        Name:      roomName,
        Type:      "private",
        CreatedBy: user1ID,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    
    // Insert room
    insertRoom := `
        INSERT INTO rooms (id, name, description, type, created_by, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
    _, err = tx.Exec(insertRoom, room.ID, room.Name, room.Description, room.Type, room.CreatedBy, room.CreatedAt, room.UpdatedAt)
    if err != nil {
        return nil, err
    }
    
    // Add both members
    insertMember := `
        INSERT INTO room_members (id, room_id, user_id, role, joined_at)
        VALUES ($1, $2, $3, $4, $5)
    `
    _, err = tx.Exec(insertMember, uuid.New(), room.ID, user1ID, "member", time.Now())
    if err != nil {
        return nil, err
    }
    
    _, err = tx.Exec(insertMember, uuid.New(), room.ID, user2ID, "member", time.Now())
    if err != nil {
        return nil, err
    }
    
    if err = tx.Commit(); err != nil {
        return nil, err
    }
    
    return room, nil
}