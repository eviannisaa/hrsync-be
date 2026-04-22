package router

import (
	"hrsync-backend/internal/handler"
	"hrsync-backend/internal/middleware"
	"net/http"
)

func RegisterDepartmentRoutes(mux *http.ServeMux, departmentHandler *handler.DepartmentHandler) {
	mux.Handle("GET /api/departments", middleware.AuthMiddleware(http.HandlerFunc(departmentHandler.GetDepartments)))
}
