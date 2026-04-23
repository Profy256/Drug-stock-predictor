from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.orm import Session
from app.db import get_db
from app.models import UserRole, BatchResponse
from app.models.database import Batch
from app.middleware.auth import get_current_user
from app.models import TokenData
import logging

logger = logging.getLogger(__name__)

router = APIRouter(prefix="/api/v1/dho", tags=["dho"])


@router.get("/batches", response_model=list[BatchResponse])
async def list_dho_batches(
    current_user: TokenData = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    """List all batches for DHO review (DHO role only)"""
    # Check if user is DHO
    if current_user.role != UserRole.DHO.value:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Only DHO users can view DHO batches"
        )
    
    # DHO sees all batches from all pharmacies in the region
    batches = db.query(Batch).all()
    
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


@router.get("/batches/{batch_id}/details", response_model=BatchResponse)
async def get_dho_batch_details(
    batch_id: str,
    current_user: TokenData = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    """Get batch details for DHO review"""
    # Check if user is DHO
    if current_user.role != UserRole.DHO.value:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Only DHO users can view batch details"
        )
    
    batch = db.query(Batch).filter(Batch.id == batch_id).first()
    
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
