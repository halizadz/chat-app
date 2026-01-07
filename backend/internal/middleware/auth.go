package middleware

import (
    "context"
    "net/http"
    "strings"

    "github.com/halizadz/chat-app-backend/internal/utils"
)

type contextKey string

const UserContextKey contextKey = "user"

func AuthMiddleware(secret string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Skip auth for OPTIONS requests (preflight)
            if r.Method == "OPTIONS" {
                next.ServeHTTP(w, r)
                return
            }

            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                http.Error(w, "Authorization header required", http.StatusUnauthorized)
                return
            }

            bearerToken := strings.Split(authHeader, " ")
            if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
                http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
                return
            }

            claims, err := utils.ValidateToken(bearerToken[1], secret)
            if err != nil {
                http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
                return
            }

            ctx := context.WithValue(r.Context(), UserContextKey, claims)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func GetUserFromContext(ctx context.Context) (*utils.Claims, bool) {
    claims, ok := ctx.Value(UserContextKey).(*utils.Claims)
    return claims, ok
}