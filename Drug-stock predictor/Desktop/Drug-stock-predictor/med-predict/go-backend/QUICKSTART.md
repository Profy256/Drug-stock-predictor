# Quick Start Guide

## Local Development (5 minutes)

### Prerequisites
- Go 1.21+
- PostgreSQL 12+
- Git

### Setup

```bash
# 1. Clone and navigate to backend
cd go-backend

# 2. Install dependencies
go mod download

# 3. Create .env from template
cp .env.example .env

# 4. Create database (if not using Docker)
createdb -U postgres medpredict

# 5. Run migrations
psql -U postgres -d medpredict -f migrations/001_init_schema.sql

# 6. Start the server
go run cmd/server/main.go
```

Server runs at `http://localhost:8080`

Check health: `curl http://localhost:8080/health`

---

## Docker Development (2 commands)

```bash
# Start everything
docker-compose up -d

# View logs
docker-compose logs -f
```

Then access:
- Backend API: `http://localhost:8080`
- PostgreSQL: `localhost:5432`

Stop: `docker-compose down`

---

## API Testing

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@pharmacy.com","password":"password123"}'
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user_id": "...",
  "role": "admin"
}
```

### Use Token in Requests
```bash
TOKEN="<token-from-login>"

curl http://localhost:8080/api/v1/stock \
  -H "Authorization: Bearer $TOKEN"
```

---

## Project Structure

```
go-backend/
├── cmd/server/           # Application entry point
├── internal/
│   ├── config/          # Configuration
│   ├── db/              # Database layer
│   ├── handlers/        # HTTP handlers
│   ├── middleware/      # Auth, logging, rate limit
│   ├── models/          # Data structures
│   └── services/        # Business logic
├── migrations/          # SQL migrations
├── .env.example         # Environment template
├── docker-compose.yml   # Docker dev setup
├── Dockerfile           # Container image
├── Makefile             # Common tasks
└── README.md
```

---

## Common Tasks

### Create a new handler
Create file `internal/handlers/feature.go` following the pattern in existing handlers.

### Add a database query
Add method to `internal/db/queries.go` following existing patterns.

### Run tests
```bash
go test ./...
```

### Format code
```bash
go fmt ./...
```

### Build for production
```bash
go build -o med-predict-backend cmd/server/main.go
```

---

## Environment Variables

See `.env.example` for all options:
- `DB_*` - Database credentials
- `PORT` - Server port (default 8080)
- `JWT_SECRET` - JWT signing key
- `LOG_LEVEL` - Logging verbosity
- `ENV` - `development` or `production`

---

## Troubleshooting

**"Connection refused"**
- Check PostgreSQL is running: `pg_isready`
- Verify DB credentials in `.env`

**"Port already in use"**
- Change PORT in `.env` or kill process using port 8080

**"Table doesn't exist"**
- Run migrations: `psql -U postgres -d medpredict -f migrations/001_init_schema.sql`

---

## Next Steps

1. Read [MIGRATION.md](MIGRATION.md) for production database setup
2. Read [DEPLOYMENT.md](DEPLOYMENT.md) for deployment options
3. Review [API documentation](#api-endpoints) below

---

## API Endpoints Summary

| Method | Endpoint | Auth Required | Role |
|--------|----------|---------------|------|
| POST | `/api/v1/auth/login` | ❌ | - |
| POST | `/api/v1/auth/register` | ❌ | - |
| GET | `/api/v1/auth/me` | ✅ | Any |
| GET | `/api/v1/stock` | ✅ | Any |
| POST | `/api/v1/stock` | ✅ | DataEntrant+ |
| PUT | `/api/v1/stock/:id` | ✅ | Admin+ |
| POST | `/api/v1/batches` | ✅ | DataEntrant+ |
| POST | `/api/v1/batches/:id/approve` | ✅ | Admin+ |
| GET | `/api/v1/analytics/trends` | ✅ | Any |
| GET | `/api/v1/admin/audit-logs` | ✅ | Admin+ |
| GET | `/api/v1/dho/pharmacies` | ✅ | DHO |

For complete API documentation, see README.md
