package seeding

import (
	"context"
	"hrsync-backend/internal/db"
	"log"
	"time"
)

func SeedOvertimes(ctx context.Context, client *db.PrismaClient) {
	overtimes := []struct {
		Email       string
		StartDate   time.Time
		EndDate     time.Time
		StartTime   string
		EndTime     string
		Type        string
		Description string
		Status      string
	}{
		{
			Email:       "john.doe@example.com",
			StartDate:   time.Now().AddDate(0, 0, 7),
			EndDate:     time.Now().AddDate(0, 0, 7),
			StartTime:   "09:00",
			EndTime:     "11:00",
			Type:        "WEEKEND",
			Description: "test",
			Status:      "APPROVED",
		},
		{
			Email:       "jane.smith@example.com",
			StartDate:   time.Now().AddDate(0, 0, 7),
			EndDate:     time.Now().AddDate(0, 0, 7),
			StartTime:   "09:00",
			EndTime:     "11:00",
			Type:        "WEEKDAY",
			Description: "test",
			Status:      "APPROVED",
		},
		{
			Email:       "bob.wilson@example.com",
			StartDate:   time.Now().AddDate(0, 0, 7),
			EndDate:     time.Now().AddDate(0, 0, 7),
			StartTime:   "09:00",
			EndTime:     "11:00",
			Type:        "HOLIDAY",
			Description: "test",
			Status:      "WAITING",
		},
	}

	for _, l := range overtimes {
		// Cek apakah sudah ada (berdasarkan email dan type untuk menghindari duplikasi data seed)
		existing, _ := client.Overtime.FindFirst(
			db.Overtime.Email.Equals(l.Email),
			db.Overtime.Type.Equals(l.Type),
		).Exec(ctx)
		if existing != nil {
			log.Printf("Overtime record already exists for %s, skipping", l.Email)
			continue
		}

		_, err := client.Overtime.CreateOne(
			db.Overtime.Type.Set(l.Type),
			db.Overtime.Description.Set(l.Description),
			db.Overtime.EndDate.Set(l.EndDate),
			db.Overtime.EndTime.Set(l.EndTime),
			db.Overtime.StartDate.Set(l.StartDate),
			db.Overtime.StartTime.Set(l.StartTime),
			db.Overtime.Employee.Link(db.Employee.Email.Equals(l.Email)),
			db.Overtime.Status.Set(l.Status),
		).Exec(ctx)
		if err != nil {
			log.Printf("failed to create overtime for %s: %v", l.Email, err)
		} else {
			log.Printf("Created overtime for email: %s", l.Email)
		}
	}
}
