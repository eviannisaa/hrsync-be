package router

import (
	"hrsync-backend/internal/handler"
	"hrsync-backend/internal/middleware"
	"net/http"
)

func RegisterOvertimeRoutes(mux *http.ServeMux, overtimeHandler *handler.OvertimeHandler) {
	mux.Handle("GET /api/overtime", middleware.AuthMiddleware(http.HandlerFunc(overtimeHandler.GetOvertimes)))
	mux.Handle("POST /api/overtime", middleware.AuthMiddleware(http.HandlerFunc(overtimeHandler.CreateOvertime)))
	mux.Handle("PUT /api/overtime/{id}", middleware.AuthMiddleware(http.HandlerFunc(overtimeHandler.UpdateOvertime)))
	mux.Handle("PUT /api/overtime/{id}/approve", middleware.AuthMiddleware(http.HandlerFunc(overtimeHandler.ApproveOvertime)))
	mux.Handle("PUT /api/overtime/{id}/reject", middleware.AuthMiddleware(http.HandlerFunc(overtimeHandler.RejectOvertimes)))
	mux.Handle("DELETE /api/overtime/{id}", middleware.AuthMiddleware(http.HandlerFunc(overtimeHandler.DeleteOvertime)))
}
