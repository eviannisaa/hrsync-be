package router

import (
	"hrsync-backend/internal/handler"
	"hrsync-backend/internal/middleware"
	"net/http"
)

func RegisterKPIRoutes(mux *http.ServeMux, h *handler.TemplateKPIHandler) {
	mux.Handle("GET /api/kpi", middleware.AuthMiddleware(http.HandlerFunc(h.GetTemplatesKPI)))
	mux.Handle("GET /api/kpi/published", middleware.AuthMiddleware(http.HandlerFunc(h.GetPublishedKPIByDepartment)))
	mux.Handle("POST /api/kpi", middleware.AuthMiddleware(http.HandlerFunc(h.CreateTemplateKPI)))
	mux.Handle("PUT /api/kpi/{id}", middleware.AuthMiddleware(http.HandlerFunc(h.UpdateTemplateKPI)))
	mux.Handle("DELETE /api/kpi/{id}", middleware.AuthMiddleware(http.HandlerFunc(h.DeleteTemplateKPI)))
}
