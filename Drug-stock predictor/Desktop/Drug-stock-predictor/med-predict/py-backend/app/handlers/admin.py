from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.orm import Session
from app.db import get_db
from app.models import UserRole, UserResponse
from app.models.database import User, Batch, PendingRecord, ApprovedVisit
from app.middleware.auth import get_current_user, get_client_ip
from app.models import TokenData
from app.services.audit import AuditService
from fastapi import Request
from datetime import datetime
import logging

logger = logging.getLogger(__name__)

router = APIRouter(prefix="/api/v1/admin", tags=["admin"])


@router.post("/batches/{batch_id}/approve")
async def approve_batch(
    request: Request,
    batch_id: str,
    current_user: TokenData = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    """Approve a batch and move records to approved visits"""
    # Check if user is admin
    if current_user.role != UserRole.ADMIN.value:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Only admins can approve batches"
        )
    
    audit_service = AuditService(db)
    client_ip = get_client_ip(request)
    
    try:
        batch = db.query(Batch).filter(Batch.id == batch_id).first()
        if not batch:
            raise HTTPException(
                status_code=status.HTTP_404_NOT_FOUND,
                detail="Batch not found"
            )
        
        # Get pending records
        pending_records = db.query(PendingRecord).filter(
            PendingRecord.batch_id == batch_id
        ).all()
        
        # Move to approved visits
        for record in pending_records:
            approved_visit = ApprovedVisit(
                pharmacy_id=batch.pharmacy_id,
                medicine_id=record.medicine_id,
                quantity_dispensed=record.quantity_dispensed,
                diagnosis=record.diagnosis,
                patient_data=record.patient_data,
                visit_date=datetime.now(),
            )
            db.add(approved_visit)
        
        # Update batch status
        batch.status = "approved"
        batch.approved_by = current_user.user_id
        batch.record_count = len(pending_records)
        
        db.commit()
        
        # Log action
        audit_service.log_action(
            current_user.user_id,
            current_user.pharmacy_id,
            "BATCH_APPROVED",
            "BATCH",
            batch_id,
            client_ip,
        )
        
        return {"status": "approved", "records_moved": len(pending_records)}
    except HTTPException:
        raise
    except Exception as e:
        db.rollback()
        logger.error(f"Error approving batch: {str(e)}")
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail="Failed to approve batch"
        )


@router.post("/batches/{batch_id}/reject")
async def reject_batch(
    request: Request,
    batch_id: str,
    reason: str,
    current_user: TokenData = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    """Reject a batch"""
    # Check if user is admin
    if current_user.role != UserRole.ADMIN.value:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Only admins can reject batches"
        )
    
    audit_service = AuditService(db)
    client_ip = get_client_ip(request)
    
    try:
        batch = db.query(Batch).filter(Batch.id == batch_id).first()
        if not batch:
            raise HTTPException(
                status_code=status.HTTP_404_NOT_FOUND,
                detail="Batch not found"
            )
        
        batch.status = "rejected"
        batch.rejection_reason = reason
        batch.approved_by = current_user.user_id
        
        db.commit()
        
        # Log action
        audit_service.log_action(
            current_user.user_id,
            current_user.pharmacy_id,
            "BATCH_REJECTED",
            "BATCH",
            batch_id,
            client_ip,
            {"reason": reason}
        )
        
        return {"status": "rejected", "reason": reason}
    except HTTPException:
        raise
    except Exception as e:
        db.rollback()
        logger.error(f"Error rejecting batch: {str(e)}")
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail="Failed to reject batch"
        )


@router.get("/users", response_model=list[UserResponse])
async def list_users(
    current_user: TokenData = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    """List all users in the pharmacy"""
    # Check if user is admin
    if current_user.role != UserRole.ADMIN.value:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Only admins can view all users"
        )
    
    users = db.query(User).filter(
        User.pharmacy_id == current_user.pharmacy_id
    ).all()
    
    return [
        UserResponse(
            id=str(user.id),
            pharmacy_id=str(user.pharmacy_id),
            name=user.name,
            email=user.email,
            role=user.role,
            is_active=user.is_active,
            created_at=user.created_at,
            updated_at=user.updated_at,
        )
        for user in users
    ]
