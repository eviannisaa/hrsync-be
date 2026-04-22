package handler

import (
	"encoding/json"
	"hrsync-backend/internal/dto"
	"hrsync-backend/internal/middleware"
	"hrsync-backend/internal/model"
	"hrsync-backend/internal/service"
	"hrsync-backend/internal/utils"
	"net/http"
	"strconv"
)

type LeaveHandler struct {
	srv service.LeaveService
}

func NewLeaveHandler(srv service.LeaveService) *LeaveHandler {
	return &LeaveHandler{srv: srv}
}

func (h *LeaveHandler) GetLeaves(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	email, _ := r.Context().Value(middleware.ContextKeyEmail).(string)
	role, _ := r.Context().Value(middleware.ContextKeyRole).(string)

	params := model.ListParams{
		Page:    page,
		Limit:   limit,
		Search:  r.URL.Query().Get("search"),
		SortBy:  r.URL.Query().Get("sortBy"),
		SortDir: r.URL.Query().Get("sortDir"),
		Role:    role,
	}

	// If user is non-admin, filter by their email. 
	// If admin, they see all unless they specifically ask for their own via ?mine=true
	isMine := r.URL.Query().Get("mine") == "true"
	if role != "ADMIN" || isMine {
		params.Email = email
	}

	responses, total, err := h.srv.GetLeaves(r.Context(), params)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendPaginated(w, "Leaves retrieved successfully", responses, total, page, limit)
}

func (h *LeaveHandler) CreateLeave(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateLeaveRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	email, _ := r.Context().Value(middleware.ContextKeyEmail).(string)
	role, _ := r.Context().Value(middleware.ContextKeyRole).(string)

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

	email, _ := r.Context().Value(middleware.ContextKeyEmail).(string)
	role, _ := r.Context().Value(middleware.ContextKeyRole).(string)

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
