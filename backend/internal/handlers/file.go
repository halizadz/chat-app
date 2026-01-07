package handlers

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"

    "github.com/google/uuid"
    "github.com/halizadz/chat-app-backend/internal/middleware"
)

type FileHandler struct {
    uploadDir string
}

func NewFileHandler(uploadDir string) *FileHandler {
    // Create upload directory if it doesn't exist
    os.MkdirAll(uploadDir, os.ModePerm)
    return &FileHandler{uploadDir: uploadDir}
}

type FileUploadResponse struct {
    URL      string `json:"url"`
    FileName string `json:"file_name"`
    FileSize int64  `json:"file_size"`
}

func (h *FileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
    claims, ok := middleware.GetUserFromContext(r.Context())
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // Parse multipart form
    err := r.ParseMultipartForm(10 << 20) // 10 MB max
    if err != nil {
        http.Error(w, "File too large", http.StatusBadRequest)
        return
    }

    file, handler, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "Error retrieving file", http.StatusBadRequest)
        return
    }
    defer file.Close()

    // Generate unique filename
    ext := filepath.Ext(handler.Filename)
    filename := fmt.Sprintf("%s-%s%s", claims.UserID.String(), uuid.New().String(), ext)
    filepath := filepath.Join(h.uploadDir, filename)

    // Create file
    dst, err := os.Create(filepath)
    if err != nil {
        http.Error(w, "Error saving file", http.StatusInternalServerError)
        return
    }
    defer dst.Close()

    // Copy file content
    fileSize, err := io.Copy(dst, file)
    if err != nil {
        http.Error(w, "Error saving file", http.StatusInternalServerError)
        return
    }

    // Return file URL
    fileURL := fmt.Sprintf("/uploads/%s", filename)

    response := FileUploadResponse{
        URL:      fileURL,
        FileName: handler.Filename,
        FileSize: fileSize,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}