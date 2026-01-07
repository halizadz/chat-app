package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/halizadz/chat-app-backend/internal/middleware"
	"github.com/halizadz/chat-app-backend/internal/models"
	"github.com/halizadz/chat-app-backend/internal/repository"
)

type UserHandler struct {
	userRepo *repository.UserRepository
}

func NewUserHandler(userRepo *repository.UserRepository) *UserHandler {
	return &UserHandler{userRepo: userRepo}
}

// GetAllUsers returns all users (excluding current user)
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get search query parameters
	searchQuery := r.URL.Query().Get("search")
	gmailQuery := r.URL.Query().Get("gmail") // Parameter khusus untuk pencarian Gmail

	var users []*models.User
	var err error

	// Jika ada parameter gmail, gunakan pencarian khusus Gmail
	if gmailQuery != "" {
		users, err = h.userRepo.SearchUsersByGmail(gmailQuery)
	} else if searchQuery != "" {
		users, err = h.userRepo.SearchUsers(searchQuery)
	} else {
		users, err = h.userRepo.FindAll()
	}

	if err != nil {
		http.Error(w, "Error fetching users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Filter out current user
	filteredUsers := []*models.User{}
	for _, user := range users {
		if user.ID != claims.UserID {
			filteredUsers = append(filteredUsers, user)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filteredUsers)
}

// GetUserProfile returns current user profile
func (h *UserHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.userRepo.FindByID(claims.UserID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Clear password hash
	user.PasswordHash = ""

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// UpdateUserProfile updates current user profile
func (h *UserHandler) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Username  *string `json:"username"`
		Email     *string `json:"email"`
		AvatarURL *string `json:"avatar_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get current user
	user, err := h.userRepo.FindByID(claims.UserID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Update fields if provided
	if req.Username != nil && *req.Username != "" {
		user.Username = *req.Username
	}
	if req.Email != nil && *req.Email != "" {
		// Check if email already exists (if changed)
		if *req.Email != user.Email {
			existingUser, _ := h.userRepo.FindByEmail(*req.Email)
			if existingUser != nil {
				http.Error(w, "Email already registered", http.StatusConflict)
				return
			}
		}
		user.Email = *req.Email
	}
	if req.AvatarURL != nil {
		user.AvatarURL = req.AvatarURL
	}

	// Update in database
	if err := h.userRepo.Update(user); err != nil {
		http.Error(w, "Error updating profile: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Clear password hash
	user.PasswordHash = ""

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
