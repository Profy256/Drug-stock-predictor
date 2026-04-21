from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.orm import Session
from app.db import get_db
from app.models import PendingRecordCreate, PendingRecordResponse, ApprovedVisitResponse
from app.models.database import PendingRecord, ApprovedVisit, Batch
from app.middleware.auth import get_current_user, get_client_ip
from app.models import TokenData
from app.services.audit import AuditService
from fastapi import Request
import logging

logger = logging.getLogger(__name__)

router = APIRouter(prefix="/api/v1/records", tags=["records"])


@router.post("/pending", response_model=PendingRecordResponse)
async def create_pending_record(
    request: Request,
    record_data: PendingRecordCreate,
    current_user: TokenData = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    """Create a pending record"""
    audit_service = AuditService(db)
    client_ip = get_client_ip(request)
    
    try:
        # Verify batch exists and belongs to user's pharmacy
        batch = db.query(Batch).filter(Batch.id == record_data.batch_id).first()
        if not batch:
            raise HTTPException(
                status_code=status.HTTP_404_NOT_FOUND,
                detail="Batch not found"
            )
        
        new_record = PendingRecord(
            batch_id=record_data.batch_id,
            patient_hash=record_data.patient_hash,
            medicine_id=record_data.medicine_id,
            quantity_dispensed=record_data.quantity_dispensed,
            diagnosis=record_data.diagnosis,
            patient_data=record_data.patient_data,
        )
        
        db.add(new_record)
        db.commit()
        db.refresh(new_record)
        
        # Log action
        audit_service.log_action(
            current_user.user_id,
            current_user.pharmacy_id,
            "PENDING_RECORD_CREATED",
            "PENDING_RECORD",
            str(new_record.id),
            client_ip,
        )
        
        return PendingRecordResponse(
            id=str(new_record.id),
            batch_id=str(new_record.batch_id),
            patient_hash=new_record.patient_hash,
            medicine_id=str(new_record.medicine_id),
            quantity_dispensed=new_record.quantity_dispensed,
            diagnosis=new_record.diagnosis,
            patient_data=new_record.patient_data,
            created_at=new_record.created_at,
        )
    except HTTPException:
        raise
    except Exception as e:
        db.rollback()
        logger.error(f"Error creating pending record: {str(e)}")
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail="Failed to create pending record"
        )


@router.get("/pending/{batch_id}", response_model=list[PendingRecordResponse])
async def list_pending_records(
    batch_id: str,
    current_user: TokenData = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    """List pending records for a batch"""
    records = db.query(PendingRecord).filter(
        PendingRecord.batch_id == batch_id,
    ).all()
    
    return [
        PendingRecordResponse(
            id=str(record.id),
            batch_id=str(record.batch_id),
            patient_hash=record.patient_hash,
            medicine_id=str(record.medicine_id),
            quantity_dispensed=record.quantity_dispensed,
            diagnosis=record.diagnosis,
            patient_data=record.patient_data,
            created_at=record.created_at,
        )
        for record in records
    ]


@router.get("/approved", response_model=list[ApprovedVisitResponse])
async def list_approved_visits(
    current_user: TokenData = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    """List approved visits for the pharmacy"""
    visits = db.query(ApprovedVisit).filter(
        ApprovedVisit.pharmacy_id == current_user.pharmacy_id,
    ).all()
    
    return [
        ApprovedVisitResponse(
            id=str(visit.id),
            pharmacy_id=str(visit.pharmacy_id),
            medicine_id=str(visit.medicine_id),
            quantity_dispensed=visit.quantity_dispensed,
            diagnosis=visit.diagnosis,
            patient_data=visit.patient_data,
            visit_date=visit.visit_date,
            approved_at=visit.approved_at,
        )
        for visit in visits
    ]
