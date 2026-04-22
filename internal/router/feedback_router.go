package router

import (
	"hrsync-backend/internal/handler"
	"hrsync-backend/internal/middleware"
	"net/http"
)

func RegisterFeedbackRoutes(mux *http.ServeMux, h *handler.FeedbackHandler) {
	mux.Handle("GET /api/feedbacks", middleware.AuthMiddleware(http.HandlerFunc(h.GetFeedbacks)))
	mux.Handle("POST /api/feedbacks", middleware.AuthMiddleware(http.HandlerFunc(h.CreateFeedback)))
	mux.Handle("DELETE /api/feedbacks/{id}", middleware.AuthMiddleware(http.HandlerFunc(h.DeleteFeedback)))
}
