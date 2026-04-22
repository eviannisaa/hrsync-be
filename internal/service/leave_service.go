package service

import (
	"context"
	"hrsync-backend/internal/dto"
	"hrsync-backend/internal/model"
	"hrsync-backend/internal/repository"
)

type LeaveService interface {
	GetLeaves(ctx context.Context, params model.ListParams) ([]dto.LeaveResponse, int, error)
	CreateLeaves(ctx context.Context, req dto.CreateLeaveRequest) (*dto.LeaveResponse, error)
	UpdateLeaves(ctx context.Context, id string, req dto.UpdateLeaveRequest) (*dto.LeaveResponse, error)
	ApproveLeave(ctx context.Context, id string) (*dto.LeaveResponse, error)
	RejectLeave(ctx context.Context, id string) (*dto.LeaveResponse, error)
	DeleteLeaves(ctx context.Context, id string) error
	GetLeaveSummary(ctx context.Context) ([]dto.LeaveStatusSummary, error)
}

type leaveService struct {
	repo repository.LeaveRepository
}

func NewLeaveService(repo repository.LeaveRepository) LeaveService {
	return &leaveService{repo: repo}
}

func (s *leaveService) GetLeaves(ctx context.Context, params model.ListParams) ([]dto.LeaveResponse, int, error) {
	_ = s.repo.AutoUpdateOngoing(ctx)
	return s.repo.GetAll(ctx, params)
}

func (s *leaveService) CreateLeaves(ctx context.Context, req dto.CreateLeaveRequest) (*dto.LeaveResponse, error) {
	return s.repo.Create(ctx, req)
}

func (s *leaveService) UpdateLeaves(ctx context.Context, id string, req dto.UpdateLeaveRequest) (*dto.LeaveResponse, error) {
	return s.repo.Update(ctx, id, req)
}

func (s *leaveService) DeleteLeaves(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *leaveService) ApproveLeave(ctx context.Context, id string) (*dto.LeaveResponse, error) {
	return s.repo.UpdateStatus(ctx, id, "APPROVED")
}

func (s *leaveService) RejectLeave(ctx context.Context, id string) (*dto.LeaveResponse, error) {
	return s.repo.UpdateStatus(ctx, id, "REJECTED")
}

func (s *leaveService) GetLeaveSummary(ctx context.Context) ([]dto.LeaveStatusSummary, error) {
	_ = s.repo.AutoUpdateOngoing(ctx)
	return s.repo.GetSummary(ctx)
}
