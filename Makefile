# ─── Infrastructure / Services ────────────────────────────────────────────────

# Start OrbStack (Docker daemon) + hrsync-postgres container + MinIO
# Safe to run multiple times — already-running services are skipped.
start-services:
	@echo "→ Starting OrbStack..."
	@~/.orbstack/bin/orbctl start
	@echo "→ Starting hrsync-postgres container..."
	@~/.orbstack/bin/docker start hrsync-postgres
	@echo "→ Starting MinIO (port 9000)..."
	@if ! nc -z localhost 9000 2>/dev/null; then \
		minio server ~/minio/data --console-address ':9001' > /tmp/minio.log 2>&1 & \
		sleep 2; \
	fi
	@echo "✓ All services running (postgres:5432, minio:9000)"

# Stop only the hrsync containers (does NOT kill OrbStack)
stop-services:
	@echo "→ Stopping hrsync-postgres..."
	@~/.orbstack/bin/docker stop hrsync-postgres || true
	@echo "→ Stopping MinIO..."
	@pkill -f "minio server" || true
	@echo "✓ Services stopped"

# Kill only the Go backend process (NEVER kills OrbStack)
kill-backend:
	@pkill -f "go run cmd/api/main.go" || true
	@pkill -f "cmd/api/main.go" || true

# ─── Prisma ───────────────────────────────────────────────────────────────────

# Kill any orphaned prisma query engine processes
kill-prisma:
	-pkill -f prisma-query-engine

# ─── Run ──────────────────────────────────────────────────────────────────────

# Start services first, then run the API
dev: start-services kill-prisma
	@lsof -ti :8080 | xargs kill -9 2>/dev/null || true
	go run cmd/api/main.go

# Run without starting services (assumes they are already up)
run: kill-prisma
	@lsof -ti :8080 | xargs kill -9 2>/dev/null || true
	go run cmd/api/main.go

seed:
	go run cmd/seed/main.go

# ─── Migrations ───────────────────────────────────────────────────────────────

# Run the database migration in development
migrate-dev: kill-prisma
	go run github.com/steebchen/prisma-client-go migrate dev --schema prisma/schema.prisma

# Apply pending migrations in production
migrate-deploy:
	go run github.com/steebchen/prisma-client-go migrate deploy --schema prisma/schema.prisma

# Generate prisma client
generate:
	go run github.com/steebchen/prisma-client-go generate --schema prisma/schema.prisma
