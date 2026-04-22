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

type EmployeeHandler struct {
	srv service.EmployeeService
}

func NewEmployeeHandler(srv service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{srv: srv}
}

func (h *EmployeeHandler) GetEmployees(w http.ResponseWriter, r *http.Request) {
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

	responses, total, err := h.srv.GetEmployees(r.Context(), params)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendPaginated(w, "Employees retrieved successfully", responses, total, page, limit)
}

func (h *EmployeeHandler) GetEmployeeByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	response, err := h.srv.GetEmployeeByID(r.Context(), id)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusNotFound)
		return
	}
	utils.SendSuccess(w, "Employee retrieved successfully", response, http.StatusOK)
}

func (h *EmployeeHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateEmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.srv.CreateEmployee(r.Context(), req)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(w, "Employee created successfully", response, http.StatusCreated)
}

func (h *EmployeeHandler) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req dto.UpdateEmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.srv.UpdateEmployee(r.Context(), id, req)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(w, "Employee updated successfully", response, http.StatusOK)
}

func (h *EmployeeHandler) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.srv.DeleteEmployee(r.Context(), id); err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(w, "Employee deleted successfully", nil, http.StatusOK)
}

func (h *EmployeeHandler) GetOrganization(w http.ResponseWriter, r *http.Request) {
	response, err := h.srv.GetOrganization(r.Context())
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendSuccess(w, "Organization retrieved successfully", response, http.StatusOK)
}

func (h *EmployeeHandler) UpdateOrganization(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateEmployeeOrganizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.srv.UpdateOrganization(r.Context(), req)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendSuccess(w, "Organization updated successfully", response, http.StatusOK)
}
