package service

import (
	"context"
	"hrsync-backend/internal/dto"
	"hrsync-backend/internal/model"
	"hrsync-backend/internal/repository"
)

type EmployeeService interface {
	GetEmployees(ctx context.Context, params model.ListParams) ([]dto.EmployeeResponse, int, error)
	GetEmployeeByID(ctx context.Context, id string) (*dto.EmployeeResponse, error)
	CreateEmployee(ctx context.Context, req dto.CreateEmployeeRequest) (*dto.EmployeeResponse, error)
	UpdateEmployee(ctx context.Context, id string, req dto.UpdateEmployeeRequest) (*dto.EmployeeResponse, error)
	DeleteEmployee(ctx context.Context, id string) error
	GetOrganization(ctx context.Context) (*dto.EmployeeOrganizationResponse, error)
	UpdateOrganization(ctx context.Context, req dto.UpdateEmployeeOrganizationRequest) (*dto.EmployeeOrganizationResponse, error)
}

type employeeService struct {
	repo repository.EmployeeRepository
}

func NewEmployeeService(repo repository.EmployeeRepository) EmployeeService {
	return &employeeService{repo: repo}
}

func (s *employeeService) GetEmployees(ctx context.Context, params model.ListParams) ([]dto.EmployeeResponse, int, error) {
	return s.repo.GetAll(ctx, params)
}

func (s *employeeService) GetEmployeeByID(ctx context.Context, id string) (*dto.EmployeeResponse, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *employeeService) CreateEmployee(ctx context.Context, req dto.CreateEmployeeRequest) (*dto.EmployeeResponse, error) {
	return s.repo.Create(ctx, req)
}

func (s *employeeService) UpdateEmployee(ctx context.Context, id string, req dto.UpdateEmployeeRequest) (*dto.EmployeeResponse, error) {
	return s.repo.Update(ctx, id, req)
}

func (s *employeeService) DeleteEmployee(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *employeeService) GetOrganization(ctx context.Context) (*dto.EmployeeOrganizationResponse, error) {
	return s.repo.GetOrganization(ctx)
}

func (s *employeeService) UpdateOrganization(ctx context.Context, req dto.UpdateEmployeeOrganizationRequest) (*dto.EmployeeOrganizationResponse, error) {
	return s.repo.UpdateOrganization(ctx, req)
}
