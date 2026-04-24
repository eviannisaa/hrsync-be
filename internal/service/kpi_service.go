package service

import (
	"context"
	"fmt"
	"hrsync-backend/internal/dto"
	"hrsync-backend/internal/model"
	"hrsync-backend/internal/repository"
	"hrsync-backend/internal/utils"
	"strings"
	"time"
)

type TemplateKPIService interface {
	GetTemplatesKPI(ctx context.Context, params model.ListParams) ([]dto.TemplateKPIResponse, int, error)
	GetPublishedTemplatesKPIByDepartment(ctx context.Context, department string) ([]dto.TemplateKPIResponse, error)
	CreateTemplateKPI(ctx context.Context, req dto.CreateTemplateKPIRequest) (*dto.TemplateKPIResponse, error)
	UpdateTemplateKPI(ctx context.Context, id string, req dto.UpdateTemplateKPIRequest) (*dto.TemplateKPIResponse, error)
	DeleteTemplateKPI(ctx context.Context, id string) error
}

type templateKPIService struct {
	repo repository.TemplateKPIRepository
}

func NewTemplateKPIService(repo repository.TemplateKPIRepository) TemplateKPIService {
	return &templateKPIService{repo: repo}
}

func (s *templateKPIService) GetTemplatesKPI(ctx context.Context, params model.ListParams) ([]dto.TemplateKPIResponse, int, error) {
	templates, total, err := s.repo.GetAll(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	for i := range templates {
		if templates[i].Attachment != nil && *templates[i].Attachment != "" {
			url := utils.GetURL(*templates[i].Attachment)
			templates[i].Attachment = &url
		}
	}

	return templates, total, nil
}

func (s *templateKPIService) GetPublishedTemplatesKPIByDepartment(ctx context.Context, department string) ([]dto.TemplateKPIResponse, error) {
	templates, err := s.repo.GetPublishedByDepartment(ctx, department)
	if err != nil {
		return nil, err
	}

	for i := range templates {
		if templates[i].Attachment != nil && *templates[i].Attachment != "" {
			url := utils.GetURL(*templates[i].Attachment)
			templates[i].Attachment = &url
		}
	}

	return templates, nil
}

func (s *templateKPIService) CreateTemplateKPI(ctx context.Context, req dto.CreateTemplateKPIRequest) (*dto.TemplateKPIResponse, error) {
	// Handle Base64 upload if present
	if req.Attachment != "" && strings.Contains(req.Attachment, ";base64,") {
		objectName := fmt.Sprintf("kpi/%d-attachment", time.Now().UnixNano())
		path, err := utils.UploadBase64(ctx, req.Attachment, objectName)
		if err != nil {
			return nil, fmt.Errorf("storage upload failed: %w", err)
		}
		req.Attachment = path
	}

	return s.repo.Create(ctx, req)
}

func (s *templateKPIService) UpdateTemplateKPI(ctx context.Context, id string, req dto.UpdateTemplateKPIRequest) (*dto.TemplateKPIResponse, error) {
	// Handle Base64 upload if present
	if req.Attachment != nil && *req.Attachment != "" && strings.Contains(*req.Attachment, ";base64,") {
		objectName := fmt.Sprintf("kpi/%d-attachment", time.Now().UnixNano())
		path, err := utils.UploadBase64(ctx, *req.Attachment, objectName)
		if err != nil {
			return nil, fmt.Errorf("storage upload failed: %w", err)
		}
		req.Attachment = &path
	}

	return s.repo.Update(ctx, id, req)
}

func (s *templateKPIService) DeleteTemplateKPI(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
