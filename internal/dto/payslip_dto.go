package dto

import (
	"hrsync-backend/internal/db"
)

type PayslipResponse struct {
	db.InnerPayslip
	EmployeeName string `json:"employeeName"`
	Department   string `json:"department"`
	CreatedBy    string `json:"createdBy"`
}

type BulkUploadPayslipResponse struct {
	Total    int      `json:"total"`
	Success  int      `json:"success"`
	Failed   int      `json:"failed"`
	Messages []string `json:"messages"`
}

type PayslipUploadItem struct {
	Filename    string
	Content     []byte
	Size        int64
	ContentType string
}
