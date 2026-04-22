package repository

import (
	"context"
	"hrsync-backend/internal/db"
	"hrsync-backend/internal/dto"
	"hrsync-backend/internal/model"
	"log"
)

type ReimburseRepository interface {
	GetAll(ctx context.Context, params model.ListParams) ([]dto.ReimburseResponse, int, error)
	Create(ctx context.Context, req dto.CreateReimburseRequest) (*dto.ReimburseResponse, error)
	Update(ctx context.Context, id string, req dto.UpdateReimburseRequest) (*dto.ReimburseResponse, error)
	UpdateStatus(ctx context.Context, id string, status string) (*dto.ReimburseResponse, error)
	Delete(ctx context.Context, id string) error
}

type reimburseRepository struct {
	client *db.PrismaClient
}

func NewReimburseRepository(client *db.PrismaClient) ReimburseRepository {
	return &reimburseRepository{client: client}
}

func (r *reimburseRepository) GetAll(ctx context.Context, params model.ListParams) ([]dto.ReimburseResponse, int, error) {
	skip := (params.Page - 1) * params.Limit

	// Build filter
	var filters []db.ReimburseWhereParam
	if params.Search != "" {
		var rawRes []struct {
			ID string `json:"id"`
		}
		// Use ILIKE for case-insensitive search across all relevant fields
		err := r.client.Prisma.QueryRaw(`
			SELECT r.id FROM "Reimburse" r
			LEFT JOIN "Employee" e ON r."email" = e.email
			WHERE r.status ILIKE $1 
			   OR r.description ILIKE $1 
			   OR e.name ILIKE $1 
			   OR e.department ILIKE $1 
			   OR e.position ILIKE $1
		`, "%"+params.Search+"%").Exec(ctx, &rawRes)

		if err == nil {
			ids := make([]string, len(rawRes))
			for i, r := range rawRes {
				ids[i] = r.ID
			}
			filters = append(filters, db.Reimburse.ID.In(ids))
		}
	}

	// Build sort
	sortDir := db.SortOrderDesc
	if params.SortDir == "asc" {
		sortDir = db.SortOrderAsc
	}
	var orderBy []db.ReimburseOrderByParam
	switch params.SortBy {
	case "createdAt", "date":
		orderBy = append(orderBy, db.Reimburse.CreatedAt.Order(sortDir))
	case "status":
		orderBy = append(orderBy, db.Reimburse.Status.Order(sortDir))
	case "amount":
		orderBy = append(orderBy, db.Reimburse.Amount.Order(sortDir))
	case "description":
		orderBy = append(orderBy, db.Reimburse.Description.Order(sortDir))
	default:
		orderBy = append(orderBy, db.Reimburse.CreatedAt.Order(sortDir))
	}

	dbReimburse, err := r.client.Reimburse.FindMany(filters...).With(db.Reimburse.Employee.Fetch()).OrderBy(orderBy...).Skip(skip).Take(params.Limit).Exec(ctx)
	if err != nil {
		log.Printf("Error in ReimburseRepository.GetAll (FindMany): %v", err)
		return nil, 0, err
	}

	allReimburse, err := r.client.Reimburse.FindMany(filters...).Exec(ctx)
	if err != nil {
		log.Printf("Error in ReimburseRepository.GetAll (Count/FindMany): %v", err)
		return nil, 0, err
	}
	total := len(allReimburse)

	responses := make([]dto.ReimburseResponse, 0, len(dbReimburse))
	for _, du := range dbReimburse {
		res := dto.ReimburseResponse{
			InnerReimburse: du.InnerReimburse,
		}
		if emp := du.Employee(); emp != nil {
			res.EmployeeName = emp.Name
			res.Department = emp.Department
		}
		responses = append(responses, res)
	}

	return responses, total, nil
}

func (r *reimburseRepository) Create(ctx context.Context, req dto.CreateReimburseRequest) (*dto.ReimburseResponse, error) {
	optionalParams := []db.ReimburseSetParam{}
	if req.CreatorRole != nil {
		optionalParams = append(optionalParams, db.Reimburse.CreatorRole.Set(*req.CreatorRole))
	}

	du, err := r.client.Reimburse.CreateOne(
		db.Reimburse.Amount.Set(req.Amount),
		db.Reimburse.Description.Set(req.Description),
		db.Reimburse.AttachBill.Set(req.AttachBill),
		db.Reimburse.PaymentProof.Set(req.PaymentProof),
		db.Reimburse.Employee.Link(db.Employee.Email.Equals(req.Email)),
		optionalParams...,
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	// Fetch full reimburse with employee
	full, err := r.client.Reimburse.FindUnique(db.Reimburse.ID.Equals(du.ID)).With(db.Reimburse.Employee.Fetch()).Exec(ctx)
	if err != nil {
		return &dto.ReimburseResponse{InnerReimburse: du.InnerReimburse}, nil
	}

	res := &dto.ReimburseResponse{
		InnerReimburse: full.InnerReimburse,
	}
	if emp := full.Employee(); emp != nil {
		res.EmployeeName = emp.Name
		res.Department = emp.Department
	}

	return res, nil
}

func (r *reimburseRepository) Update(ctx context.Context, id string, req dto.UpdateReimburseRequest) (*dto.ReimburseResponse, error) {
	var params []db.ReimburseSetParam

	if req.Amount != nil {
		params = append(params, db.Reimburse.Amount.Set(*req.Amount))
	}

	if req.Description != nil {
		params = append(params, db.Reimburse.Description.Set(*req.Description))
	}

	if req.AttachBill != nil {
		params = append(params, db.Reimburse.AttachBill.Set(*req.AttachBill))
	}

	if req.PaymentProof != nil {
		params = append(params, db.Reimburse.PaymentProof.Set(*req.PaymentProof))
	}

	if req.Status != nil {
		params = append(params, db.Reimburse.Status.Set(*req.Status))
	}

	if req.UpdatedByRole != nil {
		params = append(params, db.Reimburse.UpdatedByRole.Set(*req.UpdatedByRole))
	}

	du, err := r.client.Reimburse.FindUnique(
		db.Reimburse.ID.Equals(id),
	).Update(params...).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &dto.ReimburseResponse{InnerReimburse: du.InnerReimburse}, nil
}

func (r *reimburseRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.Reimburse.FindUnique(
		db.Reimburse.ID.Equals(id),
	).Delete().Exec(ctx)
	return err
}

func (r *reimburseRepository) UpdateStatus(ctx context.Context, id string, status string) (*dto.ReimburseResponse, error) {
	du, err := r.client.Reimburse.FindUnique(
		db.Reimburse.ID.Equals(id),
	).Update(
		db.Reimburse.Status.Set(status),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &dto.ReimburseResponse{InnerReimburse: du.InnerReimburse}, nil
}
