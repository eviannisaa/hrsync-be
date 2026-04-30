package repository

import (
	"context"
	"fmt"
	"hrsync-backend/internal/db"
	"hrsync-backend/internal/dto"
	"hrsync-backend/internal/model"
	"log"
	"time"
)

type OvertimeRepository interface {
	GetAll(ctx context.Context, params model.ListParams) ([]dto.OvertimeResponse, int, error)
	Create(ctx context.Context, req dto.CreateOvertimeRequest) (*dto.OvertimeResponse, error)
	Update(ctx context.Context, id string, req dto.UpdateOvertimeRequest) (*dto.OvertimeResponse, error)
	UpdateStatus(ctx context.Context, id string, status string) (*dto.OvertimeResponse, error)
	Delete(ctx context.Context, id string) error
}

type overtimeRepository struct {
	client *db.PrismaClient
}

func NewOvertimeRepository(client *db.PrismaClient) OvertimeRepository {
	return &overtimeRepository{client: client}
}

func (r *overtimeRepository) GetAll(ctx context.Context, params model.ListParams) ([]dto.OvertimeResponse, int, error) {
	skip := (params.Page - 1) * params.Limit

	// Build filter
	var filters []db.OvertimeWhereParam

	if params.Email != "" {
		filters = append(filters, db.Overtime.Email.Equals(params.Email))
	}

	if params.Search != "" {
		var rawRes []struct {
			ID string `json:"id"`
		}
		// Use ILIKE for case-insensitive search across all relevant fields
		err := r.client.Prisma.QueryRaw(`
			SELECT o.id FROM "Overtime" o
			LEFT JOIN "Employee" e ON o."email" = e.email
			WHERE o.status ILIKE $1 
			   OR o.type ILIKE $1 
			   OR o.description ILIKE $1 
			   OR e.name ILIKE $1 
			   OR e.department ILIKE $1 
			   OR e.position ILIKE $1
		`, "%"+params.Search+"%").Exec(ctx, &rawRes)

		if err == nil {
			ids := make([]string, len(rawRes))
			for i, r := range rawRes {
				ids[i] = r.ID
			}
			filters = append(filters, db.Overtime.ID.In(ids))
		}
	}

	// Build sort
	sortDir := db.SortOrderAsc
	if params.SortDir == "desc" {
		sortDir = db.SortOrderDesc
	}
	var orderBy []db.OvertimeOrderByParam
	switch params.SortBy {
	case "createdAt":
		orderBy = append(orderBy, db.Overtime.CreatedAt.Order(sortDir))
	case "startDate":
		orderBy = append(orderBy, db.Overtime.StartDate.Order(sortDir))
	case "startTime":
		orderBy = append(orderBy, db.Overtime.StartTime.Order(sortDir))
	case "status":
		orderBy = append(orderBy, db.Overtime.Status.Order(sortDir))
	case "type":
		orderBy = append(orderBy, db.Overtime.Type.Order(sortDir))
	case "description":
		orderBy = append(orderBy, db.Overtime.Description.Order(sortDir))
	default:
		orderBy = append(orderBy, db.Overtime.CreatedAt.Order(sortDir))
	}

	dbOvertimes, err := r.client.Overtime.FindMany(filters...).With(db.Overtime.Employee.Fetch()).OrderBy(orderBy...).Skip(skip).Take(params.Limit).Exec(ctx)
	if err != nil {
		log.Printf("Error in OvertimeRepository.GetAll (FindMany): %v", err)
		return nil, 0, err
	}

	allOvertimes, err := r.client.Overtime.FindMany(filters...).Exec(ctx)
	if err != nil {
		log.Printf("Error in OvertimeRepository.GetAll (Count/FindMany): %v", err)
		return nil, 0, err
	}
	total := len(allOvertimes)

	responses := make([]dto.OvertimeResponse, 0, len(dbOvertimes))
	for _, du := range dbOvertimes {
		res := dto.OvertimeResponse{
			InnerOvertime: du.InnerOvertime,
		}
		if emp := du.Employee(); emp != nil {
			res.EmployeeName = emp.Name
			res.Department = emp.Department
		}
		responses = append(responses, res)
	}

	return responses, total, nil
}

func (r *overtimeRepository) Create(ctx context.Context, req dto.CreateOvertimeRequest) (*dto.OvertimeResponse, error) {
	duration := calculateDuration(req.StartDate, req.EndDate, req.StartTime, req.EndTime)

	du, err := r.client.Overtime.CreateOne(
		db.Overtime.Type.Set(req.Type),
		db.Overtime.Description.Set(req.Description),
		db.Overtime.EndDate.Set(req.EndDate),
		db.Overtime.EndTime.Set(req.EndTime),
		db.Overtime.StartDate.Set(req.StartDate),
		db.Overtime.StartTime.Set(req.StartTime),
		db.Overtime.Employee.Link(db.Employee.Email.Equals(req.Email)),
		// Optional fields below
		db.Overtime.IsDeployment.Set(req.IsDeployment),
		db.Overtime.IsAdditionalLeave.Set(req.IsAdditionalLeave),
		db.Overtime.Duration.Set(duration),
	).Exec(ctx)

	if err != nil {
		log.Printf("Error in OvertimeRepository.Create: %v", err)
		return nil, err
	}

	// Fetch full overtime with employee
	full, err := r.client.Overtime.FindUnique(db.Overtime.ID.Equals(du.ID)).With(db.Overtime.Employee.Fetch()).Exec(ctx)
	if err != nil {
		return nil, err
	}

	res := &dto.OvertimeResponse{
		InnerOvertime: full.InnerOvertime,
	}

	if emp := full.Employee(); emp != nil {
		res.EmployeeName = emp.Name
		res.Department = emp.Department
	}

	return res, nil
}

func (r *overtimeRepository) Update(ctx context.Context, id string, req dto.UpdateOvertimeRequest) (*dto.OvertimeResponse, error) {
	current, err := r.client.Overtime.FindUnique(db.Overtime.ID.Equals(id)).Exec(ctx)
	if err != nil {
		return nil, err
	}

	var params []db.OvertimeSetParam

	startDate := current.StartDate
	if req.StartDate != nil {
		startDate = *req.StartDate
		params = append(params, db.Overtime.StartDate.Set(startDate))
	}

	endDate := current.EndDate
	if req.EndDate != nil {
		endDate = *req.EndDate
		params = append(params, db.Overtime.EndDate.Set(endDate))
	}

	startTime := current.StartTime
	if req.StartTime != nil {
		startTime = *req.StartTime
		params = append(params, db.Overtime.StartTime.Set(startTime))
	}

	endTime := current.EndTime
	if req.EndTime != nil {
		endTime = *req.EndTime
		params = append(params, db.Overtime.EndTime.Set(endTime))
	}

	if req.Type != nil {
		params = append(params, db.Overtime.Type.Set(*req.Type))
	}

	if req.Description != nil {
		params = append(params, db.Overtime.Description.Set(*req.Description))
	}

	if req.IsDeployment != nil {
		params = append(params, db.Overtime.IsDeployment.Set(*req.IsDeployment))
	}

	if req.IsAdditionalLeave != nil {
		params = append(params, db.Overtime.IsAdditionalLeave.Set(*req.IsAdditionalLeave))
	}

	// Recalculate duration automatically in backend
	duration := calculateDuration(startDate, endDate, startTime, endTime)
	params = append(params, db.Overtime.Duration.Set(duration))

	du, err := r.client.Overtime.FindUnique(db.Overtime.ID.Equals(id)).Update(params...).Exec(ctx)
	if err != nil {
		log.Printf("Error in OvertimeRepository.Update: %v", err)
		return nil, err
	}

	// Fetch full
	full, err := r.client.Overtime.FindUnique(db.Overtime.ID.Equals(du.ID)).With(db.Overtime.Employee.Fetch()).Exec(ctx)
	if err != nil {
		return nil, err
	}

	res := &dto.OvertimeResponse{
		InnerOvertime: full.InnerOvertime,
	}

	if emp := full.Employee(); emp != nil {
		res.EmployeeName = emp.Name
		res.Department = emp.Department
	}

	return res, nil
}

func calculateDuration(startDate, endDate time.Time, startTime, endTime string) float64 {
	var startH, startM, endH, endM int
	fmt.Sscanf(startTime, "%d:%d", &startH, &startM)
	fmt.Sscanf(endTime, "%d:%d", &endH, &endM)

	startTotal := startH*60 + startM
	endTotal := endH*60 + endM

	hoursPerDay := float64(endTotal-startTotal) / 60.0
	// Handle midnight wrap: if end time is before start time, assume it ends the next day
	if hoursPerDay < 0 {
		hoursPerDay += 24.0
	}

	days := int(endDate.Sub(startDate).Hours()/24) + 1
	if days < 1 {
		days = 1
	}

	return float64(days) * hoursPerDay
}

func (r *overtimeRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.Overtime.FindUnique(
		db.Overtime.ID.Equals(id),
	).Delete().Exec(ctx)
	return err
}

func (r *overtimeRepository) UpdateStatus(ctx context.Context, id string, status string) (*dto.OvertimeResponse, error) {
	du, err := r.client.Overtime.FindUnique(
		db.Overtime.ID.Equals(id),
	).Update(
		db.Overtime.Status.Set(status),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	// If approved and converted to additional leave, add 1 day to the employee's Additional Leave balance
	if status == "APPROVED" && du.IsAdditionalLeave {
		// Find "Additional Leave" type
		leaveType, err := r.client.LeaveType.FindUnique(db.LeaveType.Name.Equals("Additional Leave")).Exec(ctx)
		if err == nil {
			year := time.Now().Year()
			// Find or create balance
			balance, err := r.client.LeaveBalance.FindUnique(
				db.LeaveBalance.EmailLeaveTypeIDYear(
					db.LeaveBalance.Email.Equals(du.Email),
					db.LeaveBalance.LeaveTypeID.Equals(leaveType.ID),
					db.LeaveBalance.Year.Equals(year),
				),
			).Exec(ctx)

			if err == nil {
				// Increment total by 1
				_, _ = r.client.LeaveBalance.FindUnique(db.LeaveBalance.ID.Equals(balance.ID)).Update(
					db.LeaveBalance.Total.Set(balance.Total+1),
					db.LeaveBalance.Remaining.Set(balance.Remaining+1),
				).Exec(ctx)
			} else {
				// Create the balance if it doesn't exist
				_, _ = r.client.LeaveBalance.CreateOne(
					db.LeaveBalance.Email.Set(du.Email),
					db.LeaveBalance.Total.Set(1),
					db.LeaveBalance.Remaining.Set(1),
					db.LeaveBalance.Year.Set(year),
					db.LeaveBalance.LeaveType.Link(db.LeaveType.ID.Equals(leaveType.ID)),
				).Exec(ctx)
			}
		}
	}

	return &dto.OvertimeResponse{InnerOvertime: du.InnerOvertime}, nil
}
