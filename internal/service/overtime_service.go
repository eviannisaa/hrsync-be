package service

import (
	"context"
	"hrsync-backend/internal/dto"
	"hrsync-backend/internal/model"
	"hrsync-backend/internal/repository"
)

type OvertimeService interface {
	GetOvertimes(ctx context.Context, params model.ListParams) ([]dto.OvertimeResponse, int, error)
	CreateOvertimes(ctx context.Context, req dto.CreateOvertimeRequest) (*dto.OvertimeResponse, error)
	UpdateOvertimes(ctx context.Context, id string, req dto.UpdateOvertimeRequest) (*dto.OvertimeResponse, error)
	ApproveOvertimes(ctx context.Context, id string) (*dto.OvertimeResponse, error)
	RejectOvertimes(ctx context.Context, id string) (*dto.OvertimeResponse, error)
	DeleteOvertimes(ctx context.Context, id string) error
}

type overtimeService struct {
	repo repository.OvertimeRepository
}

func NewOvertimeService(repo repository.OvertimeRepository) OvertimeService {
	return &overtimeService{repo: repo}
}

func (s *overtimeService) GetOvertimes(ctx context.Context, params model.ListParams) ([]dto.OvertimeResponse, int, error) {
	return s.repo.GetAll(ctx, params)
}

func (s *overtimeService) CreateOvertimes(ctx context.Context, req dto.CreateOvertimeRequest) (*dto.OvertimeResponse, error) {
	return s.repo.Create(ctx, req)
}

func (s *overtimeService) UpdateOvertimes(ctx context.Context, id string, req dto.UpdateOvertimeRequest) (*dto.OvertimeResponse, error) {
	return s.repo.Update(ctx, id, req)
}

func (s *overtimeService) DeleteOvertimes(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *overtimeService) ApproveOvertimes(ctx context.Context, id string) (*dto.OvertimeResponse, error) {
	return s.repo.UpdateStatus(ctx, id, "APPROVED")
}

func (s *overtimeService) RejectOvertimes(ctx context.Context, id string) (*dto.OvertimeResponse, error) {
	return s.repo.UpdateStatus(ctx, id, "REJECTED")
}
