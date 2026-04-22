package service

import (
	"context"
	"hrsync-backend/internal/dto"
	"hrsync-backend/internal/model"
	"hrsync-backend/internal/repository"
)

type ReimburseService interface {
	GetAll(ctx context.Context, params model.ListParams) ([]dto.ReimburseResponse, int, error)
	Create(ctx context.Context, req dto.CreateReimburseRequest) (*dto.ReimburseResponse, error)
	Update(ctx context.Context, id string, req dto.UpdateReimburseRequest) (*dto.ReimburseResponse, error)
	UpdateStatus(ctx context.Context, id string, status string) (*dto.ReimburseResponse, error)
	Delete(ctx context.Context, id string) error
}

type reimburseService struct {
	repo repository.ReimburseRepository
}

func NewReimburseService(repo repository.ReimburseRepository) ReimburseService {
	return &reimburseService{repo: repo}
}

func (s *reimburseService) GetAll(ctx context.Context, params model.ListParams) ([]dto.ReimburseResponse, int, error) {
	return s.repo.GetAll(ctx, params)
}

func (s *reimburseService) Create(ctx context.Context, req dto.CreateReimburseRequest) (*dto.ReimburseResponse, error) {
	return s.repo.Create(ctx, req)
}

func (s *reimburseService) Update(ctx context.Context, id string, req dto.UpdateReimburseRequest) (*dto.ReimburseResponse, error) {
	return s.repo.Update(ctx, id, req)
}

func (s *reimburseService) UpdateStatus(ctx context.Context, id string, status string) (*dto.ReimburseResponse, error) {
	return s.repo.UpdateStatus(ctx, id, status)
}

func (s *reimburseService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
