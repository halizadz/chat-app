package repository

import (
    "database/sql"
    "fmt"
    "time"

    "github.com/google/uuid"
    "github.com/halizadz/chat-app-backend/internal/models"
)

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
    query := `
        INSERT INTO users (id, username, email, password_hash, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, created_at, updated_at
    `
    
    user.ID = uuid.New()
    now := time.Now()
    
    return r.db.QueryRow(
        query,
        user.ID,
        user.Username,
        user.Email,
        user.PasswordHash,
        now,
        now,
    ).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
    user := &models.User{}
    query := `
        SELECT id, username, email, password_hash, avatar_url, status, last_seen, created_at, updated_at
        FROM users WHERE email = $1
    `
    
    err := r.db.QueryRow(query, email).Scan(
        &user.ID,
        &user.Username,
        &user.Email,
        &user.PasswordHash,
        &user.AvatarURL,
        &user.Status,
        &user.LastSeen,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("user not found")
    }
    
    return user, err
}

func (r *UserRepository) FindByID(id uuid.UUID) (*models.User, error) {
    user := &models.User{}
    query := `
        SELECT id, username, email, password_hash, avatar_url, status, last_seen, created_at, updated_at
        FROM users WHERE id = $1
    `
    
    err := r.db.QueryRow(query, id).Scan(
        &user.ID,
        &user.Username,
        &user.Email,
        &user.PasswordHash,
        &user.AvatarURL,
        &user.Status,
        &user.LastSeen,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("user not found")
    }
    
    return user, err
}

func (r *UserRepository) UpdateStatus(userID uuid.UUID, status string) error {
    query := `UPDATE users SET status = $1, last_seen = $2, updated_at = $3 WHERE id = $4`
    _, err := r.db.Exec(query, status, time.Now(), time.Now(), userID)
    return err
}

// FindAll - Method baru untuk get semua users
func (r *UserRepository) FindAll() ([]*models.User, error) {
    query := `
        SELECT id, username, email, avatar_url, status, last_seen, created_at, updated_at
        FROM users
        ORDER BY username ASC
    `
    
    rows, err := r.db.Query(query)
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
            &user.CreatedAt,
            &user.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        // Don't return password hash
        user.PasswordHash = ""
        users = append(users, user)
    }
    
    return users, nil
}

// SearchUsers - Method tambahan untuk search users by username atau email
func (r *UserRepository) SearchUsers(query string) ([]*models.User, error) {
    searchQuery := `
        SELECT id, username, email, avatar_url, status, last_seen, created_at, updated_at
        FROM users
        WHERE username ILIKE $1 OR email ILIKE $1
        ORDER BY username ASC
        LIMIT 50
    `
    
    searchPattern := "%" + query + "%"
    rows, err := r.db.Query(searchQuery, searchPattern)
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
            &user.CreatedAt,
            &user.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        user.PasswordHash = ""
        users = append(users, user)
    }
    
    return users, nil
}