# File Structure & Navigation Guide

## Complete File Listing

```
go-backend/
├── cmd/server/
│   └── main.go                    # APPLICATION ENTRY POINT
│                                  # - Gin router setup
│                                  # - Route definitions (26 endpoints)
│                                  # - Service initialization
│                                  # - CORS & middleware configuration
│
├── internal/
│   ├── config/
│   │   └── config.go              # CONFIGURATION MANAGEMENT
│   │                              # - Environment variable loading
│   │                              # - DSN generation for PostgreSQL
│   │                              # - Default values
│   │
│   ├── db/
│   │   ├── database.go            # DATABASE CONNECTION
│   │   │                          # - PostgreSQL connection pool
│   │   │                          # - Connection health checks
│   │   └── queries.go             # DATABASE QUERIES (50+ methods)
│   │                              # - User queries (5 methods)
│   │                              # - Pharmacy queries (4 methods)
│   │                              # - Medicine queries (7 methods)
│   │                              # - Batch queries (4 methods)
│   │                              # - Pending record queries (3 methods)
│   │                              # - Approved visit queries (2 methods)
│   │                              # - Audit log queries (2 methods)
│   │                              # - Notification queries (1 method)
│   │                              # - Form field queries (2 methods)
│   │
│   ├── handlers/
│   │   ├── auth.go                # AUTHENTICATION (3 endpoints)
│   │   │                          # - POST /api/v1/auth/login
│   │   │                          # - POST /api/v1/auth/register
│   │   │                          # - GET  /api/v1/auth/me
│   │   │
│   │   ├── stock.go               # INVENTORY MANAGEMENT (5 endpoints)
│   │   │                          # - GET  /api/v1/stock (list)
│   │   │                          # - POST /api/v1/stock (add)
│   │   │                          # - PUT  /api/v1/stock/:id (update)
│   │   │                          # - GET  /api/v1/stock/expiring
│   │   │                          # - GET  /api/v1/stock/search
│   │   │
│   │   ├── batch.go               # BATCH PROCESSING (6 endpoints)
│   │   │                          # - POST /api/v1/batches (submit)
│   │   │                          # - GET  /api/v1/batches (list)
│   │   │                          # - GET  /api/v1/batches/:id (get details)
│   │   │                          # - POST /api/v1/batches/:id/approve
│   │   │                          # - POST /api/v1/batches/:id/reject
│   │   │                          # - DELETE /api/v1/batches/:id/records/:rid
│   │   │
│   │   ├── analytics.go           # ANALYTICS & PREDICTIONS (4 endpoints)
│   │   │                          # - GET /api/v1/analytics/trends
│   │   │                          # - GET /api/v1/analytics/stockout-risk
│   │   │                          # - GET /api/v1/analytics/ai-summary
│   │   │                          # - GET /api/v1/analytics/regional (DHO)
│   │   │
│   │   ├── patient.go             # PATIENT DATA (1 endpoint)
│   │   │                          # - GET /api/v1/patients/search
│   │   │
│   │   ├── admin.go               # ADMIN FUNCTIONS (3 endpoints)
│   │   │                          # - GET/POST /api/v1/admin/form-fields
│   │   │                          # - GET /api/v1/admin/users
│   │   │                          # - GET /api/v1/admin/audit-logs
│   │   │
│   │   └── dho.go                 # DHO MANAGEMENT (3 endpoints)
│   │                              # - GET  /api/v1/dho/pharmacies
│   │                              # - POST /api/v1/dho/pharmacies
│   │                              # - GET  /api/v1/dho/regional-map
│   │
│   ├── middleware/
│   │   ├── auth.go                # AUTHENTICATION MIDDLEWARE
│   │   │                          # - JWT token validation
│   │   │                          # - Role-based access control
│   │   │                          # - Token generation
│   │   │
│   │   ├── logger.go              # REQUEST LOGGING MIDDLEWARE
│   │   │                          # - Structured request logging
│   │   │
│   │   ├── ratelimit.go           # RATE LIMITING MIDDLEWARE
│   │   │                          # - Auth endpoint limiting
│   │   │                          # - API endpoint limiting
│   │   │
│   │   └── cors.go                # CORS HEADERS MIDDLEWARE
│   │                              # - Cross-origin configuration
│   │
│   ├── models/
│   │   └── models.go              # DATA STRUCTURES (20+ types)
│   │                              # - Pharmacy, User, Medicine
│   │                              # - Batch, PendingRecord, ApprovedVisit
│   │                              # - PatientFormField, NotificationLog
│   │                              # - AnalyticsCache, AuditLog
│   │                              # - Request/Response DTOs
│   │                              # - Constants for statuses
│   │
│   └── services/
│       ├── logger.go              # LOGGING SERVICE
│       │                          # - Structured logging with logrus
│       │                          # - File & console output
│       │
│       ├── audit.go               # AUDIT SERVICE
│       │                          # - Immutable compliance logging
│       │                          # - Action constants
│       │                          # - Automatic JSON serialization
│       │
│       └── analytics.go           # ANALYTICS SERVICE
│                                  # - Trend computation
│                                  # - Stockout risk prediction
│                                  # - AI summary generation
│                                  # - Helper functions
│
├── migrations/
│   └── 001_init_schema.sql        # DATABASE SCHEMA
│                                  # - Full PostgreSQL setup
│                                  # - Tables: pharmacies, users, medicines, etc
│                                  # - Indexes & constraints
│                                  # - Extensions: uuid-ossp, pg_trgm
│
├── Configuration & Build Files
│   ├── go.mod                     # GO MODULE DEFINITION
│   │                              # - Dependencies
│   ├── .env.example               # ENVIRONMENT TEMPLATE
│   │                              # - Database config
│   │                              # - Server settings
│   │                              # - API keys (optional)
│   ├── Dockerfile                 # MULTI-STAGE BUILD
│   │                              # - Production container image
│   ├── docker-compose.yml         # LOCAL DEVELOPMENT
│   │                              # - PostgreSQL service
│   │                              # - Backend service
│   ├── .gitignore                 # GIT IGNORE RULES
│   ├── Makefile                   # DEVELOPMENT TASKS
│   │                              # - build, run, test, migrate
│   │                              # - docker commands
│   │                              # - fmt, lint, clean
│   │
│   └── Documentation
│       ├── README.md              # MAIN DOCUMENTATION
│       │                          # - Features overview
│       │                          # - Quick start
│       │                          # - Project structure
│       │                          # - API endpoints summary
│       │
│       ├── QUICKSTART.md          # 5-MINUTE SETUP
│       │                          # - Local development
│       │                          # - Docker development
│       │                          # - API testing examples
│       │                          # - Troubleshooting
│       │
│       ├── MIGRATION.md           # DATABASE MIGRATIONS
│       │                          # - Prerequisites
│       │                          # - Database creation
│       │                          # - Migration execution
│       │                          # - Docker setup
│       │                          # - Troubleshooting
│       │
│       ├── DEPLOYMENT.md          # PRODUCTION DEPLOYMENT (6 options)
│       │                          # - Linux/systemd
│       │                          # - Docker standalone
│       │                          # - Kubernetes
│       │                          # - AWS (ECS, EB)
│       │                          # - Google Cloud
│       │                          # - DigitalOcean
│       │                          # - Security hardening
│       │
│       ├── OVERVIEW.md            # COMPLETE FEATURE OVERVIEW
│       │                          # - Architecture highlights
│       │                          # - Performance comparison
│       │                          # - Deployment options
│       │
│       └── FILE_MAP.md            # THIS FILE
│                                  # - Navigation guide
```

---

## Quick Navigation by Task

### 🚀 Getting Started
1. Read: [QUICKSTART.md](QUICKSTART.md)
2. Copy: `.env.example` → `.env`
3. Run: `go run cmd/server/main.go`

### 📚 Understanding the Code
1. Entry point: [cmd/server/main.go](cmd/server/main.go)
2. Models: [internal/models/models.go](internal/models/models.go)
3. Database: [internal/db/queries.go](internal/db/queries.go)
4. Handlers: [internal/handlers/](internal/handlers/) (any file)

### 🔐 Security & Auth
- JWT implementation: [internal/middleware/auth.go](internal/middleware/auth.go)
- Login handler: [internal/handlers/auth.go](internal/handlers/auth.go)
- Rate limiting: [internal/middleware/ratelimit.go](internal/middleware/ratelimit.go)

### 💾 Database
- Connection: [internal/db/database.go](internal/db/database.go)
- Queries: [internal/db/queries.go](internal/db/queries.go)
- Schema: [migrations/001_init_schema.sql](migrations/001_init_schema.sql)

### 📊 Business Logic
- Stock/Inventory: [internal/handlers/stock.go](internal/handlers/stock.go)
- Batch processing: [internal/handlers/batch.go](internal/handlers/batch.go)
- Analytics: [internal/services/analytics.go](internal/services/analytics.go)
- Audit logging: [internal/services/audit.go](internal/services/audit.go)

### 📋 Configuration
- Environment setup: [.env.example](.env.example)
- Configuration loading: [internal/config/config.go](internal/config/config.go)
- Docker compose: [docker-compose.yml](docker-compose.yml)

### 🚢 Deployment
- Deployment guide: [DEPLOYMENT.md](DEPLOYMENT.md)
- Dockerfile: [Dockerfile](Dockerfile)
- Systemd service: See DEPLOYMENT.md

---

## File Statistics

| Category | Count |
|----------|-------|
| Go Source Files | 14 |
| Documentation Files | 5 |
| Configuration Files | 3 |
| Migration Files | 1 |
| Container Files | 2 |
| **Total** | **25** |

| Metrics | Value |
|---------|-------|
| API Endpoints | 26 |
| Database Query Methods | 50+ |
| Data Models | 20+ |
| Lines of Code | ~2,500 |
| Lines of Documentation | ~1,000 |

---

## Key Concepts

### Authentication Flow
User → Login → JWT Token → Include in Auth Header → Validate → Route Handler

### Data Flow
Handler → Service → Database → Response

### Error Handling
Try operation → On error → Log error → Return JSON error response

### Database Pattern
SQL query → Row scanning → Struct mapping → Return to caller

---

## Development Workflow

```
1. Modify Go code
   ↓
2. Run locally: go run cmd/server/main.go
   ↓
3. Test API: curl http://localhost:8080/...
   ↓
4. Check logs in console
   ↓
5. Repeat
```

---

## Useful Commands

```bash
# Development
go run cmd/server/main.go          # Run the server
go test ./...                       # Run tests
go fmt ./...                        # Format code

# Docker
docker-compose up -d                # Start services
docker-compose logs -f              # View logs
docker-compose down                 # Stop services

# Make
make run                            # Run server
make build                          # Build binary
make docker-up                      # Docker up
make migrate                        # Run migrations
```

---

## Adding New Features

### Add New Route
1. Create handler function in appropriate file in [internal/handlers/](internal/handlers/)
2. Add route registration in [cmd/server/main.go](cmd/server/main.go)
3. Add middleware if needed

### Add Database Query
1. Add method to [internal/db/queries.go](internal/db/queries.go)
2. Follow existing SQL patterns
3. Handle errors appropriately

### Add New Handler
1. Create new file in [internal/handlers/](internal/handlers/)
2. Follow pattern from existing handlers
3. Register handler instance in main.go
4. Add routes

---

## Support Resources

- **Go Documentation**: https://golang.org/doc/
- **Gin Framework**: https://gin-gonic.com/
- **PostgreSQL**: https://www.postgresql.org/docs/
- **JWT Go**: https://github.com/golang-jwt/jwt

---

## File Dependencies

```
main.go
├── config/config.go
├── db/database.go
│   └── db/queries.go
├── handlers/*
│   ├── models/models.go
│   ├── db/queries.go
│   └── services/*
├── middleware/auth.go
├── middleware/logger.go
├── middleware/ratelimit.go
└── services/*
    ├── db/queries.go
    └── models/models.go
```

---

This guide should help you navigate and understand the codebase. Good luck! 🚀
