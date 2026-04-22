package router

import (
	"hrsync-backend/internal/handler"
	"hrsync-backend/internal/middleware"
	"net/http"
)

func RegisterPayslipRoutes(mux *http.ServeMux, h *handler.PayslipHandler) {
	mux.Handle("GET /api/payslips", middleware.AuthMiddleware(http.HandlerFunc(h.GetAll)))
	mux.Handle("POST /api/payslips/bulk", middleware.AuthMiddleware(middleware.RoleMiddleware("ADMIN")(http.HandlerFunc(h.BulkUpload))))
	mux.Handle("DELETE /api/payslips/batch", middleware.AuthMiddleware(middleware.RoleMiddleware("ADMIN")(http.HandlerFunc(h.DeleteBatch))))
	mux.Handle("DELETE /api/payslips/{id}", middleware.AuthMiddleware(middleware.RoleMiddleware("ADMIN")(http.HandlerFunc(h.Delete))))
}
