package handler

import (
	"encoding/json"
	"hrsync-backend/internal/dto"
	"hrsync-backend/internal/service"
	"hrsync-backend/internal/utils"
	"net/http"
)

type TemplateKPIHandler struct {
	srv service.TemplateKPIService
}

func NewTemplateKPIHandler(srv service.TemplateKPIService) *TemplateKPIHandler {
	return &TemplateKPIHandler{srv: srv}
}

func (h *TemplateKPIHandler) GetTemplatesKPI(w http.ResponseWriter, r *http.Request) {
	responses, total, err := h.srv.GetTemplatesKPI(r.Context())
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendPaginated(w, "KPI Templates retrieved successfully", responses, total, 1, 10)
}

func (h *TemplateKPIHandler) GetPublishedKPIByDepartment(w http.ResponseWriter, r *http.Request) {
	department := r.URL.Query().Get("department")
	if department == "" {
		utils.SendError(w, "Department is required", http.StatusBadRequest)
		return
	}

	responses, err := h.srv.GetPublishedTemplatesKPIByDepartment(r.Context(), department)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendSuccess(w, "Published KPI Templates retrieved successfully", responses, http.StatusOK)
}

func (h *TemplateKPIHandler) CreateTemplateKPI(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTemplateKPIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.srv.CreateTemplateKPI(r.Context(), req)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(w, "KPI Template created successfully", response, http.StatusCreated)
}

func (h *TemplateKPIHandler) UpdateTemplateKPI(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req dto.UpdateTemplateKPIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.srv.UpdateTemplateKPI(r.Context(), id, req)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(w, "KPI Template updated successfully", response, http.StatusOK)
}

func (h *TemplateKPIHandler) DeleteTemplateKPI(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.srv.DeleteTemplateKPI(r.Context(), id); err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(w, "KPI Template deleted successfully", nil, http.StatusOK)
}

// Removed saveBase64File as it is now in utils
