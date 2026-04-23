from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.orm import Session
from app.db import get_db
from app.models import BatchCreate, BatchResponse
from app.models.database import Batch, User
from app.middleware.auth import get_current_user
from app.models import TokenData
from app.services.audit import AuditService
from app.middleware.auth import get_client_ip
from fastapi import Request
import logging

logger = logging.getLogger(__name__)

router = APIRouter(prefix="/api/v1/batches", tags=["batches"])


@router.post("/", response_model=BatchResponse)
async def create_batch(
    request: Request,
    batch_data: BatchCreate,
    current_user: TokenData = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    """Create a new batch"""
    audit_service = AuditService(db)
    client_ip = get_client_ip(request)
    
    try:
        new_batch = Batch(
            pharmacy_id=current_user.pharmacy_id,
            submitted_by=current_user.user_id,
            status="pending",
            record_count=0,
        )
        
        db.add(new_batch)
        db.commit()
        db.refresh(new_batch)
        
        # Log action
        audit_service.log_action(
            current_user.user_id,
            current_user.pharmacy_id,
            "BATCH_CREATED",
            "BATCH",
            str(new_batch.id),
            client_ip,
        )
        
        return BatchResponse(
            id=str(new_batch.id),
            pharmacy_id=str(new_batch.pharmacy_id),
            submitted_by=str(new_batch.submitted_by),
            status=new_batch.status,
            rejection_reason=new_batch.rejection_reason,
            approved_by=str(new_batch.approved_by) if new_batch.approved_by else None,
            record_count=new_batch.record_count,
            created_at=new_batch.created_at,
            updated_at=new_batch.updated_at,
        )
    except Exception as e:
        db.rollback()
        logger.error(f"Error creating batch: {str(e)}")
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail="Failed to create batch"
        )


@router.get("/", response_model=list[BatchResponse])
async def list_batches(
    current_user: TokenData = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    """List all batches for the pharmacy"""
    batches = db.query(Batch).filter(
        Batch.pharmacy_id == current_user.pharmacy_id
    ).all()
    
    return [
        BatchResponse(
            id=str(batch.id),
            pharmacy_id=str(batch.pharmacy_id),
            submitted_by=str(batch.submitted_by),
            status=batch.status,
            rejection_reason=batch.rejection_reason,
            approved_by=str(batch.approved_by) if batch.approved_by else None,
            record_count=batch.record_count,
            created_at=batch.created_at,
            updated_at=batch.updated_at,
        )
        for batch in batches
    ]


@router.get("/{batch_id}", response_model=BatchResponse)
async def get_batch(
    batch_id: str,
    current_user: TokenData = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    """Get a specific batch"""
    batch = db.query(Batch).filter(
        Batch.id == batch_id,
        Batch.pharmacy_id == current_user.pharmacy_id,
    ).first()
    
    if not batch:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Batch not found"
        )
    
    return BatchResponse(
        id=str(batch.id),
        pharmacy_id=str(batch.pharmacy_id),
        submitted_by=str(batch.submitted_by),
        status=batch.status,
        rejection_reason=batch.rejection_reason,
        approved_by=str(batch.approved_by) if batch.approved_by else None,
        record_count=batch.record_count,
        created_at=batch.created_at,
        updated_at=batch.updated_at,
    )
