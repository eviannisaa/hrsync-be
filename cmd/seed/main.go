package main

import (
	"context"
	"hrsync-backend/internal/db"
	"hrsync-backend/internal/seeding"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	url := os.Getenv("DATABASE_URL")
	client := db.NewClient(db.WithDatasourceURL(url))

	if err := client.Prisma.Connect(); err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			log.Fatalf("failed to disconnect: %v", err)
		}
	}()

	ctx := context.Background()

	// Run feature-specific seeders
	seeding.SeedEmployees(ctx, client)
	seeding.SeedAuth(ctx, client)
	seeding.SeedLeaves(ctx, client)
	seeding.SeedOvertimes(ctx, client)
	seeding.SeedFeedbacks(ctx, client)

	log.Println("Seeding completed!")
}
