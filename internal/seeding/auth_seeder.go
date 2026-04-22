package seeding

import (
	"context"
	"hrsync-backend/internal/db"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func SeedAuth(ctx context.Context, client *db.PrismaClient) {
	users := []struct {
		Email    string
		Password string
		Role     db.Role
	}{
		{
			Email:    "admin@hrsync.com",
			Password: "admin123",
			Role:     db.RoleAdmin,
		},
		{
			Email:    "employee@hrsync.com",
			Password: "employee123",
			Role:     db.RoleEmployee,
		},
		{
			Email:    "john.doe@example.com",
			Password: "password123",
			Role:     db.RoleEmployee,
		},
		{
			Email:    "jane.smith@example.com",
			Password: "password123",
			Role:     db.RoleEmployee,
		},
		{
			Email:    "bob.wilson@example.com",
			Password: "password123",
			Role:     db.RoleEmployee,
		},
		{
			Email:    "tester@hrsync.com",
			Password: "password123",
			Role:     db.RoleEmployee,
		},
	}
	for _, u := range users {
		// Cek apakah sudah ada
		existing, _ := client.User.FindUnique(
			db.User.Email.Equals(u.Email),
		).Exec(ctx)
		if existing != nil {
			log.Printf("User already exists, skipping: %s", u.Email)
			continue
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("failed to hash password for %s: %v", u.Email, err)
			continue
		}

		// Temukan employeeId yang sesuai
		emp, _ := client.Employee.FindUnique(
			db.Employee.Email.Equals(u.Email),
		).Exec(ctx)

		if emp != nil {
			_, err = client.User.CreateOne(
				db.User.Email.Set(u.Email),
				db.User.Password.Set(string(hashed)),
				db.User.Role.Set(u.Role),
				db.User.Employee.Link(db.Employee.ID.Equals(emp.ID)),
			).Exec(ctx)
		} else {
			_, err = client.User.CreateOne(
				db.User.Email.Set(u.Email),
				db.User.Password.Set(string(hashed)),
				db.User.Role.Set(u.Role),
			).Exec(ctx)
		}
		if err != nil {
			log.Printf("failed to create user %s: %v", u.Email, err)
		} else {
			log.Printf("Created user: %s (role: %s)", u.Email, u.Role)
		}
	}
}
