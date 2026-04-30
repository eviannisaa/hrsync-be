package repository

import (
	"context"
	"hrsync-backend/internal/db"
	"hrsync-backend/internal/dto"
	"time"
)

type LeaveTypeRepository interface {
	GetAll(ctx context.Context) ([]dto.LeaveTypeResponse, error)
	GetMyCredits(ctx context.Context, email string) ([]dto.LeaveBalanceResponse, error)
	Seed(ctx context.Context) error
	InitializeBalances(ctx context.Context, email string) error
}

type leaveTypeRepository struct {
	client *db.PrismaClient
}

func NewLeaveTypeRepository(client *db.PrismaClient) LeaveTypeRepository {
	return &leaveTypeRepository{client: client}
}

func (r *leaveTypeRepository) GetAll(ctx context.Context) ([]dto.LeaveTypeResponse, error) {
	types, err := r.client.LeaveType.FindMany().Exec(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.LeaveTypeResponse, 0, len(types))
	for _, t := range types {
		responses = append(responses, dto.LeaveTypeResponse{
			InnerLeaveType: t.InnerLeaveType,
			CreatedAt:      &t.CreatedAt,
			UpdatedAt:      &t.UpdatedAt,
		})
	}

	return responses, nil
}

func (r *leaveTypeRepository) GetMyCredits(ctx context.Context, email string) ([]dto.LeaveBalanceResponse, error) {
	balances, err := r.client.LeaveBalance.FindMany(
		db.LeaveBalance.Email.Equals(email),
	).With(
		db.LeaveBalance.LeaveType.Fetch(),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	// If no balances found, initialize them automatically
	if len(balances) == 0 {
		err := r.InitializeBalances(ctx, email)
		if err != nil {
			return nil, err
		}
		// Refetch after initialization
		balances, err = r.client.LeaveBalance.FindMany(
			db.LeaveBalance.Email.Equals(email),
		).With(
			db.LeaveBalance.LeaveType.Fetch(),
		).Exec(ctx)
		if err != nil {
			return nil, err
		}
	}

	if len(balances) == 0 {
		return []dto.LeaveBalanceResponse{}, nil
	}

	// Dynamically compute how many days have actually been consumed per leave type
	// by inspecting the Leave table for statuses that count as used:
	//   COMPLETE  → full duration (endDate - startDate + 1)
	//   ONGOING   → days elapsed from startDate up to today (inclusive)
	//   APPROVED  → full duration (future leave that has been approved)
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	var leaveRows []struct {
		LeaveTypeID string    `json:"leaveTypeId"`
		Status      string    `json:"status"`
		StartDate   time.Time `json:"startDate"`
		EndDate     time.Time `json:"endDate"`
	}

	err = r.client.Prisma.QueryRaw(`
		SELECT "leaveTypeId", status, "startDate", "endDate"
		FROM "Leave"
		WHERE email = $1
		  AND status IN ('COMPLETE', 'ONGOING', 'APPROVED')
	`, email).Exec(ctx, &leaveRows)
	if err != nil {
		// Fall back to stored values on query error
		responses := make([]dto.LeaveBalanceResponse, 0, len(balances))
		for _, b := range balances {
			responses = append(responses, dto.LeaveBalanceResponse{
				Leave:     b.LeaveType().Name,
				Total:     b.Total,
				Used:      b.Used,
				Remaining: b.Remaining,
			})
		}
		return responses, nil
	}

	// Aggregate used days per leaveTypeId
	usedByType := make(map[string]int)
	for _, row := range leaveRows {
		var days int
		switch row.Status {
		case "COMPLETE":
			// Full duration
			days = int(row.EndDate.Sub(row.StartDate).Hours()/24) + 1
		case "ONGOING":
			// Days elapsed from start up to today (inclusive)
			effectiveEnd := row.EndDate
			if today.Before(effectiveEnd) {
				effectiveEnd = today
			}
			days = int(effectiveEnd.Sub(row.StartDate).Hours()/24) + 1
		case "APPROVED":
			// Future approved — count full planned duration
			days = int(row.EndDate.Sub(row.StartDate).Hours()/24) + 1
		}
		if days > 0 {
			usedByType[row.LeaveTypeID] += days
		}
	}

	responses := make([]dto.LeaveBalanceResponse, 0, len(balances))
	for _, b := range balances {
		used := usedByType[b.LeaveTypeID]
		total := b.Total
		remaining := total - used

		responses = append(responses, dto.LeaveBalanceResponse{
			Leave:     b.LeaveType().Name,
			Total:     total,
			Used:      used,
			Remaining: remaining,
		})
	}

	return responses, nil
}

func (r *leaveTypeRepository) Seed(ctx context.Context) error {
	types := []struct {
		Name string
		Days int
	}{
		{"Annual Leave", 12},
		{"Additional Leave", 3},
		{"Sick Leave", 236},
		{"Period Leave", 236},
		{"Half Day Leave", 236},
		{"Work from Home", 236},
	}

	for _, t := range types {
		_, err := r.client.LeaveType.UpsertOne(
			db.LeaveType.Name.Equals(t.Name),
		).Create(
			db.LeaveType.Name.Set(t.Name),
			db.LeaveType.DefaultDays.Set(t.Days),
		).Update(
			db.LeaveType.DefaultDays.Set(t.Days),
		).Exec(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *leaveTypeRepository) InitializeBalances(ctx context.Context, email string) error {
	types, err := r.client.LeaveType.FindMany().Exec(ctx)
	if err != nil {
		return err
	}

	year := time.Now().Year()
	for _, t := range types {
		_, err := r.client.LeaveBalance.UpsertOne(
			db.LeaveBalance.EmailLeaveTypeIDYear(
				db.LeaveBalance.Email.Equals(email),
				db.LeaveBalance.LeaveTypeID.Equals(t.ID),
				db.LeaveBalance.Year.Equals(year),
			),
		).Create(
			db.LeaveBalance.Email.Set(email),
			db.LeaveBalance.Total.Set(t.DefaultDays),
			db.LeaveBalance.Remaining.Set(t.DefaultDays),
			db.LeaveBalance.Year.Set(year),
			db.LeaveBalance.LeaveType.Link(db.LeaveType.ID.Equals(t.ID)),
		).Update(
			db.LeaveBalance.Total.Set(t.DefaultDays),
			db.LeaveBalance.Remaining.Set(t.DefaultDays),
		).Exec(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
