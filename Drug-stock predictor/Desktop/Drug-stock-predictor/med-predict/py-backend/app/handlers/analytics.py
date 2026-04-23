from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.orm import Session
from app.db import get_db
from app.middleware.auth import get_current_user
from app.models import TokenData
from app.services.analytics import AnalyticsService
import logging

logger = logging.getLogger(__name__)

router = APIRouter(prefix="/api/v1/analytics", tags=["analytics"])


@router.get("/stockout-predictions")
async def get_stockout_predictions(
    current_user: TokenData = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    """Get medicines likely to stockout"""
    analytics_service = AnalyticsService(db)
    predictions = analytics_service.get_stockout_predictions(current_user.pharmacy_id)
    return {"predictions": predictions}


@router.get("/trends")
async def get_trends(
    current_user: TokenData = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    """Get dispensing trends"""
    analytics_service = AnalyticsService(db)
    trends = analytics_service.get_trends(current_user.pharmacy_id)
    return {"trends": trends}


@router.get("/expiry-alerts")
async def get_expiry_alerts(
    current_user: TokenData = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    """Get medicines expiring soon"""
    analytics_service = AnalyticsService(db)
    alerts = analytics_service.get_expiry_alerts(current_user.pharmacy_id)
    return {"alerts": alerts}
