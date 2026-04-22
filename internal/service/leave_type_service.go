package service

import (
	"context"
	"hrsync-backend/internal/dto"
	"hrsync-backend/internal/repository"
)

type LeaveTypeService interface {
	GetLeaveTypes(ctx context.Context) ([]dto.LeaveTypeResponse, error)
	GetMyCredits(ctx context.Context, email string) ([]dto.LeaveBalanceResponse, error)
	SeedLeaveTypes(ctx context.Context) error
}

type leaveTypeService struct {
	repo repository.LeaveTypeRepository
}

func NewLeaveTypeService(repo repository.LeaveTypeRepository) LeaveTypeService {
	return &leaveTypeService{repo: repo}
}

func (s *leaveTypeService) GetLeaveTypes(ctx context.Context) ([]dto.LeaveTypeResponse, error) {
	return s.repo.GetAll(ctx)
}

func (s *leaveTypeService) GetMyCredits(ctx context.Context, email string) ([]dto.LeaveBalanceResponse, error) {
	return s.repo.GetMyCredits(ctx, email)
}

func (s *leaveTypeService) SeedLeaveTypes(ctx context.Context) error {
	return s.repo.Seed(ctx)
}
