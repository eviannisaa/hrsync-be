package router

import (
	"hrsync-backend/internal/handler"
	"hrsync-backend/internal/middleware"
	"net/http"
)

func RegisterTeamRoutes(mux *http.ServeMux, employeeHandler *handler.EmployeeHandler) {
	mux.Handle("GET /api/employees", middleware.AuthMiddleware(http.HandlerFunc(employeeHandler.GetEmployees)))
	mux.Handle("POST /api/employees", middleware.AuthMiddleware(http.HandlerFunc(employeeHandler.CreateEmployee)))
	mux.Handle("GET /api/employees/{id}", middleware.AuthMiddleware(http.HandlerFunc(employeeHandler.GetEmployeeByID)))
	mux.Handle("PUT /api/employees/{id}", middleware.AuthMiddleware(http.HandlerFunc(employeeHandler.UpdateEmployee)))
	mux.Handle("DELETE /api/employees/{id}", middleware.AuthMiddleware(http.HandlerFunc(employeeHandler.DeleteEmployee)))

	// Organization
	mux.Handle("GET /api/employee-organization", middleware.AuthMiddleware(http.HandlerFunc(employeeHandler.GetOrganization)))
	mux.Handle("PUT /api/employee-organization", middleware.AuthMiddleware(http.HandlerFunc(employeeHandler.UpdateOrganization)))
}
