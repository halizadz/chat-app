package main

import (
    "log"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/halizadz/chat-app-backend/internal/config"
    "github.com/halizadz/chat-app-backend/internal/database"
    "github.com/halizadz/chat-app-backend/internal/handlers"
    "github.com/halizadz/chat-app-backend/internal/middleware"
    "github.com/halizadz/chat-app-backend/internal/repository"
    "github.com/halizadz/chat-app-backend/internal/websocket"
)

func main() {
    cfg, err := config.Load()
    if err != nil {
        log.Fatal("Error loading config:", err)
    }

    db, err := database.NewDatabase(cfg.DatabaseURL)
    if err != nil {
        log.Fatal("Error connecting to database:", err)
    }
    defer db.Close()

    userRepo := repository.NewUserRepository(db.DB)
    roomRepo := repository.NewRoomRepository(db.DB)
    messageRepo := repository.NewMessageRepository(db.DB)

    hub := websocket.NewHub()
    go hub.Run()

    authHandler := handlers.NewAuthHandler(userRepo, cfg.JWTSecret)
    chatHandler := handlers.NewChatHandler(roomRepo, messageRepo, userRepo)
    wsHandler := handlers.NewWebSocketHandler(hub, roomRepo, messageRepo, cfg.JWTSecret)
    fileHandler := handlers.NewFileHandler("./uploads")
    userHandler := handlers.NewUserHandler(userRepo)

    r := mux.NewRouter()
    r.Use(middleware.CORS)

    // Public routes
    r.HandleFunc("/api/auth/register", authHandler.Register).Methods("POST", "OPTIONS")
    r.HandleFunc("/api/auth/login", authHandler.Login).Methods("POST", "OPTIONS")

    // WebSocket route - TANPA auth middleware, karena token di query param
    r.HandleFunc("/api/ws/{roomId}", wsHandler.HandleWebSocket).Methods("GET")

    r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

    // Protected routes
    api := r.PathPrefix("/api").Subrouter()
    api.Use(middleware.AuthMiddleware(cfg.JWTSecret))

    api.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET", "OPTIONS")
    api.HandleFunc("/users/me", userHandler.GetUserProfile).Methods("GET", "OPTIONS")

    api.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET", "OPTIONS")
    api.HandleFunc("/users/me", userHandler.GetUserProfile).Methods("GET", "OPTIONS")
    api.HandleFunc("/users/me", userHandler.UpdateUserProfile).Methods("PUT", "OPTIONS")

    api.HandleFunc("/rooms", chatHandler.GetUserRooms).Methods("GET", "OPTIONS")
    api.HandleFunc("/rooms", chatHandler.CreateRoom).Methods("POST", "OPTIONS")
    api.HandleFunc("/rooms/private", chatHandler.CreateOrGetPrivateRoom).Methods("POST", "OPTIONS")
    api.HandleFunc("/rooms/{roomId}", chatHandler.GetRoom).Methods("GET", "OPTIONS")
    api.HandleFunc("/rooms/{roomId}", chatHandler.UpdateRoom).Methods("PUT", "OPTIONS")
    api.HandleFunc("/rooms/{roomId}", chatHandler.DeleteRoom).Methods("DELETE", "OPTIONS")
    api.HandleFunc("/rooms/{roomId}/messages", chatHandler.GetRoomMessages).Methods("GET", "OPTIONS")
    api.HandleFunc("/rooms/{roomId}/messages/search", chatHandler.SearchMessages).Methods("GET", "OPTIONS")
    api.HandleFunc("/rooms/{roomId}/read", chatHandler.MarkRoomAsRead).Methods("POST", "OPTIONS")
    api.HandleFunc("/rooms/{roomId}/members", chatHandler.GetRoomMembers).Methods("GET", "OPTIONS")
    api.HandleFunc("/rooms/{roomId}/members", chatHandler.AddRoomMember).Methods("POST", "OPTIONS")
    api.HandleFunc("/rooms/{roomId}/members/{userId}", chatHandler.RemoveRoomMember).Methods("DELETE", "OPTIONS")
    api.HandleFunc("/rooms/{roomId}/leave", chatHandler.LeaveRoom).Methods("POST", "OPTIONS")

    api.HandleFunc("/messages/{messageId}", chatHandler.UpdateMessage).Methods("PUT", "OPTIONS")
    api.HandleFunc("/messages/{messageId}", chatHandler.DeleteMessage).Methods("DELETE", "OPTIONS")

    api.HandleFunc("/upload", fileHandler.UploadFile).Methods("POST", "OPTIONS")

    log.Printf("Server starting on port %s", cfg.Port)
    log.Printf("CORS enabled for all origins")
    log.Printf("WebSocket available at ws://localhost:%s/api/ws/{roomId}?token=...", cfg.Port)
    
    if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
        log.Fatal("Error starting server:", err)
    }
}