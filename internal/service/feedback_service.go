package service

import (
	"context"
	"hrsync-backend/internal/dto"
	"hrsync-backend/internal/model"
	"hrsync-backend/internal/repository"
)

type FeedbackService interface {
	GetAll(ctx context.Context, params model.ListParams) ([]dto.FeedbackResponse, int, error)
	Create(ctx context.Context, req dto.CreateFeedbackRequest) (*dto.FeedbackResponse, error)
	Delete(ctx context.Context, id string) error
}

type feedbackService struct {
	repo repository.FeedbackRepository
}

func NewFeedbackService(repo repository.FeedbackRepository) FeedbackService {
	return &feedbackService{repo: repo}
}

func (s *feedbackService) GetAll(ctx context.Context, params model.ListParams) ([]dto.FeedbackResponse, int, error) {
	return s.repo.GetAll(ctx, params)
}

func (s *feedbackService) Create(ctx context.Context, req dto.CreateFeedbackRequest) (*dto.FeedbackResponse, error) {
	return s.repo.Create(ctx, req)
}

func (s *feedbackService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
