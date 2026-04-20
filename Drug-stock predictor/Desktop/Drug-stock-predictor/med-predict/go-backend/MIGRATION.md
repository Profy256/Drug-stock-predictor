# Go Backend Migration Guide

## Database Setup

### Prerequisites
- PostgreSQL 12+
- psql CLI tool (or any PostgreSQL client)

### Step 1: Create Database

```bash
createdb -U postgres medpredict
```

Or using psql:
```bash
psql -U postgres
CREATE DATABASE medpredict;
```

### Step 2: Run Migrations

Navigate to the go-backend directory and run the migration SQL:

```bash
psql -U postgres -d medpredict -f migrations/001_init_schema.sql
```

### Step 3: Verify Schema

```bash
psql -U postgres -d medpredict -c "\dt"
```

You should see all the tables listed.

## Development Setup

### 1. Install Dependencies

```bash
go mod download
go mod tidy
```

### 2. Configure Environment

```bash
cp .env.example .env
```

Edit `.env` with your database credentials:
```
DB_HOST=localhost
DB_PORT=5432
DB_NAME=medpredict
DB_USER=postgres
DB_PASSWORD=postgres
```

### 3. Run Application

Development mode with hot reload (requires `air`):
```bash
go install github.com/cosmtrek/air@latest
air
```

Or direct run:
```bash
go run cmd/server/main.go
```

Server will be available at `http://localhost:8080`

### 4. Test Health Check

```bash
curl http://localhost:8080/health
```

## Docker Development

### Using Docker Compose

```bash
docker-compose up -d
```

This will:
- Start PostgreSQL database
- Run migrations automatically
- Start the Go backend
- Expose backend at `http://localhost:8080`
- Expose PostgreSQL at `localhost:5432`

### Accessing Database in Docker

```bash
docker exec -it medpredict-db psql -U postgres -d medpredict
```

### View Logs

```bash
docker-compose logs -f backend
docker-compose logs -f postgres
```

### Stop Services

```bash
docker-compose down
```

To also remove the database volume:
```bash
docker-compose down -v
```

## Production Deployment

### Build Binary

```bash
go build -o med-predict-backend cmd/server/main.go
```

### Build Docker Image

```bash
docker build -t med-predict-backend:latest .
```

### Run with Docker

```bash
docker run -p 8080:8080 \
  -e DB_HOST=postgres-host \
  -e DB_PORT=5432 \
  -e DB_NAME=medpredict \
  -e DB_USER=postgres \
  -e DB_PASSWORD=your-secure-password \
  -e JWT_SECRET=your-jwt-secret \
  -e GIN_MODE=release \
  -e ENV=production \
  med-predict-backend:latest
```

## Troubleshooting

### "Connection refused" Error

- Ensure PostgreSQL is running: `pg_isready -h localhost`
- Check DB_HOST in .env (should be `localhost` for local, `postgres` in Docker)
- Verify PostgreSQL user/password

### Migration Already Exists

If running migrations multiple times:
```bash
psql -U postgres -d medpredict -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
```

Then re-run migrations.

### Port Already in Use

If port 8080 is already in use:
```bash
lsof -i :8080  # Find what's using it
kill -9 <PID>  # Kill the process
```

Or change PORT in .env to a different value.
