package dto

import (
	"hrsync-backend/internal/db"
)

type KPIItemDTO struct {
	ID         string `json:"id"`
	TemplateID string `json:"templateId"`
	NameResult string `json:"nameResult"`
	KpiResult  string `json:"kpiResult"`
	Weight     float64 `json:"weight"`
	Target     float64 `json:"target"`
	Actual     float64 `json:"actual"`
	Score      float64 `json:"score"`
	FinalScore float64 `json:"finalScore"`
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
