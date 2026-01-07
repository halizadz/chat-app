package handlers

import (
    "encoding/json"
    "net/http"

    "github.com/halizadz/chat-app-backend/internal/repository"
	 "github.com/halizadz/chat-app-backend/internal/models" 
    "github.com/halizadz/chat-app-backend/internal/middleware"
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

    // Get search query if provided
    searchQuery := r.URL.Query().Get("search")
    
    var users []*models.User
    var err error
    
    if searchQuery != "" {
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