package repository

import (
	"context"
	"hrsync-backend/internal/db"
	"hrsync-backend/internal/dto"
	"log"
)

type PayslipRepository interface {
	GetAll(ctx context.Context, email string) ([]dto.PayslipResponse, error)
	Create(ctx context.Context, email, fileUrl, month, year, createdBy string) (*dto.PayslipResponse, error)
	Delete(ctx context.Context, id string) error
	DeleteBatch(ctx context.Context, month, year string) error
}

type payslipRepository struct {
	client *db.PrismaClient
}

func NewPayslipRepository(client *db.PrismaClient) PayslipRepository {
	return &payslipRepository{client: client}
}

func (r *payslipRepository) GetAll(ctx context.Context, email string) ([]dto.PayslipResponse, error) {
	var filters []db.PayslipWhereParam
	if email != "" {
		filters = append(filters, db.Payslip.Email.Equals(email))
	}

	dbPayslips, err := r.client.Payslip.FindMany(filters...).With(db.Payslip.Employee.Fetch()).OrderBy(db.Payslip.CreatedAt.Order(db.SortOrderDesc)).Exec(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.PayslipResponse, 0, len(dbPayslips))
	for _, p := range dbPayslips {
		createdBy, _ := p.CreatedBy()
		res := dto.PayslipResponse{
			InnerPayslip: p.InnerPayslip,
			CreatedBy:    createdBy,
		}
		if emp := p.Employee(); emp != nil {
			res.EmployeeName = emp.Name
			res.Department = emp.Department
		}
		responses = append(responses, res)
	}

	return responses, nil
}

func (r *payslipRepository) Create(ctx context.Context, email, fileUrl, month, year, createdBy string) (*dto.PayslipResponse, error) {
	p, err := r.client.Payslip.CreateOne(
		db.Payslip.FileURL.Set(fileUrl),
		db.Payslip.Month.Set(month),
		db.Payslip.Year.Set(year),
		db.Payslip.Employee.Link(db.Employee.Email.Equals(email)),
		db.Payslip.CreatedBy.Set(createdBy),
	).With(db.Payslip.Employee.Fetch()).Exec(ctx)

	if err != nil {
		log.Printf("Error creating payslip: %v", err)
		return nil, err
	}

	sender, _ := p.CreatedBy()
	res := dto.PayslipResponse{
		InnerPayslip: p.InnerPayslip,
		CreatedBy:    sender,
	}
	if emp := p.Employee(); emp != nil {
		res.EmployeeName = emp.Name
		res.Department = emp.Department
	}

	return &res, nil
}

func (r *payslipRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.Payslip.FindUnique(db.Payslip.ID.Equals(id)).Delete().Exec(ctx)
	return err
}

func (r *payslipRepository) DeleteBatch(ctx context.Context, month, year string) error {
	_, err := r.client.Payslip.FindMany(
		db.Payslip.Month.Equals(month),
		db.Payslip.Year.Equals(year),
	).Delete().Exec(ctx)
	return err
}
