package handler

import (
	"hrsync-backend/internal/dto"
	"hrsync-backend/internal/model"
	"hrsync-backend/internal/service"
	"hrsync-backend/internal/utils"
	"io"
	"net/http"
)

type PayslipHandler struct {
	srv service.PayslipService
}

func NewPayslipHandler(srv service.PayslipService) *PayslipHandler {
	return &PayslipHandler{srv: srv}
}

func (h *PayslipHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	
	role, _ := r.Context().Value(model.ContextKeyRole).(string)
	userEmail, _ := r.Context().Value(model.ContextKeyEmail).(string)

	// Security: If user is EMPLOYEE, override email param to their own email
	if role == "EMPLOYEE" {
		email = userEmail
	}
	
	payslips, err := h.srv.GetPayslips(r.Context(), email)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(w, "Payslips retrieved successfully", payslips, http.StatusOK)
}

func (h *PayslipHandler) BulkUpload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(50 << 20); err != nil { // 50MB limit
		utils.SendError(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	month := r.FormValue("month")
	year := r.FormValue("year")

	if month == "" || year == "" {
		utils.SendError(w, "Month and year are required", http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		utils.SendError(w, "No files uploaded", http.StatusBadRequest)
		return
	}

	uploadItems := make([]dto.PayslipUploadItem, 0, len(files))
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			continue
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			continue
		}

		uploadItems = append(uploadItems, dto.PayslipUploadItem{
			Filename:    fileHeader.Filename,
			Content:     content,
			Size:        fileHeader.Size,
			ContentType: fileHeader.Header.Get("Content-Type"),
		})
	}

	createdBy, _ := r.Context().Value(model.ContextKeyEmail).(string)
	if createdBy == "" {
		createdBy = "Administrator"
	}

	result, err := h.srv.BulkUpload(r.Context(), uploadItems, month, year, createdBy)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(w, "Bulk upload completed", result, http.StatusOK)
}

func (h *PayslipHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		utils.SendError(w, "ID is required", http.StatusBadRequest)
		return
	}

	err := h.srv.Delete(r.Context(), id)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(w, "Payslip deleted successfully", nil, http.StatusOK)
}

func (h *PayslipHandler) DeleteBatch(w http.ResponseWriter, r *http.Request) {
	month := r.URL.Query().Get("month")
	year := r.URL.Query().Get("year")

	if month == "" || year == "" {
		utils.SendError(w, "Month and year are required", http.StatusBadRequest)
		return
	}

	err := h.srv.DeleteBatch(r.Context(), month, year)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(w, "Batch payslips deleted successfully", nil, http.StatusOK)
}
