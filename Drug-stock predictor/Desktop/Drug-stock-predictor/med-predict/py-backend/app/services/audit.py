from sqlalchemy.orm import Session
from app.models import User, AuditLog
import logging
from datetime import datetime
from typing import Optional, Dict, Any
import uuid

logger = logging.getLogger(__name__)


class AuditService:
    """Service for audit logging"""
    
    def __init__(self, db: Session):
        self.db = db
    
    def log_action(
        self,
        user_id: Optional[str],
        pharmacy_id: Optional[str],
        action: str,
        entity_type: str,
        entity_id: Optional[str],
        ip_address: str,
        changes: Optional[Dict[str, Any]] = None,
    ):
        """Log an action to the audit log"""
        try:
            audit_log = AuditLog(
                id=str(uuid.uuid4()),
                user_id=user_id,
                pharmacy_id=pharmacy_id,
                action=action,
                entity_type=entity_type,
                entity_id=entity_id,
                ip_address=ip_address,
                changes=changes,
            )
            self.db.add(audit_log)
            self.db.commit()
        except Exception as e:
            logger.error(f"Error logging audit action: {str(e)}")
            self.db.rollback()


class LoggerService:
    """Centralized logging service"""
    
    def __init__(self):
        self.logger = logging.getLogger(__name__)
    
    def info(self, message: str, **kwargs):
        self.logger.info(f"{message} {kwargs}")
    
    def warning(self, message: str, **kwargs):
        self.logger.warning(f"{message} {kwargs}")
    
    def error(self, message: str, **kwargs):
        self.logger.error(f"{message} {kwargs}")
    
    def debug(self, message: str, **kwargs):
        self.logger.debug(f"{message} {kwargs}")
