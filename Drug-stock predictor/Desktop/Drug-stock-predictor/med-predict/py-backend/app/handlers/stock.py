from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.orm import Session
from app.db import get_db
from app.models import (
    MedicineCreate,
    MedicineUpdate,
    MedicineResponse,
)
from app.models.database import Medicine
from app.middleware.auth import get_current_user
from app.models import TokenData
from datetime import datetime
import logging

logger = logging.getLogger(__name__)

router = APIRouter(prefix="/api/v1/stock", tags=["stock"])


@router.get("/medicines", response_model=list[MedicineResponse])
async def list_medicines(
    current_user: TokenData = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    """List all medicines for the user's pharmacy"""
    medicines = db.query(Medicine).filter(
        Medicine.pharmacy_id == current_user.pharmacy_id
    ).all()
    
    response = []
    for medicine in medicines:
        # Calculate status
        now = datetime.now().date()
        status_val = "ok"
        if medicine.expiry_date < now:
            status_val = "expired"
        elif (medicine.expiry_date - now).days <= medicine.notification_days:
            status_val = "expiring"
        elif medicine.quantity_remaining <= medicine.reorder_level:
            status_val = "low"
        
        response.append(MedicineResponse(
            id=str(medicine.id),
            pharmacy_id=str(medicine.pharmacy_id),
            name=medicine.name,
            generic_name=medicine.generic_name,
            category=medicine.category,
            unit=medicine.unit,
            quantity_total=medicine.quantity_total,
            quantity_remaining=medicine.quantity_remaining,
            expiry_date=medicine.expiry_date.isoformat(),
            batch_number=medicine.batch_number,
            supplier=medicine.supplier,
            unit_cost=medicine.unit_cost,
            reorder_level=medicine.reorder_level,
            notification_days=medicine.notification_days,
            status=status_val,
            created_by=str(medicine.created_by),
            created_at=medicine.created_at,
            updated_at=medicine.updated_at,
        ))
    
    return response


@router.post("/medicines", response_model=MedicineResponse)
async def create_medicine(
    medicine_data: MedicineCreate,
    current_user: TokenData = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    """Create a new medicine entry"""
    try:
        expiry_date = datetime.fromisoformat(medicine_data.expiry_date).date()
        
        new_medicine = Medicine(
            pharmacy_id=current_user.pharmacy_id,
            name=medicine_data.name,
            generic_name=medicine_data.generic_name,
            category=medicine_data.category,
            unit=medicine_data.unit,
            quantity_total=medicine_data.quantity_total,
            quantity_remaining=medicine_data.quantity_remaining,
            expiry_date=expiry_date,
            batch_number=medicine_data.batch_number,
            supplier=medicine_data.supplier,
            unit_cost=medicine_data.unit_cost,
            reorder_level=medicine_data.reorder_level,
            notification_days=medicine_data.notification_days,
            created_by=current_user.user_id,
        )
        
        db.add(new_medicine)
        db.commit()
        db.refresh(new_medicine)
        
        # Calculate status
        now = datetime.now().date()
        status_val = "ok"
        if new_medicine.expiry_date < now:
            status_val = "expired"
        elif (new_medicine.expiry_date - now).days <= new_medicine.notification_days:
            status_val = "expiring"
        elif new_medicine.quantity_remaining <= new_medicine.reorder_level:
            status_val = "low"
        
        return MedicineResponse(
            id=str(new_medicine.id),
            pharmacy_id=str(new_medicine.pharmacy_id),
            name=new_medicine.name,
            generic_name=new_medicine.generic_name,
            category=new_medicine.category,
            unit=new_medicine.unit,
            quantity_total=new_medicine.quantity_total,
            quantity_remaining=new_medicine.quantity_remaining,
            expiry_date=new_medicine.expiry_date.isoformat(),
            batch_number=new_medicine.batch_number,
            supplier=new_medicine.supplier,
            unit_cost=new_medicine.unit_cost,
            reorder_level=new_medicine.reorder_level,
            notification_days=new_medicine.notification_days,
            status=status_val,
            created_by=str(new_medicine.created_by),
            created_at=new_medicine.created_at,
            updated_at=new_medicine.updated_at,
        )
    except Exception as e:
        db.rollback()
        logger.error(f"Error creating medicine: {str(e)}")
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail="Failed to create medicine"
        )


@router.get("/medicines/{medicine_id}", response_model=MedicineResponse)
async def get_medicine(
    medicine_id: str,
    current_user: TokenData = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    """Get a specific medicine"""
    medicine = db.query(Medicine).filter(
        Medicine.id == medicine_id,
        Medicine.pharmacy_id == current_user.pharmacy_id,
    ).first()
    
    if not medicine:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Medicine not found"
        )
    
    # Calculate status
    now = datetime.now().date()
    status_val = "ok"
    if medicine.expiry_date < now:
        status_val = "expired"
    elif (medicine.expiry_date - now).days <= medicine.notification_days:
        status_val = "expiring"
    elif medicine.quantity_remaining <= medicine.reorder_level:
        status_val = "low"
    
    return MedicineResponse(
        id=str(medicine.id),
        pharmacy_id=str(medicine.pharmacy_id),
        name=medicine.name,
        generic_name=medicine.generic_name,
        category=medicine.category,
        unit=medicine.unit,
        quantity_total=medicine.quantity_total,
        quantity_remaining=medicine.quantity_remaining,
        expiry_date=medicine.expiry_date.isoformat(),
        batch_number=medicine.batch_number,
        supplier=medicine.supplier,
        unit_cost=medicine.unit_cost,
        reorder_level=medicine.reorder_level,
        notification_days=medicine.notification_days,
        status=status_val,
        created_by=str(medicine.created_by),
        created_at=medicine.created_at,
        updated_at=medicine.updated_at,
    )


@router.put("/medicines/{medicine_id}", response_model=MedicineResponse)
async def update_medicine(
    medicine_id: str,
    medicine_data: MedicineUpdate,
    current_user: TokenData = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    """Update a medicine entry"""
    medicine = db.query(Medicine).filter(
        Medicine.id == medicine_id,
        Medicine.pharmacy_id == current_user.pharmacy_id,
    ).first()
    
    if not medicine:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Medicine not found"
        )
    
    try:
        # Update fields if provided
        if medicine_data.name:
            medicine.name = medicine_data.name
        if medicine_data.generic_name:
            medicine.generic_name = medicine_data.generic_name
        if medicine_data.category:
            medicine.category = medicine_data.category
        if medicine_data.quantity_remaining is not None:
            medicine.quantity_remaining = medicine_data.quantity_remaining
        if medicine_data.expiry_date:
            medicine.expiry_date = datetime.fromisoformat(medicine_data.expiry_date).date()
        if medicine_data.batch_number:
            medicine.batch_number = medicine_data.batch_number
        if medicine_data.supplier:
            medicine.supplier = medicine_data.supplier
        if medicine_data.unit_cost is not None:
            medicine.unit_cost = medicine_data.unit_cost
        if medicine_data.reorder_level is not None:
            medicine.reorder_level = medicine_data.reorder_level
        if medicine_data.notification_days is not None:
            medicine.notification_days = medicine_data.notification_days
        
        db.commit()
        db.refresh(medicine)
        
        # Calculate status
        now = datetime.now().date()
        status_val = "ok"
        if medicine.expiry_date < now:
            status_val = "expired"
        elif (medicine.expiry_date - now).days <= medicine.notification_days:
            status_val = "expiring"
        elif medicine.quantity_remaining <= medicine.reorder_level:
            status_val = "low"
        
        return MedicineResponse(
            id=str(medicine.id),
            pharmacy_id=str(medicine.pharmacy_id),
            name=medicine.name,
            generic_name=medicine.generic_name,
            category=medicine.category,
            unit=medicine.unit,
            quantity_total=medicine.quantity_total,
            quantity_remaining=medicine.quantity_remaining,
            expiry_date=medicine.expiry_date.isoformat(),
            batch_number=medicine.batch_number,
            supplier=medicine.supplier,
            unit_cost=medicine.unit_cost,
            reorder_level=medicine.reorder_level,
            notification_days=medicine.notification_days,
            status=status_val,
            created_by=str(medicine.created_by),
            created_at=medicine.created_at,
            updated_at=medicine.updated_at,
        )
    except Exception as e:
        db.rollback()
        logger.error(f"Error updating medicine: {str(e)}")
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail="Failed to update medicine"
        )


@router.delete("/medicines/{medicine_id}", status_code=status.HTTP_204_NO_CONTENT)
async def delete_medicine(
    medicine_id: str,
    current_user: TokenData = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    """Delete a medicine entry"""
    medicine = db.query(Medicine).filter(
        Medicine.id == medicine_id,
        Medicine.pharmacy_id == current_user.pharmacy_id,
    ).first()
    
    if not medicine:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Medicine not found"
        )
    
    try:
        db.delete(medicine)
        db.commit()
    except Exception as e:
        db.rollback()
        logger.error(f"Error deleting medicine: {str(e)}")
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail="Failed to delete medicine"
        )
