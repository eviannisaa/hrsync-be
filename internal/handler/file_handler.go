package handler

import (
	"fmt"
	"hrsync-backend/internal/utils"
	"net/http"
	"path/filepath"
	"time"
)

type FileHandler struct{}

func NewFileHandler() *FileHandler {
	return &FileHandler{}
}

func (h *FileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.SendError(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		utils.SendError(w, "File is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Generate object name
	ext := filepath.Ext(header.Filename)
	objectName := fmt.Sprintf("uploads/%d%s", time.Now().UnixNano(), ext)
	contentType := header.Header.Get("Content-Type")
	// Upload to storage directly without conversion
	key, err := utils.Upload(r.Context(), file, header.Size, objectName, contentType)

	if err != nil {
		utils.SendError(w, "Upload failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(w, "File uploaded successfully", map[string]string{
		"key": key,
		"url": utils.GetURL(key),
	}, http.StatusOK)
}
