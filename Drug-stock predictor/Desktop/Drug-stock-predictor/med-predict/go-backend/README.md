# Med Predict Go Backend

Go backend for the Med Predict pharmaceutical management system.

## Quick Start

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 12+ (or use Docker)

### Development Setup

1. **Clone and navigate to backend directory:**
```bash
cd med-predict/go-backend
```

2. **Create .env file:**
```bash
cp .env.example .env
```

3. **Start database:**
```bash
make docker-up
```

Or use Docker Compose directly:
```bash
docker-compose up -d db
```

4. **Install dependencies:**
```bash
go mod download
```

5. **Run development server:**
```bash
make dev
# Or directly:
go run ./cmd/api/main.go
```

Server will be available at `http://localhost:8000`

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login user
- `GET /api/v1/auth/me` - Get current user

### Stock Management
- `GET /api/v1/stock/medicines` - List medicines
- `POST /api/v1/stock/medicines` - Create medicine
- `GET /api/v1/stock/medicines/:id` - Get medicine details
- `PUT /api/v1/stock/medicines/:id` - Update medicine
- `DELETE /api/v1/stock/medicines/:id` - Delete medicine

### Batch Management
- `GET /api/v1/batch` - List batches
- `POST /api/v1/batch` - Create batch
- `POST /api/v1/batch/:id/approve` - Approve batch (Admin)
- `POST /api/v1/batch/:id/reject` - Reject batch (Admin)

### Analytics
- `GET /api/v1/analytics/predictions` - Get predictions
- `GET /api/v1/analytics/alerts` - Get alerts
- `GET /api/v1/analytics/trends` - Get trends

### Patient Management
- `GET /api/v1/patient/form-fields` - Get form fields
- `PUT /api/v1/patient/form-fields` - Update form fields (Admin)
- `POST /api/v1/patient/pending-records` - Create pending record
- `GET /api/v1/patient/pending-records` - List pending records

### Records
- `GET /api/v1/records/pending` - List pending records
- `GET /api/v1/records/approved` - List approved visits
- `GET /api/v1/records/:id` - Get record details

### Admin (Admin Role)
- `GET /api/v1/admin/users` - List users
- `DELETE /api/v1/admin/users/:id` - Deactivate user
- `GET /api/v1/admin/audit-logs` - Get audit logs
- `GET /api/v1/admin/pharmacies` - List pharmacies

### DHO (DHO Role)
- `GET /api/v1/dho/batches/:id/review` - Review batch
- `GET /api/v1/dho/batches-for-review` - List batches for review
- `GET /api/v1/dho/stats` - Get review stats

## Database

PostgreSQL database schema includes:
- `pharmacies` - Pharmacy information
- `users` - User accounts with roles (data_entrant, admin, dho)
- `medicines` - Medicine/stock items
- `batches` - Data batches (pending/approved/rejected)
- `pending_records` - Pending patient records
- `approved_visits` - Approved patient visits
- `audit_logs` - Activity audit trail

## Available Commands

```bash
make help              # Show all available commands
make dev              # Run in development mode
make build            # Build the application
make run              # Run the built application
make test             # Run tests
make docker-up        # Start Docker containers
make docker-down      # Stop Docker containers
make docker-logs      # View container logs
make clean            # Clean build artifacts
```

## Configuration

Environment variables (create `.env` file):
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=med_predict
JWT_SECRET=your-secret-key-change-in-production
PORT=8000
GIN_MODE=debug
```

## Project Structure

```
go-backend/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── handlers/                # HTTP handlers
│   │   ├── auth.go
│   │   ├── stock.go
│   │   ├── batch.go
│   │   ├── analytics.go
│   │   ├── patient.go
│   │   ├── records.go
│   │   ├── admin.go
│   │   ├── dho.go
│   │   └── routes.go
│   ├── models/
│   │   └── models.go            # Data models
│   ├── db/
│   │   └── db.go                # Database connection
│   └── middleware/
│       ├── middleware.go         # General middleware
│       └── auth.go               # JWT authentication
├── migrations/
│   └── 001_init_schema.sql      # Database schema
├── go.mod                        # Go module definition
├── docker-compose.yml            # Docker composition
├── Dockerfile                    # Docker image definition
├── Makefile                      # Build commands
└── README.md                     # This file
```

## Authentication

The API uses JWT (JSON Web Tokens) for authentication.

1. Register: `POST /api/v1/auth/register`
2. Login: `POST /api/v1/auth/login` - Returns JWT token
3. Use token in Authorization header: `Authorization: Bearer <token>`

## Testing

```bash
# Run all tests
make test

# Run specific test file
go test ./internal/handlers/...

# Run with coverage
go test -cover ./...
```

## Building for Production

```bash
# Build Docker image
make docker-build

# Run with Docker
make docker-up

# Or build binary
make build
./bin/med-predict
```

## Troubleshooting

### Database connection issues:
```bash
# Check if database is running
docker-compose ps

# View logs
docker-compose logs db

# Restart database
docker-compose restart db
```

### Port already in use:
Change port in `.env`:
```
PORT=8001
```

### Go module issues:
```bash
go mod tidy
go mod download
```

## License

Proprietary - Med Predict Project

## Support

For issues and questions, please contact the development team.
