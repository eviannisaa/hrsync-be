package handler

import (
	"hrsync-backend/internal/utils"
	"net/http"
)

type DepartmentHandler struct{}

func NewDepartmentHandler() *DepartmentHandler {
	return &DepartmentHandler{}
}

func (h *DepartmentHandler) GetDepartments(w http.ResponseWriter, r *http.Request) {
	departments := []string{"engineering", "finance", "ensol", "hrd"}
	utils.SendSuccess(w, "Departments retrieved successfully", departments, http.StatusOK)
}
