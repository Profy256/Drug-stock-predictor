from pydantic import BaseModel, EmailStr, Field
from typing import Optional, List, Dict, Any
from datetime import datetime
from enum import Enum


class UserRole(str, Enum):
    """User roles"""
    DATA_ENTRANT = "data_entrant"
    ADMIN = "admin"
    DHO = "dho"


# ============================================================
# Authentication
# ============================================================

class LoginRequest(BaseModel):
    email: EmailStr
    password: str


class RegisterRequest(BaseModel):
    pharmacy_id: str
    name: str
    email: EmailStr
    password: str


class LoginResponse(BaseModel):
    token: str
    user_id: str
    pharmacy_id: str
    name: str
    email: str
    role: UserRole


class TokenData(BaseModel):
    user_id: str
    pharmacy_id: str
    role: UserRole


# ============================================================
# Pharmacy
# ============================================================

class PharmacyCreate(BaseModel):
    name: str
    region: str
    district: str
    lat: Optional[float] = None
    lng: Optional[float] = None
    contact_phone: Optional[str] = None
    whatsapp_number: Optional[str] = None


class PharmacyUpdate(BaseModel):
    name: Optional[str] = None
    region: Optional[str] = None
    district: Optional[str] = None
    lat: Optional[float] = None
    lng: Optional[float] = None
    contact_phone: Optional[str] = None
    whatsapp_number: Optional[str] = None
    is_active: Optional[bool] = None


class PharmacyResponse(BaseModel):
    id: str
    name: str
    region: str
    district: str
    lat: Optional[float]
    lng: Optional[float]
    contact_phone: Optional[str]
    whatsapp_number: Optional[str]
    is_active: bool
    created_at: datetime
    updated_at: datetime

    class Config:
        from_attributes = True


# ============================================================
# User
# ============================================================

class UserCreate(BaseModel):
    pharmacy_id: str
    name: str
    email: EmailStr
    password: str
    role: UserRole


class UserResponse(BaseModel):
    id: str
    pharmacy_id: str
    name: str
    email: str
    role: UserRole
    is_active: bool
    created_at: datetime
    updated_at: datetime

    class Config:
        from_attributes = True


class UserMe(BaseModel):
    id: str
    pharmacy_id: str
    name: str
    email: str
    role: UserRole

    class Config:
        from_attributes = True


# ============================================================
# Medicine / Stock
# ============================================================

class MedicineCreate(BaseModel):
    name: str
    generic_name: Optional[str] = None
    category: str
    unit: str = "boxes"
    quantity_total: int
    quantity_remaining: int
    expiry_date: str  # ISO format date
    batch_number: Optional[str] = None
    supplier: Optional[str] = None
    unit_cost: Optional[float] = None
    reorder_level: int = 10
    notification_days: int = 14


class MedicineUpdate(BaseModel):
    name: Optional[str] = None
    generic_name: Optional[str] = None
    category: Optional[str] = None
    quantity_remaining: Optional[int] = None
    expiry_date: Optional[str] = None
    batch_number: Optional[str] = None
    supplier: Optional[str] = None
    unit_cost: Optional[float] = None
    reorder_level: Optional[int] = None
    notification_days: Optional[int] = None


class MedicineResponse(BaseModel):
    id: str
    pharmacy_id: str
    name: str
    generic_name: Optional[str]
    category: str
    unit: str
    quantity_total: int
    quantity_remaining: int
    expiry_date: str
    batch_number: Optional[str]
    supplier: Optional[str]
    unit_cost: Optional[float]
    reorder_level: int
    notification_days: int
    status: str  # ok, expiring, low, expired
    created_by: str
    created_at: datetime
    updated_at: datetime

    class Config:
        from_attributes = True


# ============================================================
# Patient Form Fields
# ============================================================

class PatientFormFieldCreate(BaseModel):
    field_key: str
    label: str
    field_type: str  # text, number, date, select
    options: Optional[List[str]] = None
    is_required: bool = False
    sort_order: int = 0


class PatientFormFieldResponse(BaseModel):
    id: str
    pharmacy_id: str
    field_key: str
    label: str
    field_type: str
    options: Optional[List[str]]
    is_required: bool
    is_active: bool
    sort_order: int
    created_at: datetime

    class Config:
        from_attributes = True


# ============================================================
# Batch (Daily Data Submission)
# ============================================================

class BatchCreate(BaseModel):
    pharmacy_id: str


class BatchResponse(BaseModel):
    id: str
    pharmacy_id: str
    submitted_by: str
    status: str  # pending, approved, rejected
    rejection_reason: Optional[str] = None
    approved_by: Optional[str] = None
    record_count: int
    created_at: datetime
    updated_at: datetime

    class Config:
        from_attributes = True


# ============================================================
# Pending Records
# ============================================================

class PendingRecordCreate(BaseModel):
    batch_id: str
    patient_hash: str
    medicine_id: str
    quantity_dispensed: int
    diagnosis: Optional[str] = None
    patient_data: Dict[str, Any] = {}


class PendingRecordResponse(BaseModel):
    id: str
    batch_id: str
    patient_hash: str
    medicine_id: str
    quantity_dispensed: int
    diagnosis: Optional[str]
    patient_data: Dict[str, Any]
    created_at: datetime

    class Config:
        from_attributes = True


# ============================================================
# Approved Visits
# ============================================================

class ApprovedVisitResponse(BaseModel):
    id: str
    pharmacy_id: str
    medicine_id: str
    quantity_dispensed: int
    diagnosis: Optional[str]
    patient_data: Dict[str, Any]
    visit_date: datetime
    approved_at: datetime

    class Config:
        from_attributes = True


# ============================================================
# Notifications
# ============================================================

class NotificationLogResponse(BaseModel):
    id: str
    pharmacy_id: str
    type: str
    channel: str
    recipient: str
    message: str
    status: str  # sent, failed
    created_at: datetime

    class Config:
        from_attributes = True


# ============================================================
# Audit
# ============================================================

class AuditLogResponse(BaseModel):
    id: str
    user_id: str
    pharmacy_id: str
    action: str
    entity_type: str
    entity_id: str
    ip_address: str
    changes: Optional[Dict[str, Any]]
    created_at: datetime

    class Config:
        from_attributes = True
