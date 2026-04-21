from fastapi import FastAPI, Request
from fastapi.responses import JSONResponse
from fastapi.middleware.cors import CORSMiddleware
import logging
from app.core.config import settings
from app.db import init_db
from app.middleware.middleware import setup_middleware
from app.handlers import auth, stock, analytics, batch, patient, records, admin, dho

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
)
logger = logging.getLogger(__name__)

# Create FastAPI app
app = FastAPI(
    title=settings.PROJECT_NAME,
    description="Pharmaceutical Management System Backend",
    version="1.0.0",
)

# Setup middleware
setup_middleware(app)

# Initialize database
init_db()
logger.info("Database initialized")


# ============================================================
# Global Exception Handlers
# ============================================================

@app.exception_handler(Exception)
async def general_exception_handler(request: Request, exc: Exception):
    logger.error(f"Unhandled exception: {str(exc)}", exc_info=True)
    return JSONResponse(
        status_code=500,
        content={"detail": "Internal server error"},
    )


# ============================================================
# Health Check
# ============================================================

@app.get("/health")
async def health_check():
    """Health check endpoint"""
    return {
        "status": "ok",
        "service": settings.PROJECT_NAME,
    }


# ============================================================
# API Routes
# ============================================================

# Auth routes
app.include_router(auth.router)

# Stock routes
app.include_router(stock.router)

# Analytics routes
app.include_router(analytics.router)

# Batch routes
app.include_router(batch.router)

# Patient routes
app.include_router(patient.router)

# Records routes
app.include_router(records.router)

# Admin routes
app.include_router(admin.router)

# DHO routes
app.include_router(dho.router)


# ============================================================
# Root endpoint
# ============================================================

@app.get("/")
async def root():
    """Root endpoint"""
    return {
        "message": "Welcome to Med Predict Backend",
        "version": "1.0.0",
        "docs": "/docs",
        "openapi": "/openapi.json",
    }


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(
        "app.main:app",
        host=settings.HOST,
        port=settings.PORT,
        reload=settings.DEBUG,
    )
