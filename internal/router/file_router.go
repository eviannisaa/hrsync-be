package router

import (
	"hrsync-backend/internal/handler"
	"hrsync-backend/internal/middleware"
	"net/http"
)

func RegisterFileRoutes(mux *http.ServeMux, h *handler.FileHandler) {
	mux.Handle("POST /api/upload", middleware.AuthMiddleware(http.HandlerFunc(h.UploadFile)))
}
