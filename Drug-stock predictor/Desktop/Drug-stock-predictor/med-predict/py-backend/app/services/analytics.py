from sqlalchemy.orm import Session
from app.models import Medicine, ApprovedVisit
from datetime import datetime
import logging

logger = logging.getLogger(__name__)


class AnalyticsService:
    """Service for analytics operations"""
    
    def __init__(self, db: Session):
        self.db = db
    
    def get_stockout_predictions(self, pharmacy_id: str):
        """Get medicines that are likely to stockout"""
        medicines = self.db.query(Medicine).filter(
            Medicine.pharmacy_id == pharmacy_id,
        ).all()
        
        predictions = []
        for medicine in medicines:
            if medicine.quantity_remaining <= medicine.reorder_level:
                predictions.append({
                    "id": str(medicine.id),
                    "name": medicine.name,
                    "quantity_remaining": medicine.quantity_remaining,
                    "reorder_level": medicine.reorder_level,
                    "status": "critical" if medicine.quantity_remaining == 0 else "low",
                })
        
        return predictions
    
    def get_trends(self, pharmacy_id: str):
        """Get dispensing trends"""
        visits = self.db.query(ApprovedVisit).filter(
            ApprovedVisit.pharmacy_id == pharmacy_id,
        ).all()
        
        trends = {}
        for visit in visits:
            medicine_id = str(visit.medicine_id)
            if medicine_id not in trends:
                trends[medicine_id] = 0
            trends[medicine_id] += visit.quantity_dispensed
        
        return trends
    
    def get_expiry_alerts(self, pharmacy_id: str):
        """Get medicines expiring soon"""
        medicines = self.db.query(Medicine).filter(
            Medicine.pharmacy_id == pharmacy_id,
        ).all()
        
        alerts = []
        now = datetime.now().date()
        for medicine in medicines:
            days_to_expiry = (medicine.expiry_date - now).days
            if 0 <= days_to_expiry <= medicine.notification_days:
                alerts.append({
                    "id": str(medicine.id),
                    "name": medicine.name,
                    "expiry_date": medicine.expiry_date.isoformat(),
                    "days_to_expiry": days_to_expiry,
                    "quantity_remaining": medicine.quantity_remaining,
                })
        
        return alerts
