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

type FeedbackHandler struct {
	srv service.FeedbackService
}

func NewFeedbackHandler(srv service.FeedbackService) *FeedbackHandler {
	return &FeedbackHandler{srv: srv}
}

func (h *FeedbackHandler) GetFeedbacks(w http.ResponseWriter, r *http.Request) {
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
	utils.SendPaginated(w, "Feedbacks retrieved successfully", responses, total, page, limit)
}

func (h *FeedbackHandler) CreateFeedback(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateFeedbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.srv.Create(r.Context(), req)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(w, "Feedback created successfully", response, http.StatusCreated)
}

func (h *FeedbackHandler) DeleteFeedback(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.srv.Delete(r.Context(), id); err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(w, "Feedback deleted successfully", nil, http.StatusOK)
}
