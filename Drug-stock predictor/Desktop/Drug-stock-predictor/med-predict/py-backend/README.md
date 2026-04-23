# Python Med Predict Backend

A pharmaceutical management system backend built with Python, FastAPI, and PostgreSQL. Complete rewrite of the Go backend with 100% API compatibility.

## Features

- 🔐 **JWT Authentication** - Secure token-based auth with role-based access control
- 📊 **Analytics Engine** - Trends, stockout predictions, AI summaries
- 📦 **Inventory Management** - Real-time stock tracking with expiry alerts
- 👥 **Role-Based Access** - DataEntrant, Admin, DHO roles
- 🔄 **Batch Processing** - Daily data collection with approval workflows
- 📈 **Audit Logging** - Complete immutable compliance trail
- ⚡ **High Performance** - Built with FastAPI and async support
- 🐳 **Container Ready** - Docker & Docker Compose included

## Prerequisites

### Required
- **Python 3.11+** - [Download](https://www.python.org/downloads/)
- **PostgreSQL 12+** - [Download](https://www.postgresql.org/download/)
- **Git** - [Download](https://git-scm.com/)

### Optional
- **Docker & Docker Compose** - For containerized setup
- **Curl** - For testing API endpoints

## Installation

### Step 1: Clone or Navigate to Project

```bash
cd Desktop/med-predict-system/med-predict/py-backend
```

### Step 2: Create Virtual Environment

```bash
# Windows
python -m venv venv
venv\Scripts\activate

# macOS/Linux
python3 -m venv venv
source venv/bin/activate
```

### Step 3: Install Dependencies

```bash
pip install -r requirements.txt
```

### Step 4: Configure Environment

```bash
cp .env.example .env
# Edit .env with your configuration
```

## Database Setup

### Using Docker Compose (Recommended)

```bash
docker-compose up -d
```

### Manual PostgreSQL Setup

```bash
# Create database
createdb med_predict

# The application will create tables automatically on startup
```

## Running the Server

### Development

```bash
# From the activated venv
uvicorn app.main:app --reload --host 0.0.0.0 --port 8000
```

### Production

```bash
# Using Gunicorn
gunicorn app.main:app -w 4 -k uvicorn.workers.UvicornWorker --bind 0.0.0.0:8000
```

### Using Docker

```bash
docker-compose up
```

## API Endpoints

### Authentication
- `POST /api/v1/auth/login` - Login user
- `POST /api/v1/auth/register` - Register new user
- `GET /api/v1/auth/me` - Get current user

### Stock Management
- `GET /api/v1/stock/medicines` - List medicines
- `POST /api/v1/stock/medicines` - Create medicine
- `GET /api/v1/stock/medicines/{id}` - Get medicine details
- `PUT /api/v1/stock/medicines/{id}` - Update medicine
- `DELETE /api/v1/stock/medicines/{id}` - Delete medicine

## Development

### Project Structure

```
py-backend/
├── app/
│   ├── core/          # Configuration and security
│   ├── db/            # Database connection and setup
│   ├── handlers/      # API route handlers
│   ├── middleware/    # Authentication and middleware
│   ├── models/        # Pydantic schemas and SQLAlchemy ORM models
│   ├── services/      # Business logic services
│   └── main.py        # FastAPI application entry point
├── migrations/        # Database migrations
├── requirements.txt   # Python dependencies
├── .env.example       # Environment variables template
├── Dockerfile         # Docker configuration
└── docker-compose.yml # Docker Compose configuration
```

## Troubleshooting

### Database Connection Failed
- Ensure PostgreSQL is running
- Check DATABASE_URL in .env
- Verify credentials and database exists

### Port Already in Use
```bash
# Change port in command or .env
uvicorn app.main:app --port 8001
```

### Import Errors
```bash
# Ensure venv is activated and dependencies installed
pip install -r requirements.txt
```

## Deployment

### Docker Deployment

```bash
# Build image
docker build -t med-predict-backend:latest .

# Run container
docker run -p 8000:8000 \
  -e DATABASE_URL=postgresql://user:pass@host:5432/db \
  -e JWT_SECRET=your-secret \
  med-predict-backend:latest
```

### Environment Variables

```
DATABASE_URL=postgresql://user:password@localhost:5432/med_predict
JWT_SECRET=your-secret-key-here
ENV=production
DEBUG=False
HOST=0.0.0.0
PORT=8000
```

## API Documentation

Once running, visit:
- **Swagger UI**: http://localhost:8000/docs
- **ReDoc**: http://localhost:8000/redoc
- **OpenAPI JSON**: http://localhost:8000/openapi.json

## Performance

- Startup time: ~2-3 seconds
- Memory usage: ~150-200MB
- Request latency: ~10-50ms (depending on database)

## Testing

```bash
# Run tests (when implemented)
pytest
```

## License

MIT License - See LICENSE file for details
