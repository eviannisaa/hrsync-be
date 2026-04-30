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

type LeaveRepository interface {
	GetAll(ctx context.Context, params model.ListParams) ([]dto.LeaveResponse, int, error)
	Create(ctx context.Context, req dto.CreateLeaveRequest) (*dto.LeaveResponse, error)
	Update(ctx context.Context, id string, req dto.UpdateLeaveRequest) (*dto.LeaveResponse, error)
	UpdateStatus(ctx context.Context, id string, status string) (*dto.LeaveResponse, error)
	AutoUpdateOngoing(ctx context.Context) error
	Delete(ctx context.Context, id string) error
	GetSummary(ctx context.Context) ([]dto.LeaveStatusSummary, error)
}

type leaveRepository struct {
	client *db.PrismaClient
}

func NewLeaveRepository(client *db.PrismaClient) LeaveRepository {
	return &leaveRepository{client: client}
}

func (r *leaveRepository) GetAll(ctx context.Context, params model.ListParams) ([]dto.LeaveResponse, int, error) {
	skip := (params.Page - 1) * params.Limit
	log.Printf("[LeaveRepository] GetAll: Page=%d, Limit=%d, SortBy=%s, SortDir=%s", params.Page, params.Limit, params.SortBy, params.SortDir)

	// Build filter
	var filters []db.LeaveWhereParam

	// Owner filter
	if params.Email != "" {
		filters = append(filters, db.Leave.Email.Equals(params.Email))
	}

	if params.Search != "" {
		var rawRes []struct {
			ID string `json:"id"`
		}
		// Use ILIKE for case-insensitive search across all relevant fields
		err := r.client.Prisma.QueryRaw(`
			SELECT l.id FROM "Leave" l
			LEFT JOIN "LeaveType" lt ON l."leaveTypeId" = lt.id
			LEFT JOIN "Employee" e ON l."email" = e.email
			WHERE l.status ILIKE $1 
			   OR l.reason ILIKE $1 
			   OR lt.name ILIKE $1 
			   OR e.name ILIKE $1 
			   OR e.department ILIKE $1 
			   OR e.position ILIKE $1
		`, "%"+params.Search+"%").Exec(ctx, &rawRes)

		if err == nil {
			ids := make([]string, len(rawRes))
			for i, r := range rawRes {
				ids[i] = r.ID
			}
			filters = append(filters, db.Leave.ID.In(ids))
		}
	}

	// Build sort
	sortDir := db.SortOrderDesc
	if params.SortDir == "asc" {
		sortDir = db.SortOrderAsc
	}
	var orderBy []db.LeaveOrderByParam
	switch params.SortBy {
	case "email":
		orderBy = append(orderBy, db.Leave.Email.Order(sortDir))
	case "status":
		orderBy = append(orderBy, db.Leave.Status.Order(sortDir))
	case "leaveType":
		orderBy = append(orderBy, db.Leave.LeaveTypeID.Order(sortDir))
	case "reason":
		orderBy = append(orderBy, db.Leave.Reason.Order(sortDir))
	case "createdAt":
		orderBy = append(orderBy, db.Leave.CreatedAt.Order(sortDir))
	case "startDate", "period":
		orderBy = append(orderBy, db.Leave.StartDate.Order(sortDir))
	case "endDate":
		orderBy = append(orderBy, db.Leave.EndDate.Order(sortDir))
	default:
		orderBy = append(orderBy, db.Leave.CreatedAt.Order(sortDir))
	}

	dbLeaves, err := r.client.Leave.FindMany(filters...).With(
		db.Leave.LeaveType.Fetch(),
		db.Leave.Employee.Fetch(),
	).OrderBy(orderBy...).Skip(skip).Take(params.Limit).Exec(ctx)
	if err != nil {
		log.Printf("Error in LeaveRepository.GetAll (FindMany): %v", err)
		return nil, 0, err
	}

	allLeaves, err := r.client.Leave.FindMany(filters...).Exec(ctx)
	if err != nil {
		log.Printf("Error in LeaveRepository.GetAll (Count/FindMany): %v", err)
		return nil, 0, err
	}
	total := len(allLeaves)

	responses := make([]dto.LeaveResponse, 0, len(dbLeaves))
	for _, du := range dbLeaves {
		res := dto.LeaveResponse{
			InnerLeave: du.InnerLeave,
			CreatedBy: func() *string {
				v, ok := du.CreatedBy()
				if !ok {
					return nil
				}
				return &v
			}(),
			UpdatedBy: func() *string {
				v, ok := du.UpdatedBy()
				if !ok {
					return nil
				}
				return &v
			}(),
			CreatedAt: &du.CreatedAt,
			UpdatedAt: &du.UpdatedAt,
		}
		if lt := du.LeaveType(); lt != nil {
			res.LeaveType = lt.Name
		}
		if emp := du.Employee(); emp != nil {
			res.EmployeeName = emp.Name
			res.Department = emp.Department
		}
		responses = append(responses, res)
	}

	return responses, total, nil
}

func (r *leaveRepository) Create(ctx context.Context, req dto.CreateLeaveRequest) (*dto.LeaveResponse, error) {
	// 1. Normalize dates to start-of-day (UTC midnight) for strict duration math
	startDate := time.Date(req.StartDate.Year(), req.StartDate.Month(), req.StartDate.Day(), 0, 0, 0, 0, time.UTC)
	endDate := time.Date(req.EndDate.Year(), req.EndDate.Month(), req.EndDate.Day(), 0, 0, 0, 0, time.UTC)

	requestedDays := int(endDate.Sub(startDate).Hours()/24) + 1
	if requestedDays <= 0 {
		return nil, fmt.Errorf("invalid leave period: end date must be after or same as start date")
	}

	// 2. Fetch the leave type and check for unlimited status (-1)
	lt, err := r.client.LeaveType.FindUnique(db.LeaveType.ID.Equals(req.LeaveTypeId)).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch leave type: %v", err)
	}

	// 3. Validation Logic
	year := startDate.Year()
	balance, err := r.client.LeaveBalance.FindUnique(
		db.LeaveBalance.EmailLeaveTypeIDYear(
			db.LeaveBalance.Email.Equals(req.Email),
			db.LeaveBalance.LeaveTypeID.Equals(req.LeaveTypeId),
			db.LeaveBalance.Year.Equals(year),
		),
	).Exec(ctx)

	totalAllowed := lt.DefaultDays
	if err == nil {
		totalAllowed = balance.Total
	}

	// Calculate consumed days (including SUBMITTED)
	var consumedRows []struct {
		StartDate time.Time `json:"startDate"`
		EndDate   time.Time `json:"endDate"`
	}
	err = r.client.Prisma.QueryRaw(`
		SELECT "startDate", "endDate"
		FROM "Leave"
		WHERE email = $1
		  AND "leaveTypeId" = $2
		  AND status IN ('COMPLETE', 'ONGOING', 'APPROVED', 'SUBMITTED')
		  AND EXTRACT(YEAR FROM "startDate") = $3
	`, req.Email, req.LeaveTypeId, year).Exec(ctx, &consumedRows)

	consumedDays := 0
	if err == nil {
		for _, row := range consumedRows {
			// Normalize rows for calculation too
			s := time.Date(row.StartDate.Year(), row.StartDate.Month(), row.StartDate.Day(), 0, 0, 0, 0, time.UTC)
			e := time.Date(row.EndDate.Year(), row.EndDate.Month(), row.EndDate.Day(), 0, 0, 0, 0, time.UTC)
			d := int(e.Sub(s).Hours()/24) + 1
			if d > 0 {
				consumedDays += d
			}
		}
	}

	remaining := totalAllowed - consumedDays
	if requestedDays > remaining {
		return nil, fmt.Errorf("requested leave duration (%d days) exceeds your remaining balance (%d days left)", requestedDays, remaining)
	}

	du, err := r.client.Leave.CreateOne(
		db.Leave.Reason.Set(req.Reason),
		db.Leave.EndDate.Set(req.EndDate),
		db.Leave.StartDate.Set(req.StartDate),
		db.Leave.Employee.Link(db.Employee.Email.Equals(req.Email)),
		db.Leave.LeaveType.Link(db.LeaveType.ID.Equals(req.LeaveTypeId)),
		db.Leave.Status.Set("SUBMITTED"),
		db.Leave.CreatedBy.Set(req.CreatedBy),
		db.Leave.UpdatedBy.Set(req.UpdatedBy),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	// Fetch the full leave with leave type and employee after creation
	fullLeave, err := r.client.Leave.FindUnique(db.Leave.ID.Equals(du.ID)).With(
		db.Leave.LeaveType.Fetch(),
		db.Leave.Employee.Fetch(),
	).Exec(ctx)
	if err != nil {
		return &dto.LeaveResponse{InnerLeave: du.InnerLeave}, nil
	}

	res := &dto.LeaveResponse{
		InnerLeave: fullLeave.InnerLeave,
		CreatedBy: func() *string {
			v, ok := fullLeave.CreatedBy()
			if !ok {
				return nil
			}
			return &v
		}(),
		UpdatedBy: func() *string {
			v, ok := fullLeave.UpdatedBy()
			if !ok {
				return nil
			}
			return &v
		}(),
		CreatedAt: &fullLeave.CreatedAt,
		UpdatedAt: &fullLeave.UpdatedAt,
	}
	if lt := fullLeave.LeaveType(); lt != nil {
		res.LeaveType = lt.Name
	}

	return res, nil
}

func (r *leaveRepository) Update(ctx context.Context, id string, req dto.UpdateLeaveRequest) (*dto.LeaveResponse, error) {
	var params []db.LeaveSetParam

	if req.LeaveTypeId != nil {
		params = append(params, db.Leave.LeaveType.Link(db.LeaveType.ID.Equals(*req.LeaveTypeId)))
	}

	if req.StartDate != nil {
		params = append(params, db.Leave.StartDate.Set(*req.StartDate))
	}

	if req.EndDate != nil {
		params = append(params, db.Leave.EndDate.Set(*req.EndDate))
	}

	if req.Reason != nil {
		params = append(params, db.Leave.Reason.Set(*req.Reason))
	}

	if req.UpdatedBy != "" {
		params = append(params, db.Leave.UpdatedBy.Set(req.UpdatedBy))
	}

	du, err := r.client.Leave.FindUnique(
		db.Leave.ID.Equals(id),
	).Update(params...).Exec(ctx)
	if err != nil {
		return nil, err
	}

	// Fetch the full leave with leave type after update
	fullLeave, err := r.client.Leave.FindUnique(db.Leave.ID.Equals(du.ID)).With(db.Leave.LeaveType.Fetch()).Exec(ctx)
	if err != nil {
		return &dto.LeaveResponse{InnerLeave: du.InnerLeave}, nil
	}

	res := &dto.LeaveResponse{
		InnerLeave: fullLeave.InnerLeave,
		CreatedBy: func() *string {
			v, ok := fullLeave.CreatedBy()
			if !ok {
				return nil
			}
			return &v
		}(),
		UpdatedBy: func() *string {
			v, ok := fullLeave.UpdatedBy()
			if !ok {
				return nil
			}
			return &v
		}(),
		CreatedAt: &fullLeave.CreatedAt,
		UpdatedAt: &fullLeave.UpdatedAt,
	}
	if lt := fullLeave.LeaveType(); lt != nil {
		res.LeaveType = lt.Name
	}

	return res, nil
}

func (r *leaveRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.Leave.FindUnique(
		db.Leave.ID.Equals(id),
	).Delete().Exec(ctx)
	return err
}

func (r *leaveRepository) UpdateStatus(ctx context.Context, id string, status string) (*dto.LeaveResponse, error) {
	// EVENT-DRIVEN: Instantly transition to ONGOING if approved today or after start date
	if status == "APPROVED" {
		current, err := r.client.Leave.FindUnique(db.Leave.ID.Equals(id)).Exec(ctx)
		if err == nil {
			now := time.Now()
			today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			tomorrow := today.Add(24 * time.Hour)

			// If period has started or starts today, and has not ended yet
			if !current.StartDate.After(tomorrow) && !current.EndDate.Before(today) {
				status = "ONGOING"
			}
		}
	}

	du, err := r.client.Leave.FindUnique(
		db.Leave.ID.Equals(id),
	).Update(
		db.Leave.Status.Set(status),
	).Exec(ctx)

	// Fetch the full leave with leave type after update
	fullLeave, err := r.client.Leave.FindUnique(db.Leave.ID.Equals(du.ID)).With(db.Leave.LeaveType.Fetch()).Exec(ctx)
	if err != nil {
		return &dto.LeaveResponse{InnerLeave: du.InnerLeave}, nil
	}

	res := &dto.LeaveResponse{
		InnerLeave: fullLeave.InnerLeave,
		CreatedBy: func() *string {
			v, ok := fullLeave.CreatedBy()
			if !ok {
				return nil
			}
			return &v
		}(),
		UpdatedBy: func() *string {
			v, ok := fullLeave.UpdatedBy()
			if !ok {
				return nil
			}
			return &v
		}(),
		CreatedAt: &fullLeave.CreatedAt,
		UpdatedAt: &fullLeave.UpdatedAt,
	}
	if lt := fullLeave.LeaveType(); lt != nil {
		res.LeaveType = lt.Name
	}

	return res, nil
}

func (r *leaveRepository) AutoUpdateOngoing(ctx context.Context) error {
	now := time.Now()
	// Use midnight of today for strict day-based comparison
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	tomorrow := today.Add(24 * time.Hour)

	// Update APPROVED -> ONGOING (if today is within the range)
	// We only update if status is not already ONGOING
	res, err := r.client.Leave.FindMany(
		db.Leave.Status.Equals("APPROVED"),
		db.Leave.StartDate.Lte(tomorrow), // Period has started or starts today
		db.Leave.EndDate.Gte(today),      // Period has not ended yet
	).Update(
		db.Leave.Status.Set("ONGOING"),
	).Exec(ctx)

	if err == nil && res.Count > 0 {
		log.Printf("[LeaveRepository] Auto-updated %d requests to ONGOING based on request period", res.Count)
	}

	// Update APPROVED/ONGOING -> COMPLETE (if endDate is in the past)
	res, err = r.client.Leave.FindMany(
		db.Leave.Status.In([]string{"APPROVED", "ONGOING"}),
		db.Leave.EndDate.Lt(today), // Period ended before today
	).Update(
		db.Leave.Status.Set("COMPLETE"),
	).Exec(ctx)

	if err == nil && res.Count > 0 {
		log.Printf("[LeaveRepository] Auto-updated %d requests to COMPLETE based on request period", res.Count)
	}

	return err
}

func (r *leaveRepository) GetSummary(ctx context.Context) ([]dto.LeaveStatusSummary, error) {
	var rawRes []dto.LeaveStatusSummary

	// Explicitly cast COUNT(*) to int to match the Go struct field type
	err := r.client.Prisma.QueryRaw(`
		SELECT status, COUNT(*)::int as total 
		FROM "Leave" 
		GROUP BY status
	`).Exec(ctx, &rawRes)

	if err != nil {
		log.Printf("[LeaveRepository] Error getting leave summary: %v", err)
		return nil, err
	}

	log.Printf("[LeaveRepository] Retrieved %d leave summary statuses", len(rawRes))
	for _, s := range rawRes {
		log.Printf("[LeaveRepository] - %s: %d", s.Status, s.Total)
	}

	return rawRes, nil
}
