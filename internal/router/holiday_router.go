package router

import (
	"hrsync-backend/internal/handler"
	"hrsync-backend/internal/middleware"
	"net/http"
)

func RegisterHolidayRoutes(mux *http.ServeMux, h *handler.HolidayHandler) {
	mux.Handle("GET /api/holidays", middleware.AuthMiddleware(http.HandlerFunc(h.GetHolidays)))
	mux.Handle("POST /api/holidays/sync", middleware.AuthMiddleware(http.HandlerFunc(h.SyncHolidays)))
}
