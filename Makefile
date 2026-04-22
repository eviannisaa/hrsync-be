# Kill any orphaned prisma query engine processes
kill-prisma:
	-pkill -f prisma-query-engine

run: kill-prisma
	go run cmd/api/main.go

seed:
	go run cmd/seed/main.go

# Run the database migration in development
migrate-dev: kill-prisma
	go run github.com/steebchen/prisma-client-go migrate dev --schema prisma/schema.prisma

# Apply pending migrations in production
migrate-deploy:
	go run github.com/steebchen/prisma-client-go migrate deploy --schema prisma/schema.prisma

# Generate prisma client
generate:
	go run github.com/steebchen/prisma-client-go generate --schema prisma/schema.prisma
