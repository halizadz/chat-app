package handlers

import (
    "encoding/json"
    "net/http"

    "github.com/halizadz/chat-app-backend/internal/models"
    "github.com/halizadz/chat-app-backend/internal/repository"
    "github.com/halizadz/chat-app-backend/internal/utils"
)

type AuthHandler struct {
    userRepo  *repository.UserRepository
    jwtSecret string
}

func NewAuthHandler(userRepo *repository.UserRepository, jwtSecret string) *AuthHandler {
    return &AuthHandler{
        userRepo:  userRepo,
        jwtSecret: jwtSecret,
    }
}

type RegisterRequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type AuthResponse struct {
    Token string       `json:"token"`
    User  *models.User `json:"user"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    var req RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Validate input
    if req.Username == "" || req.Email == "" || req.Password == "" {
        http.Error(w, "All fields are required", http.StatusBadRequest)
        return
    }

    // Hash password
    hashedPassword, err := utils.HashPassword(req.Password)
    if err != nil {
        http.Error(w, "Error hashing password", http.StatusInternalServerError)
        return
    }

    // Create user
    user := &models.User{
        Username:     req.Username,
        Email:        req.Email,
        PasswordHash: hashedPassword,
        Status:       "offline",
    }

    if err := h.userRepo.Create(user); err != nil {
        http.Error(w, "Error creating user: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Generate token
    token, err := utils.GenerateToken(user.ID, user.Username, user.Email, h.jwtSecret)
    if err != nil {
        http.Error(w, "Error generating token", http.StatusInternalServerError)
        return
    }

    // Clear password hash before sending
    user.PasswordHash = ""

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(AuthResponse{
        Token: token,
        User:  user,
    })
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Find user by email
    user, err := h.userRepo.FindByEmail(req.Email)
    if err != nil {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    // Check password
    if !utils.CheckPassword(req.Password, user.PasswordHash) {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    // Generate token
    token, err := utils.GenerateToken(user.ID, user.Username, user.Email, h.jwtSecret)
    if err != nil {
        http.Error(w, "Error generating token", http.StatusInternalServerError)
        return
    }

    // Update user status to online
    h.userRepo.UpdateStatus(user.ID, "online")

    // Clear password hash before sending
    user.PasswordHash = ""

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(AuthResponse{
        Token: token,
        User:  user,
    })
}