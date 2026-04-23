from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.orm import Session
from app.db import get_db
from app.models import PatientFormFieldCreate, PatientFormFieldResponse
from app.models.database import PatientFormField
from app.middleware.auth import get_current_user
from app.models import TokenData
import logging

logger = logging.getLogger(__name__)

router = APIRouter(prefix="/api/v1/patient", tags=["patient"])


@router.get("/form-fields", response_model=list[PatientFormFieldResponse])
async def list_form_fields(
    current_user: TokenData = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    """List all form fields for the pharmacy"""
    fields = db.query(PatientFormField).filter(
        PatientFormField.pharmacy_id == current_user.pharmacy_id,
        PatientFormField.is_active == True,
    ).order_by(PatientFormField.sort_order).all()
    
    return [
        PatientFormFieldResponse(
            id=str(field.id),
            pharmacy_id=str(field.pharmacy_id),
            field_key=field.field_key,
            label=field.label,
            field_type=field.field_type,
            options=field.options,
            is_required=field.is_required,
            is_active=field.is_active,
            sort_order=field.sort_order,
            created_at=field.created_at,
        )
        for field in fields
    ]


@router.post("/form-fields", response_model=PatientFormFieldResponse)
async def create_form_field(
    field_data: PatientFormFieldCreate,
    current_user: TokenData = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    """Create a new form field"""
    try:
        new_field = PatientFormField(
            pharmacy_id=current_user.pharmacy_id,
            field_key=field_data.field_key,
            label=field_data.label,
            field_type=field_data.field_type,
            options=field_data.options,
            is_required=field_data.is_required,
            sort_order=field_data.sort_order,
        )
        
        db.add(new_field)
        db.commit()
        db.refresh(new_field)
        
        return PatientFormFieldResponse(
            id=str(new_field.id),
            pharmacy_id=str(new_field.pharmacy_id),
            field_key=new_field.field_key,
            label=new_field.label,
            field_type=new_field.field_type,
            options=new_field.options,
            is_required=new_field.is_required,
            is_active=new_field.is_active,
            sort_order=new_field.sort_order,
            created_at=new_field.created_at,
        )
    except Exception as e:
        db.rollback()
        logger.error(f"Error creating form field: {str(e)}")
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail="Failed to create form field"
        )
