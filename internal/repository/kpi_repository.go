package repository

import (
	"context"
	"hrsync-backend/internal/db"
	"hrsync-backend/internal/dto"
)

type TemplateKPIRepository interface {
	GetAll(ctx context.Context) ([]dto.TemplateKPIResponse, int, error)
	GetPublishedByDepartment(ctx context.Context, department string) ([]dto.TemplateKPIResponse, error)
	Create(ctx context.Context, req dto.CreateTemplateKPIRequest) (*dto.TemplateKPIResponse, error)
	Update(ctx context.Context, id string, req dto.UpdateTemplateKPIRequest) (*dto.TemplateKPIResponse, error)
	Delete(ctx context.Context, id string) error
}

type templateKPIRepository struct {
	client *db.PrismaClient
}

func NewTemplateKPIRepository(client *db.PrismaClient) TemplateKPIRepository {
	return &templateKPIRepository{client: client}
}

func (r *templateKPIRepository) GetAll(ctx context.Context) ([]dto.TemplateKPIResponse, int, error) {
	allKPI, err := r.client.TemplateKPI.FindMany().Exec(ctx)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.TemplateKPIResponse, 0, len(allKPI))
	for _, du := range allKPI {
		responses = append(responses, dto.TemplateKPIResponse{
			InnerTemplateKPI: du.InnerTemplateKPI,
		})
	}

	return responses, len(allKPI), nil
}

func (r *templateKPIRepository) GetPublishedByDepartment(ctx context.Context, department string) ([]dto.TemplateKPIResponse, error) {
	allKPI, err := r.client.TemplateKPI.FindMany(
		db.TemplateKPI.Department.Equals(department),
		db.TemplateKPI.IsPublished.Equals(true),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.TemplateKPIResponse, 0, len(allKPI))
	for _, du := range allKPI {
		responses = append(responses, dto.TemplateKPIResponse{
			InnerTemplateKPI: du.InnerTemplateKPI,
		})
	}

	return responses, nil
}

func (r *templateKPIRepository) Create(ctx context.Context, req dto.CreateTemplateKPIRequest) (*dto.TemplateKPIResponse, error) {
	du, err := r.client.TemplateKPI.CreateOne(
		db.TemplateKPI.Email.Set(req.Email),
		db.TemplateKPI.Department.Set(req.Department),
		db.TemplateKPI.TemplateName.Set(req.TemplateName),
		db.TemplateKPI.Description.Set(req.Description),
		db.TemplateKPI.Attachment.Set(req.Attachment),
		db.TemplateKPI.IsPublished.Set(req.IsPublished),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &dto.TemplateKPIResponse{InnerTemplateKPI: du.InnerTemplateKPI}, nil
}

func (r *templateKPIRepository) Update(ctx context.Context, id string, req dto.UpdateTemplateKPIRequest) (*dto.TemplateKPIResponse, error) {
	du, err := r.client.TemplateKPI.FindUnique(
		db.TemplateKPI.ID.Equals(id),
	).Update(
		db.TemplateKPI.Department.Set(req.Department),
		db.TemplateKPI.TemplateName.Set(req.TemplateName),
		db.TemplateKPI.Description.Set(req.Description),
		db.TemplateKPI.Attachment.Set(req.Attachment),
		db.TemplateKPI.IsPublished.Set(req.IsPublished),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &dto.TemplateKPIResponse{InnerTemplateKPI: du.InnerTemplateKPI}, nil
}

func (r *templateKPIRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.TemplateKPI.FindUnique(
		db.TemplateKPI.ID.Equals(id),
	).Delete().Exec(ctx)
	return err
}
