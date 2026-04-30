package service

import (
	"context"
	"hrsync-backend/internal/dto"
	"hrsync-backend/internal/repository"
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error)
	GeneratePassword(ctx context.Context, employeeId string) (string, error)
	HandleGoogleAuth(ctx context.Context, email string) (*dto.AuthResponse, error)
}

type authService struct {
	repo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) AuthService {
	return &authService{repo: repo}
}

func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	return s.repo.Register(ctx, req)
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	return s.repo.Login(ctx, req)
}

func (s *authService) GeneratePassword(ctx context.Context, employeeId string) (string, error) {
	return s.repo.GeneratePassword(ctx, employeeId)
}

func (s *authService) HandleGoogleAuth(ctx context.Context, email string) (*dto.AuthResponse, error) {
	return s.repo.HandleGoogleAuth(ctx, email)
}
