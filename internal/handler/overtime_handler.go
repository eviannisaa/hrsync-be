package handler

import (
	"encoding/json"
	"hrsync-backend/internal/dto"
	"hrsync-backend/internal/model"
	"hrsync-backend/internal/service"
	"hrsync-backend/internal/utils"
	"net/http"
	"strconv"
)

type OvertimeHandler struct {
	srv service.OvertimeService
}

func NewOvertimeHandler(srv service.OvertimeService) *OvertimeHandler {
	return &OvertimeHandler{srv: srv}
}

func (h *OvertimeHandler) GetOvertimes(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	params := model.ListParams{
		Page:    page,
		Limit:   limit,
		Search:  r.URL.Query().Get("search"),
		SortBy:  r.URL.Query().Get("sortBy"),
		SortDir: r.URL.Query().Get("sortDir"),
	}

	responses, total, err := h.srv.GetOvertimes(r.Context(), params)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendPaginated(w, "Overtime retrieved successfully", responses, total, page, limit)
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
