package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"unicode"

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

	// Validate email format
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	// Validate password strength (min 6 characters, at least one letter and one number)
	if len(req.Password) < 6 {
		http.Error(w, "Password must be at least 6 characters long", http.StatusBadRequest)
		return
	}

	hasLetter := false
	hasNumber := false
	for _, char := range req.Password {
		if unicode.IsLetter(char) {
			hasLetter = true
		}
		if unicode.IsNumber(char) {
			hasNumber = true
		}
	}

	if !hasLetter || !hasNumber {
		http.Error(w, "Password must contain at least one letter and one number", http.StatusBadRequest)
		return
	}

	// Validate username (alphanumeric and underscore, 3-20 characters)
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
	if !usernameRegex.MatchString(req.Username) {
		http.Error(w, "Username must be 3-20 characters and contain only letters, numbers, and underscores", http.StatusBadRequest)
		return
	}

	// Check if email already exists
	existingUser, _ := h.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		http.Error(w, "Email already registered", http.StatusConflict)
		return
	}

	// Check if username already exists (need to add FindByUsername method)
	// For now, we'll let database constraint handle it and return better error

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
		// Check for duplicate key error
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			if strings.Contains(err.Error(), "email") {
				http.Error(w, "Email already registered", http.StatusConflict)
			} else if strings.Contains(err.Error(), "username") {
				http.Error(w, "Username already taken", http.StatusConflict)
			} else {
				http.Error(w, "User already exists", http.StatusConflict)
			}
			return
		}
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
