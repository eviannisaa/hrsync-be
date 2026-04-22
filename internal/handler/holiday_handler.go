package handler

import (
	"hrsync-backend/internal/service"
	"hrsync-backend/internal/utils"
	"net/http"
	"strconv"
	"time"
)

type HolidayHandler struct {
	srv service.HolidayService
}

func NewHolidayHandler(srv service.HolidayService) *HolidayHandler {
	return &HolidayHandler{srv: srv}
}

func (h *HolidayHandler) GetHolidays(w http.ResponseWriter, r *http.Request) {
	holidays, err := h.srv.GetHolidays(r.Context())
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(w, "Holidays retrieved successfully", holidays, http.StatusOK)
}

func (h *HolidayHandler) SyncHolidays(w http.ResponseWriter, r *http.Request) {
	yearStr := r.URL.Query().Get("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		year = time.Now().Year()
	}

	if err := h.srv.SyncHolidays(r.Context(), year); err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(w, "Holidays synced successfully", nil, http.StatusOK)
}
