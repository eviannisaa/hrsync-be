package dto

import (
	"hrsync-backend/internal/db"
)

type TemplateKPIResponse struct {
	db.InnerTemplateKPI
}

type CreateTemplateKPIRequest struct {
	Email        string `json:"email"`
	Department   string `json:"department"`
	TemplateName string `json:"templateName"`
	Description  string `json:"description"`
	Attachment   string `json:"attachment"`
	IsPublished  bool   `json:"isPublished"`
}

type UpdateTemplateKPIRequest struct {
	Department   string `json:"department"`
	TemplateName string `json:"templateName"`
	Description  string `json:"description"`
	Attachment   string `json:"attachment"`
	IsPublished  bool   `json:"isPublished"`
}
