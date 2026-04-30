package service

import (
	"bytes"
	"context"
	"hrsync-backend/internal/dto"
	"hrsync-backend/internal/model"
	"hrsync-backend/internal/repository"
	"hrsync-backend/internal/utils"
	"strings"
)

type PayslipService interface {
	GetPayslips(ctx context.Context, email string) ([]dto.PayslipResponse, error)
	BulkUpload(ctx context.Context, items []dto.PayslipUploadItem, month, year, createdBy string) (*dto.BulkUploadPayslipResponse, error)
	Delete(ctx context.Context, id string) error
	DeleteBatch(ctx context.Context, month, year string) error
}

type payslipService struct {
	repo         repository.PayslipRepository
	employeeRepo repository.EmployeeRepository
}

func NewPayslipService(repo repository.PayslipRepository, employeeRepo repository.EmployeeRepository) PayslipService {
	return &payslipService{
		repo:         repo,
		employeeRepo: employeeRepo,
	}
}

func (s *payslipService) GetPayslips(ctx context.Context, email string) ([]dto.PayslipResponse, error) {
	return s.repo.GetAll(ctx, email)
}

func (s *payslipService) BulkUpload(ctx context.Context, items []dto.PayslipUploadItem, month, year, createdBy string) (*dto.BulkUploadPayslipResponse, error) {
	resp := &dto.BulkUploadPayslipResponse{
		Total:    len(items),
		Messages: []string{},
	}

	// 1. Get all employees for matching
	employees, _, err := s.employeeRepo.GetAll(ctx, model.ListParams{Limit: 1000, Page: 1})
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		// Clean filename to match (e.g. "John Doe.pdf" -> "john doe")
		cleanName := strings.TrimSuffix(strings.ToLower(item.Filename), ".pdf")
		cleanName = strings.TrimSpace(cleanName)

		var matchedEmail string
		for _, emp := range employees {
			if strings.ToLower(emp.Name) == cleanName {
				matchedEmail = emp.Email
				break
			}
		}

		if matchedEmail == "" {
			resp.Failed++
			resp.Messages = append(resp.Messages, "Failed: Could not match employee for file "+item.Filename)
			continue
		}

		// 2. Upload to MinIO
		reader := bytes.NewReader(item.Content)
		key, err := utils.Upload(ctx, reader, item.Size, "payslips/"+item.Filename, item.ContentType)
		if err != nil {
			resp.Failed++
			resp.Messages = append(resp.Messages, "Failed: MinIO upload error for "+item.Filename+": "+err.Error())
			continue
		}
		fileUrl := utils.GetURL(key)

		// 3. Save to DB
		_, err = s.repo.Create(ctx, matchedEmail, fileUrl, month, year, createdBy)
		if err != nil {
			resp.Failed++
			resp.Messages = append(resp.Messages, "Failed: Could not save record for "+item.Filename+": "+err.Error())
			continue
		}

		resp.Success++
		resp.Messages = append(resp.Messages, "Success: Uploaded for "+cleanName)
	}

	return resp, nil
}

func (s *payslipService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *payslipService) DeleteBatch(ctx context.Context, month, year string) error {
	return s.repo.DeleteBatch(ctx, month, year)
}
