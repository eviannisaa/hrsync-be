package router

import (
	"hrsync-backend/internal/handler"
	"hrsync-backend/internal/middleware"
	"net/http"
)

func RegisterLeaveRoutes(mux *http.ServeMux, leaveHandler *handler.LeaveHandler, leaveTypeHandler *handler.LeaveTypeHandler) {
	mux.Handle("GET /api/leaves", middleware.AuthMiddleware(http.HandlerFunc(leaveHandler.GetLeaves)))
	mux.Handle("POST /api/leaves", middleware.AuthMiddleware(http.HandlerFunc(leaveHandler.CreateLeave)))
	mux.Handle("PUT /api/leaves/{id}", middleware.AuthMiddleware(http.HandlerFunc(leaveHandler.UpdateLeave)))
	mux.Handle("PUT /api/leaves/{id}/approve", middleware.AuthMiddleware(http.HandlerFunc(leaveHandler.ApproveLeave)))
	mux.Handle("PUT /api/leaves/{id}/reject", middleware.AuthMiddleware(http.HandlerFunc(leaveHandler.RejectLeave)))
	mux.Handle("GET /api/leaves/summary", middleware.AuthMiddleware(http.HandlerFunc(leaveHandler.GetLeaveSummary)))
	mux.Handle("DELETE /api/leaves/{id}", middleware.AuthMiddleware(http.HandlerFunc(leaveHandler.DeleteLeave)))

	// My routes
	mux.Handle("GET /api/leaves/credits", middleware.AuthMiddleware(http.HandlerFunc(leaveTypeHandler.GetMyCredits)))

	// Leave types
	mux.Handle("GET /api/leaves/types", middleware.AuthMiddleware(http.HandlerFunc(leaveTypeHandler.GetLeaveTypes)))
	mux.Handle("POST /api/leaves/types/seed", middleware.AuthMiddleware(http.HandlerFunc(leaveTypeHandler.SeedLeaveTypes)))
}
