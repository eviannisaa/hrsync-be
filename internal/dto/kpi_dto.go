package dto

import (
	"hrsync-backend/internal/db"
)

type KPIItemDTO struct {
	db.InnerKPIItem
}

type TemplateKPIResponse struct {
	db.InnerTemplateKPI
	Items []KPIItemDTO `json:"items"`
}

type CreateTemplateKPIRequest struct {
	Email        string       `json:"email"`
	Department   string       `json:"department"`
	TemplateName string       `json:"templateName"`
	Description  string       `json:"description"`
	Attachment   string       `json:"attachment"`
	IsPublished  bool         `json:"isPublished"`
	Items        []KPIItemDTO `json:"items"`
}

type UpdateTemplateKPIRequest struct {
	Department   string       `json:"department"`
	TemplateName string       `json:"templateName"`
	Description  string       `json:"description"`
	Attachment   *string      `json:"attachment"`
	IsPublished  *bool        `json:"isPublished"`
	Items        []KPIItemDTO `json:"items"`
}
