package handler

import (
	"encoding/json"
	"hrsync-backend/internal/dto"
	"hrsync-backend/internal/service"
	"hrsync-backend/internal/utils"
	"net/http"
)

type OvertimeHandler struct {
	srv service.OvertimeService
}

func NewOvertimeHandler(srv service.OvertimeService) *OvertimeHandler {
	return &OvertimeHandler{srv: srv}
}

func (h *OvertimeHandler) GetOvertimes(w http.ResponseWriter, r *http.Request) {
	params := utils.GetListParams(r)

	// Filter by email if employee, or if admin provides email param.
	// If admin and no email param provided, Email will be cleared to show all.
	if params.Role == "ADMIN" {
		params.Email = r.URL.Query().Get("email")
	}

	responses, total, err := h.srv.GetOvertimes(r.Context(), params)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendPaginated(w, "Overtime retrieved successfully", responses, total, params.Page, params.Limit)
}

func (h *OvertimeHandler) CreateOvertime(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateOvertimeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.srv.CreateOvertimes(r.Context(), req)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(w, "Overtime created successfully", response, http.StatusCreated)
}

func (h *OvertimeHandler) UpdateOvertime(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req dto.UpdateOvertimeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.srv.UpdateOvertimes(r.Context(), id, req)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(w, "Overtime updated successfully", response, http.StatusOK)
}

func (h *OvertimeHandler) DeleteOvertime(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.srv.DeleteOvertimes(r.Context(), id); err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(w, "Overtime deleted successfully", nil, http.StatusOK)
}

func (h *OvertimeHandler) ApproveOvertime(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	response, err := h.srv.ApproveOvertimes(r.Context(), id)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(w, "Overime approved successfully", response, http.StatusOK)
}

func (h *OvertimeHandler) RejectOvertimes(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	response, err := h.srv.RejectOvertimes(r.Context(), id)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(w, "Overtime rejected successfully", response, http.StatusOK)
}
