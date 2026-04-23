# Python Backend Configuration Files

This directory contains configuration and setup files for the Python backend.

## Files

- `.env.example` - Environment variables template
- `requirements.txt` - Python dependencies
- `Dockerfile` - Docker image configuration
- `docker-compose.yml` - Docker Compose orchestration
- `Makefile` - Development commands
- `migrations/` - Database migration scripts

## Quick Start

1. Copy `.env.example` to `.env`
2. Configure database connection in `.env`
3. Run `make install` to install dependencies
4. Run `make dev` to start the development server
5. Visit http://localhost:8000/docs for API documentation

## Documentation

See `README.md` for complete documentation.
