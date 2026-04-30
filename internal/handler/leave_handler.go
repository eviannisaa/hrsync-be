package handler

import (
	"encoding/json"
	"hrsync-backend/internal/dto"
	"hrsync-backend/internal/model"
	"hrsync-backend/internal/service"
	"hrsync-backend/internal/utils"
	"net/http"
)

type LeaveHandler struct {
	srv service.LeaveService
}

func NewLeaveHandler(srv service.LeaveService) *LeaveHandler {
	return &LeaveHandler{srv: srv}
}

func (h *LeaveHandler) GetLeaves(w http.ResponseWriter, r *http.Request) {
	params := utils.GetListParams(r)

	// If admin, they see all unless they specifically ask for their own via ?mine=true
	// or ask for a specific email via ?email=...
	if params.Role == "ADMIN" {
		isMine := r.URL.Query().Get("mine") == "true"
		if !isMine {
			params.Email = r.URL.Query().Get("email")
		}
	}

	responses, total, err := h.srv.GetLeaves(r.Context(), params)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendPaginated(w, "Leaves retrieved successfully", responses, total, params.Page, params.Limit)
}

func (h *LeaveHandler) CreateLeave(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateLeaveRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	email, _ := r.Context().Value(model.ContextKeyEmail).(string)
	role, _ := r.Context().Value(model.ContextKeyRole).(string)

	// If admin and email specified in body, use it. Otherwise use user's email.
	if role != "ADMIN" || req.Email == "" {
		req.Email = email
	}

	req.CreatedBy = email
	req.UpdatedBy = email

	response, err := h.srv.CreateLeaves(r.Context(), req)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(w, "Leave created successfully", response, http.StatusCreated)
}

func (h *LeaveHandler) UpdateLeave(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req dto.UpdateLeaveRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	email, _ := r.Context().Value(model.ContextKeyEmail).(string)
	role, _ := r.Context().Value(model.ContextKeyRole).(string)

	// If admin and email specified in body, use it. Otherwise use user's email.
	if role != "ADMIN" || req.Email == nil || *req.Email == "" {
		req.Email = &email
	}

	req.UpdatedBy = email

	response, err := h.srv.UpdateLeaves(r.Context(), id, req)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(w, "Leave updated successfully", response, http.StatusOK)
}

func (h *LeaveHandler) DeleteLeave(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.srv.DeleteLeaves(r.Context(), id); err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(w, "Leave deleted successfully", nil, http.StatusOK)
}

func (h *LeaveHandler) ApproveLeave(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	response, err := h.srv.ApproveLeave(r.Context(), id)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(w, "Leave approved successfully", response, http.StatusOK)
}

func (h *LeaveHandler) RejectLeave(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	response, err := h.srv.RejectLeave(r.Context(), id)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(w, "Leave rejected successfully", response, http.StatusOK)
}

func (h *LeaveHandler) GetLeaveSummary(w http.ResponseWriter, r *http.Request) {
	summary, err := h.srv.GetLeaveSummary(r.Context())
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(w, "Leave summary retrieved successfully", summary, http.StatusOK)
}
