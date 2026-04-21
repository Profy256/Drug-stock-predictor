# QUICKSTART Guide

## Prerequisites

- Python 3.11+
- PostgreSQL 12+
- pip (Python package manager)

## Installation

### 1. Clone/Navigate to Project

```bash
cd Desktop/med-predict-system/med-predict/py-backend
```

### 2. Create Virtual Environment

**Windows:**
```bash
python -m venv venv
venv\Scripts\activate
```

**macOS/Linux:**
```bash
python3 -m venv venv
source venv/bin/activate
```

### 3. Install Dependencies

```bash
pip install -r requirements.txt
```

### 4. Setup Environment

```bash
cp .env.example .env
# Edit .env with your database credentials
```

## Database Setup

### Option A: Docker Compose (Recommended)

```bash
docker-compose up -d
```

This will:
- Start PostgreSQL database
- Create the med_predict database
- Start the Python backend

### Option B: Manual PostgreSQL

```bash
# Create database
createdb -U postgres med_predict

# The application will create tables automatically on startup
```

## Run the Server

### Development Mode

```bash
make dev
# or
uvicorn app.main:app --reload --host 0.0.0.0 --port 8000
```

### Production Mode

```bash
make prod
# or
gunicorn app.main:app -w 4 -k uvicorn.workers.UvicornWorker --bind 0.0.0.0:8000
```

## Access the Application

Once running, visit:

- **API Documentation (Swagger)**: http://localhost:8000/docs
- **Alternative Docs (ReDoc)**: http://localhost:8000/redoc
- **OpenAPI Schema**: http://localhost:8000/openapi.json
- **Health Check**: http://localhost:8000/health

## Test the API

### Login

```bash
curl -X POST "http://localhost:8000/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### Get Current User

```bash
curl -X GET "http://localhost:8000/api/v1/auth/me" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### List Medicines

```bash
curl -X GET "http://localhost:8000/api/v1/stock/medicines" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

## Useful Commands

```bash
# Run tests
make test

# Format code
make format

# Lint code
make lint

# Clean up
make clean

# Docker commands
make docker-build    # Build image
make docker-up       # Start containers
make docker-down     # Stop containers
make docker-logs     # View logs
```

## Troubleshooting

### Database Connection Error
- Ensure PostgreSQL is running
- Check DATABASE_URL in .env file
- Verify database credentials

### Port 8000 Already in Use
```bash
# Change port
uvicorn app.main:app --port 8001
```

### Module Import Errors
```bash
# Reinstall dependencies
pip install -r requirements.txt
```

### Permission Denied (Linux/Mac)
```bash
# Make scripts executable
chmod +x scripts/*.sh
```

## Next Steps

1. Configure environment variables in `.env`
2. Set up your database
3. Run the development server
4. Explore API documentation at `/docs`
5. Implement additional handlers as needed

## Support

For issues or questions:
- Check the README.md for detailed documentation
- Review error messages in the terminal
- Check logs for debugging information
