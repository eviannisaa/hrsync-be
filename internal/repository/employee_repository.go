package repository

import (
	"context"
	"hrsync-backend/internal/db"
	"hrsync-backend/internal/dto"
	"hrsync-backend/internal/model"
	"hrsync-backend/internal/utils"
)

type EmployeeRepository interface {
	GetAll(ctx context.Context, params model.ListParams) ([]dto.EmployeeResponse, int, error)
	GetByID(ctx context.Context, id string) (*dto.EmployeeResponse, error)
	Create(ctx context.Context, req dto.CreateEmployeeRequest) (*dto.EmployeeResponse, error)
	Update(ctx context.Context, id string, req dto.UpdateEmployeeRequest) (*dto.EmployeeResponse, error)
	Delete(ctx context.Context, id string) error
	GetByEmail(ctx context.Context, email string) (*dto.EmployeeResponse, error)
	GetOrganization(ctx context.Context) (*dto.EmployeeOrganizationResponse, error)
	UpdateOrganization(ctx context.Context, req dto.UpdateEmployeeOrganizationRequest) (*dto.EmployeeOrganizationResponse, error)
}

type employeeRepository struct {
	client *db.PrismaClient
}

func NewEmployeeRepository(client *db.PrismaClient) EmployeeRepository {
	return &employeeRepository{client: client}
}

func (r *employeeRepository) GetAll(ctx context.Context, params model.ListParams) ([]dto.EmployeeResponse, int, error) {
	skip := (params.Page - 1) * params.Limit

	// Build filter
	var filters []db.EmployeeWhereParam
	if params.Search != "" {
		var rawRes []struct {
			ID string `json:"id"`
		}
		// Use ILIKE for case-insensitive search
		err := r.client.Prisma.QueryRaw(`
			SELECT id FROM "Employee"
			WHERE name ILIKE $1 
			   OR email ILIKE $1 
			   OR department ILIKE $1 
			   OR position ILIKE $1 
			   OR location ILIKE $1
		`, "%"+params.Search+"%").Exec(ctx, &rawRes)

		if err == nil {
			ids := make([]string, len(rawRes))
			for i, r := range rawRes {
				ids[i] = r.ID
			}
			filters = append(filters, db.Employee.ID.In(ids))
		}
	}

	// Build sort
	sortDir := db.SortOrderAsc
	if params.SortDir == "desc" {
		sortDir = db.SortOrderDesc
	}
	var orderBy []db.EmployeeOrderByParam
	switch params.SortBy {
	case "name":
		orderBy = append(orderBy, db.Employee.Name.Order(sortDir))
	case "email":
		orderBy = append(orderBy, db.Employee.Email.Order(sortDir))
	case "department":
		orderBy = append(orderBy, db.Employee.Department.Order(sortDir))
	case "joinDate":
		orderBy = append(orderBy, db.Employee.JoinDate.Order(sortDir))
	default:
		orderBy = append(orderBy, db.Employee.Name.Order(sortDir))
	}

	dbEmployees, err := r.client.Employee.FindMany(filters...).OrderBy(orderBy...).Skip(skip).Take(params.Limit).Exec(ctx)
	if err != nil {
		return nil, 0, err
	}

	allEmployees, err := r.client.Employee.FindMany(filters...).Exec(ctx)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.EmployeeResponse, 0, len(dbEmployees))
	for _, du := range dbEmployees {
		responses = append(responses, dto.EmployeeResponse{
			InnerEmployee: du.InnerEmployee,
			Status:        du.Status,
		})
	}

	return responses, len(allEmployees), nil
}

func (r *employeeRepository) GetByID(ctx context.Context, id string) (*dto.EmployeeResponse, error) {
	du, err := r.client.Employee.FindUnique(
		db.Employee.ID.Equals(id),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &dto.EmployeeResponse{
		InnerEmployee: du.InnerEmployee,
		Status:        du.Status,
	}, nil
}

func (r *employeeRepository) Create(ctx context.Context, req dto.CreateEmployeeRequest) (*dto.EmployeeResponse, error) {
	id := utils.GenerateEmployeeID()

	du, err := r.client.Employee.CreateOne(
		db.Employee.Name.Set(req.Name),
		db.Employee.Email.Set(req.Email),
		db.Employee.Phone.Set(req.Phone),
		db.Employee.Department.Set(req.Department),
		db.Employee.Position.Set(req.Position),
		db.Employee.JoinDate.Set(req.JoinDate),
		db.Employee.ID.Set(id),
		db.Employee.Location.Set(req.Location),
		db.Employee.Status.Set(req.Status),
		db.Employee.CreatedBy.Set(req.CreatedBy),
		db.Employee.UpdatedBy.Set(req.CreatedBy),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &dto.EmployeeResponse{
		InnerEmployee: du.InnerEmployee,
		Status:        du.Status,
	}, nil
}

func (r *employeeRepository) Update(ctx context.Context, id string, req dto.UpdateEmployeeRequest) (*dto.EmployeeResponse, error) {
	var params []db.EmployeeSetParam
	if req.Name != nil {
		params = append(params, db.Employee.Name.Set(*req.Name))
	}
	if req.Email != nil {
		params = append(params, db.Employee.Email.Set(*req.Email))
	}
	if req.Phone != nil {
		params = append(params, db.Employee.Phone.Set(*req.Phone))
	}
	if req.Department != nil {
		params = append(params, db.Employee.Department.Set(*req.Department))
	}
	if req.Position != nil {
		params = append(params, db.Employee.Position.Set(*req.Position))
	}
	if req.JoinDate != nil {
		params = append(params, db.Employee.JoinDate.Set(*req.JoinDate))
	}
	if req.IsActive != nil {
		params = append(params, db.Employee.IsActive.Set(*req.IsActive))
	}
	if req.Latitude != nil {
		params = append(params, db.Employee.Latitude.Set(*req.Latitude))
	}
	if req.Longitude != nil {
		params = append(params, db.Employee.Longitude.Set(*req.Longitude))
	}
	if req.Location != nil {
		params = append(params, db.Employee.Location.Set(*req.Location))
	}
	if req.Status != nil {
		params = append(params, db.Employee.Status.Set(*req.Status))
	}
	if req.UpdatedBy != nil {
		params = append(params, db.Employee.UpdatedBy.Set(*req.UpdatedBy))
	}

	du, err := r.client.Employee.FindUnique(
		db.Employee.ID.Equals(id),
	).Update(params...).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &dto.EmployeeResponse{
		InnerEmployee: du.InnerEmployee,
		Status:        du.Status,
	}, nil
}

func (r *employeeRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.Employee.FindUnique(
		db.Employee.ID.Equals(id),
	).Delete().Exec(ctx)
	return err
}

func (r *employeeRepository) GetByEmail(ctx context.Context, email string) (*dto.EmployeeResponse, error) {
	du, err := r.client.Employee.FindUnique(
		db.Employee.Email.Equals(email),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &dto.EmployeeResponse{
		InnerEmployee: du.InnerEmployee,
		Status:        du.Status,
	}, nil
}

func (r *employeeRepository) GetOrganization(ctx context.Context) (*dto.EmployeeOrganizationResponse, error) {
	org, err := r.client.EmployeeOrganization.FindMany().Exec(ctx)
	if err != nil {
		return nil, err
	}
	if len(org) == 0 {
		return &dto.EmployeeOrganizationResponse{
			ID:                "",
			OrganizationImage: "",
		}, nil
	}
	image, _ := org[0].OrganizationImage()
	return &dto.EmployeeOrganizationResponse{
		ID:                org[0].ID,
		OrganizationImage: image,
		UpdatedAt:         org[0].UpdatedAt,
	}, nil
}

func (r *employeeRepository) UpdateOrganization(ctx context.Context, req dto.UpdateEmployeeOrganizationRequest) (*dto.EmployeeOrganizationResponse, error) {
	// Find if any exists
	orgs, err := r.client.EmployeeOrganization.FindMany().Exec(ctx)
	if err != nil {
		return nil, err
	}

	var res *db.EmployeeOrganizationModel
	if len(orgs) == 0 {
		res, err = r.client.EmployeeOrganization.CreateOne(
			db.EmployeeOrganization.OrganizationImage.Set(req.OrganizationImage),
		).Exec(ctx)
	} else {
		res, err = r.client.EmployeeOrganization.FindUnique(
			db.EmployeeOrganization.ID.Equals(orgs[0].ID),
		).Update(
			db.EmployeeOrganization.OrganizationImage.Set(req.OrganizationImage),
		).Exec(ctx)
	}

	if err != nil {
		return nil, err
	}

	image, _ := res.OrganizationImage()
	return &dto.EmployeeOrganizationResponse{
		ID:                res.ID,
		OrganizationImage: image,
		UpdatedAt:         res.UpdatedAt,
	}, nil
}
