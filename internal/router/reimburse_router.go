package router

import (
	"hrsync-backend/internal/handler"
	"hrsync-backend/internal/middleware"
	"net/http"
)

func RegisterReimburseRoutes(mux *http.ServeMux, h *handler.ReimburseHandler) {
	mux.Handle("GET /api/reimbursements", middleware.AuthMiddleware(http.HandlerFunc(h.GetReimbursements)))
	mux.Handle("POST /api/reimbursements", middleware.AuthMiddleware(http.HandlerFunc(h.CreateReimbursement)))
	mux.Handle("PUT /api/reimbursements/{id}", middleware.AuthMiddleware(http.HandlerFunc(h.UpdateReimbursement)))
	mux.Handle("PUT /api/reimbursements/{id}/approve", middleware.AuthMiddleware(http.HandlerFunc(h.ApproveReimbursement)))
	mux.Handle("PUT /api/reimbursements/{id}/reject", middleware.AuthMiddleware(http.HandlerFunc(h.RejectReimbursement)))
	mux.Handle("DELETE /api/reimbursements/{id}", middleware.AuthMiddleware(http.HandlerFunc(h.DeleteReimbursement)))
}
