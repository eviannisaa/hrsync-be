# HR Sync Backend

Go backend with Prisma ORM and PostgreSQL.

## Prerequisites

- [Go](https://golang.org/doc/install)
- [PostgreSQL](https://www.postgresql.org/download/)
- [Prisma](https://www.prisma.io/docs/getting-started)

## Setup

1. Copy `.env.example` to `.env` and update the `DATABASE_URL`.
   ```bash
   cp .env.example .env
   ```

2. Push the schema to the database:
   ```bash
   go run github.com/steebchen/prisma-client-go db push
   ```

3. Seed data:
   ```bash
   go run cmd/seed/main.go
   ```

4. Run the server:
   ```bash
   go run cmd/api/main.go
   ```

## Project Structure

- `cmd/api/main.go`: API application entry point (Thin).
- `cmd/seed/main.go`: Database seeder.
- `internal/handler/`: HTTP handlers.
- `internal/router/`: Dedicated routing logic.
- `internal/service/`: Business logic layer.
- `internal/repository/`: Data access layer (Prisma).
- `internal/model/`: Domain models.


- `db/`: Generated Prisma Client.
- `schema.prisma`: Prisma schema definition.

## API Endpoints

### GET /api/users
Returns a list of all users.

#### Response
```json
[
  {
    "id": "ck...",
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "08123456789",
    "department": "IT",
    "position": "Software Engineer",
    "joinDate": "2023-01-01T00:00:00Z",
    "isActive": true,
    "createdAt": "2023-01-01T00:00:00Z",
    "updatedAt": "2023-01-01T00:00:00Z"
  }
]
```
