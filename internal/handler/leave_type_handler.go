package handler

import (
	"hrsync-backend/internal/middleware"
	"hrsync-backend/internal/service"
	"hrsync-backend/internal/utils"
	"net/http"
)

type LeaveTypeHandler struct {
	srv service.LeaveTypeService
}

func NewLeaveTypeHandler(srv service.LeaveTypeService) *LeaveTypeHandler {
	return &LeaveTypeHandler{srv: srv}
}

func (h *LeaveTypeHandler) GetLeaveTypes(w http.ResponseWriter, r *http.Request) {
	responses, err := h.srv.GetLeaveTypes(r.Context())
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendSuccess(w, "Leave types retrieved successfully", responses, http.StatusOK)
}

func (h *LeaveTypeHandler) GetMyCredits(w http.ResponseWriter, r *http.Request) {
	email, ok := r.Context().Value(middleware.ContextKeyEmail).(string)
	if !ok {
		utils.SendError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	responses, err := h.srv.GetMyCredits(r.Context(), email)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendSuccess(w, "Leave credits retrieved successfully", responses, http.StatusOK)
}

func (h *LeaveTypeHandler) SeedLeaveTypes(w http.ResponseWriter, r *http.Request) {
	if err := h.srv.SeedLeaveTypes(r.Context()); err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendSuccess(w, "Leave types seeded successfully", nil, http.StatusOK)
}
