package seeding

import (
	"context"
	"hrsync-backend/internal/db"
	"hrsync-backend/internal/utils"
	"log"
	"time"
)

func SeedEmployees(ctx context.Context, client *db.PrismaClient) {
	// Seed sample employees
	employees := []struct {
		Name       string
		Email      string
		Phone      string
		Department string
		Position   string
		JoinDate   time.Time
		IsActive   bool
		Location   string
	}{
		{
			Name:       "John Doe",
			Email:      "john.doe@example.com",
			Phone:      "08123456789",
			Department: "Engineer",
			Position:   "Software Engineer",
			JoinDate:   time.Now(),
			IsActive:   true,
			Location:   "Jakarta Selatan",
		},
		{
			Name:       "Jane Smith",
			Email:      "jane.smith@example.com",
			Phone:      "08987654321",
			Department: "Finance",
			Position:   "Finance Manager",
			JoinDate:   time.Now().AddDate(0, -6, 0),
			IsActive:   true,
			Location:   "Tangerang Selatan",
		},
		{
			Name:       "Bob Wilson",
			Email:      "bob.wilson@example.com",
			Phone:      "08112233445",
			Department: "Marketing",
			Position:   "Marketing Lead",
			JoinDate:   time.Now().AddDate(-1, 0, 0),
			IsActive:   false,
			Location:   "Bandung",
		},
		{
			Name:       "Super Admin",
			Email:      "admin@hrsync.com",
			Phone:      "08000000001",
			Department: "Management",
			Position:   "Administrator",
			JoinDate:   time.Now().AddDate(-2, 0, 0),
			IsActive:   true,
			Location:   "Jakarta Barat",
		},
		{
			Name:       "Sample Employee",
			Email:      "employee@hrsync.com",
			Phone:      "08000000002",
			Department: "Ensol",
			Position:   "Field Engineer",
			JoinDate:   time.Now().AddDate(-1, 0, 0),
			IsActive:   true,
			Location:   "Bekasi",
		},
		{
			Name:       "Andi Pratama",
			Email:      "andi.pratama@hrsync.com",
			Phone:      "08111222333",
			Department: "Engineer",
			Position:   "Backend Developer",
			JoinDate:   time.Now().AddDate(-1, -3, 0),
			IsActive:   true,
			Location:   "Depok",
		},
		{
			Name:       "Sari Dewi",
			Email:      "sari.dewi@hrsync.com",
			Phone:      "08222333444",
			Department: "Marketing",
			Position:   "Digital Marketing Specialist",
			JoinDate:   time.Now().AddDate(0, -8, 0),
			IsActive:   true,
			Location:   "Bogor",
		},
		{
			Name:       "Rizky Firmansyah",
			Email:      "rizky.firmansyah@hrsync.com",
			Phone:      "08333444555",
			Department: "Finance",
			Position:   "Financial Analyst",
			JoinDate:   time.Now().AddDate(-2, -1, 0),
			IsActive:   true,
			Location:   "Jakarta Timur",
		},
		{
			Name:       "Putri Handayani",
			Email:      "putri.handayani@hrsync.com",
			Phone:      "08444555666",
			Department: "Ensol",
			Position:   "Project Coordinator",
			JoinDate:   time.Now().AddDate(0, -4, 0),
			IsActive:   true,
			Location:   "Tangerang",
		},
		{
			Name:       "Budi Santoso",
			Email:      "budi.santoso@hrsync.com",
			Phone:      "08555666777",
			Department: "Management",
			Position:   "Operations Manager",
			JoinDate:   time.Now().AddDate(-3, 0, 0),
			IsActive:   true,
			Location:   "Jakarta Utara",
		},
		{
			Name:       "Tester Employee",
			Email:      "tester@hrsync.com",
			Phone:      "08998877665",
			Department: "Engineer",
			Position:   "QA Engineer",
			JoinDate:   time.Now(),
			IsActive:   true,
			Location:   "Jakarta Pusat",
		},
	}

	for _, u := range employees {
		existing, _ := client.Employee.FindUnique(
			db.Employee.Email.Equals(u.Email),
		).Exec(ctx)

		if existing != nil {
			// Employee exists — always sync department and location from seeder data
			_, err := client.Employee.FindUnique(
				db.Employee.Email.Equals(u.Email),
			).Update(
				db.Employee.Department.Set(u.Department),
				db.Employee.Position.Set(u.Position),
				db.Employee.Location.Set(u.Location),
			).Exec(ctx)
			if err != nil {
				log.Printf("Failed to update %s: %v", u.Email, err)
			} else {
				log.Printf("Updated department/location for %s → %s / %s", u.Email, u.Department, u.Location)
			}
			continue
		}

		id := utils.GenerateEmployeeID()
		_, err := client.Employee.CreateOne(
			db.Employee.Name.Set(u.Name),
			db.Employee.Email.Set(u.Email),
			db.Employee.Phone.Set(u.Phone),
			db.Employee.Department.Set(u.Department),
			db.Employee.Position.Set(u.Position),
			db.Employee.JoinDate.Set(u.JoinDate),
			db.Employee.ID.Set(id),
			db.Employee.IsActive.Set(u.IsActive),
			db.Employee.Location.Set(u.Location),
		).Exec(ctx)
		if err != nil {
			log.Printf("Failed to create employee %s: %v", u.Name, err)
		} else {
			log.Printf("Created employee: %s (%s)", u.Name, u.Location)
		}
	}
}
