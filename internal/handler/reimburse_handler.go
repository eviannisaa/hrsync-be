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

type ReimburseHandler struct {
	srv service.ReimburseService
}

func NewReimburseHandler(srv service.ReimburseService) *ReimburseHandler {
	return &ReimburseHandler{srv: srv}
}

func (h *ReimburseHandler) GetReimbursements(w http.ResponseWriter, r *http.Request) {
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

	responses, total, err := h.srv.GetAll(r.Context(), params)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendPaginated(w, "Reimbursements retrieved successfully", responses, total, page, limit)
}

func (h *ReimburseHandler) CreateReimbursement(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateReimburseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}
	response, err := h.srv.Create(r.Context(), req)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendSuccess(w, "Reimbursement created successfully", response, http.StatusCreated)
}

func (h *ReimburseHandler) UpdateReimbursement(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req dto.UpdateReimburseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}
	response, err := h.srv.Update(r.Context(), id, req)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendSuccess(w, "Reimbursement updated successfully", response, http.StatusOK)
}

func (h *ReimburseHandler) ApproveReimbursement(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	response, err := h.srv.UpdateStatus(r.Context(), id, "APPROVED")
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendSuccess(w, "Reimbursement approved successfully", response, http.StatusOK)
}

func (h *ReimburseHandler) RejectReimbursement(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	response, err := h.srv.UpdateStatus(r.Context(), id, "REJECTED")
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendSuccess(w, "Reimbursement rejected successfully", response, http.StatusOK)
}

func (h *ReimburseHandler) DeleteReimbursement(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.srv.Delete(r.Context(), id); err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendSuccess(w, "Reimbursement deleted successfully", nil, http.StatusOK)
}
