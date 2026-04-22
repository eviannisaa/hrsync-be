package seeding

import (
	"context"
	"hrsync-backend/internal/db"
	"log"
	"time"
)

func SeedLeaves(ctx context.Context, client *db.PrismaClient) {
	// 1. Seed Leave Types
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

	typeMap := make(map[string]string)
	for _, t := range types {
		lt, err := client.LeaveType.UpsertOne(
			db.LeaveType.Name.Equals(t.Name),
		).Create(
			db.LeaveType.Name.Set(t.Name),
			db.LeaveType.DefaultDays.Set(t.Days),
		).Update(
			db.LeaveType.DefaultDays.Set(t.Days),
		).Exec(ctx)
		if err != nil {
			log.Printf("failed to seed leave type %s: %v", t.Name, err)
			continue
		}
		typeMap[t.Name] = lt.ID
	}

	// 2. Seed Leave Balances for all users
	users, _ := client.User.FindMany().Exec(ctx)
	for _, u := range users {
		for _, t := range types {
			used := 0
			if t.Name == "Annual Leave" {
				used = 4
			}
			
			remaining := t.Days
			if t.Days != -1 {
				remaining = t.Days - used
			}

			_, err := client.LeaveBalance.UpsertOne(
				db.LeaveBalance.EmailLeaveTypeIDYear(
					db.LeaveBalance.Email.Equals(u.Email),
					db.LeaveBalance.LeaveTypeID.Equals(typeMap[t.Name]),
					db.LeaveBalance.Year.Equals(2026),
				),
			).Create(
				db.LeaveBalance.Email.Set(u.Email),
				db.LeaveBalance.Total.Set(t.Days),
				db.LeaveBalance.Remaining.Set(remaining),
				db.LeaveBalance.Year.Set(2026),
				db.LeaveBalance.LeaveType.Link(db.LeaveType.ID.Equals(typeMap[t.Name])),
				db.LeaveBalance.Used.Set(used),
			).Update(
				db.LeaveBalance.Total.Set(t.Days),
				db.LeaveBalance.Used.Set(used),
				db.LeaveBalance.Remaining.Set(remaining),
			).Exec(ctx)
			if err != nil {
				log.Printf("failed to seed balance for %s: %v", u.Email, err)
			}
		}
	}

	// 3. Seed some leaves
	leaves := []struct {
		Email     string
		Reason    string
		TypeName  string
		StartDate time.Time
		EndDate   time.Time
		Status    string
	}{
		{
			Email:     "employee@hrsync.com",
			Reason:    "Family Vacation",
			TypeName:  "Annual Leave",
			StartDate: time.Now().AddDate(0, 0, -10),
			EndDate:   time.Now().AddDate(0, 0, -7),
			Status:    "COMPLETE",
		},
		{
			Email:     "employee@hrsync.com",
			Reason:    "Stomach flu",
			TypeName:  "Sick Leave",
			StartDate: time.Now().AddDate(0, 0, -1),
			EndDate:   time.Now().AddDate(0, 0, 1),
			Status:    "ONGOING",
		},
		{
			Email:     "employee@hrsync.com",
			Reason:    "Doctor appointment",
			TypeName:  "Half Day Leave",
			StartDate: time.Now().AddDate(0, 0, 3),
			EndDate:   time.Now().AddDate(0, 0, 3),
			Status:    "APPROVED",
		},
		{
			Email:     "employee@hrsync.com",
			Reason:    "Renewing passport",
			TypeName:  "Annual Leave",
			StartDate: time.Now().AddDate(0, 0, 5),
			EndDate:   time.Now().AddDate(0, 0, 5),
			Status:    "SUBMITTED",
		},
		{
			Email:     "employee@hrsync.com",
			Reason:    "Personal business",
			TypeName:  "Additional Leave",
			StartDate: time.Now().AddDate(0, 0, 10),
			EndDate:   time.Now().AddDate(0, 0, 12),
			Status:    "REJECTED",
		},
		{
			Email:     "admin@hrsync.com",
			Reason:    "Company offsite prep",
			TypeName:  "Work from Home",
			StartDate: time.Now().AddDate(0, 0, 1),
			EndDate:   time.Now().AddDate(0, 0, 2),
			Status:    "APPROVED",
		},
	}

	for _, l := range leaves {
		existing, _ := client.Leave.FindFirst(
			db.Leave.Email.Equals(l.Email),
			db.Leave.Reason.Equals(l.Reason),
			db.Leave.StartDate.Equals(l.StartDate),
		).Exec(ctx)
		if existing != nil {
			continue
		}

		typeID, ok := typeMap[l.TypeName]
		if !ok {
			log.Printf("Type %s not found for seeding", l.TypeName)
			continue
		}

		_, err := client.Leave.CreateOne(
			db.Leave.Reason.Set(l.Reason),
			db.Leave.EndDate.Set(l.EndDate),
			db.Leave.StartDate.Set(l.StartDate),
			db.Leave.Employee.Link(db.Employee.Email.Equals(l.Email)),
			db.Leave.LeaveType.Link(db.LeaveType.ID.Equals(typeID)),
			db.Leave.Status.Set(l.Status),
		).Exec(ctx)
		if err != nil {
			log.Printf("failed to create leave for %s: %v", l.Email, err)
		}
	}
}
