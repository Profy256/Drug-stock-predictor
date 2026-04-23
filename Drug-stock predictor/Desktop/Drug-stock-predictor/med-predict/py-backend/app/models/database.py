from sqlalchemy import Column, String, Integer, Float, Boolean, DateTime, Text, Date, JSON, ForeignKey, Index, Enum
from sqlalchemy.dialects.postgresql import UUID, ARRAY
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.sql import func
import uuid
from datetime import datetime
import enum as python_enum

Base = declarative_base()


class UserRole(python_enum.Enum):
    """User role enumeration"""
    data_entrant = "data_entrant"
    admin = "admin"
    dho = "dho"


# ============================================================
# Pharmacy
# ============================================================

class Pharmacy(Base):
    __tablename__ = "pharmacies"

    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    name = Column(String(255), nullable=False)
    region = Column(String(100))
    district = Column(String(100))
    lat = Column(Float)
    lng = Column(Float)
    contact_phone = Column(String(30))
    whatsapp_number = Column(String(30))
    is_active = Column(Boolean, default=True)
    created_at = Column(DateTime(timezone=True), server_default=func.now())
    updated_at = Column(DateTime(timezone=True), server_default=func.now(), onupdate=func.now())


# ============================================================
# User
# ============================================================

class User(Base):
    __tablename__ = "users"

    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    pharmacy_id = Column(UUID(as_uuid=True), ForeignKey("pharmacies.id", ondelete="CASCADE"))
    name = Column(String(255), nullable=False)
    email = Column(String(255), unique=True, nullable=False)
    password_hash = Column(Text, nullable=False)
    role = Column(Enum(UserRole), nullable=False)
    is_active = Column(Boolean, default=True)
    created_at = Column(DateTime(timezone=True), server_default=func.now())
    updated_at = Column(DateTime(timezone=True), server_default=func.now(), onupdate=func.now())

    __table_args__ = (
        Index("idx_users_email", "email"),
        Index("idx_users_pharmacy", "pharmacy_id"),
    )


# ============================================================
# Patient Form Fields
# ============================================================

class PatientFormField(Base):
    __tablename__ = "patient_form_fields"

    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    pharmacy_id = Column(UUID(as_uuid=True), ForeignKey("pharmacies.id", ondelete="CASCADE"))
    field_key = Column(String(100), nullable=False)
    label = Column(String(200), nullable=False)
    field_type = Column(String(50), default="text", nullable=False)
    options = Column(ARRAY(String))
    is_required = Column(Boolean, default=False)
    is_active = Column(Boolean, default=True)
    sort_order = Column(Integer, default=0)
    created_at = Column(DateTime(timezone=True), server_default=func.now())


# ============================================================
# Medicine / Stock
# ============================================================

class Medicine(Base):
    __tablename__ = "medicines"

    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    pharmacy_id = Column(UUID(as_uuid=True), ForeignKey("pharmacies.id", ondelete="CASCADE"))
    name = Column(String(255), nullable=False)
    generic_name = Column(String(255))
    category = Column(String(100))
    unit = Column(String(50), default="boxes")
    quantity_total = Column(Integer, default=0)
    quantity_remaining = Column(Integer, default=0)
    expiry_date = Column(Date, nullable=False)
    batch_number = Column(String(100))
    supplier = Column(String(255))
    unit_cost = Column(Float)
    reorder_level = Column(Integer, default=10)
    notification_days = Column(Integer, default=14)
    created_by = Column(UUID(as_uuid=True), ForeignKey("users.id"))
    created_at = Column(DateTime(timezone=True), server_default=func.now())
    updated_at = Column(DateTime(timezone=True), server_default=func.now(), onupdate=func.now())

    __table_args__ = (
        Index("idx_medicines_pharmacy", "pharmacy_id"),
        Index("idx_medicines_expiry", "expiry_date"),
    )


# ============================================================
# Batches (Daily Data Submissions)
# ============================================================

class Batch(Base):
    __tablename__ = "batches"

    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    pharmacy_id = Column(UUID(as_uuid=True), ForeignKey("pharmacies.id", ondelete="CASCADE"))
    submitted_by = Column(UUID(as_uuid=True), ForeignKey("users.id"))
    status = Column(String(50), default="pending")
    rejection_reason = Column(Text)
    approved_by = Column(UUID(as_uuid=True), ForeignKey("users.id"))
    record_count = Column(Integer, default=0)
    created_at = Column(DateTime(timezone=True), server_default=func.now())
    updated_at = Column(DateTime(timezone=True), server_default=func.now(), onupdate=func.now())

    __table_args__ = (
        Index("idx_batches_pharmacy", "pharmacy_id"),
        Index("idx_batches_status", "status"),
    )


# ============================================================
# Pending Records
# ============================================================

class PendingRecord(Base):
    __tablename__ = "pending_records"

    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    batch_id = Column(UUID(as_uuid=True), ForeignKey("batches.id", ondelete="CASCADE"))
    patient_hash = Column(String(255), nullable=False)
    medicine_id = Column(UUID(as_uuid=True), ForeignKey("medicines.id"))
    quantity_dispensed = Column(Integer, nullable=False)
    diagnosis = Column(Text)
    patient_data = Column(JSON, default={})
    created_at = Column(DateTime(timezone=True), server_default=func.now())


# ============================================================
# Approved Visits
# ============================================================

class ApprovedVisit(Base):
    __tablename__ = "approved_visits"

    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    pharmacy_id = Column(UUID(as_uuid=True), ForeignKey("pharmacies.id", ondelete="CASCADE"))
    medicine_id = Column(UUID(as_uuid=True), ForeignKey("medicines.id"))
    quantity_dispensed = Column(Integer, nullable=False)
    diagnosis = Column(Text)
    patient_data = Column(JSON, default={})
    visit_date = Column(DateTime(timezone=True))
    approved_at = Column(DateTime(timezone=True), server_default=func.now())


# ============================================================
# Notifications
# ============================================================

class NotificationLog(Base):
    __tablename__ = "notification_logs"

    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    pharmacy_id = Column(UUID(as_uuid=True), ForeignKey("pharmacies.id", ondelete="CASCADE"))
    type = Column(String(100), nullable=False)  # expiry_alert, low_stock
    channel = Column(String(50), nullable=False)  # whatsapp, email
    recipient = Column(String(255), nullable=False)
    message = Column(Text, nullable=False)
    status = Column(String(50), default="pending")  # pending, sent, failed
    created_at = Column(DateTime(timezone=True), server_default=func.now())


# ============================================================
# Audit Log
# ============================================================

class AuditLog(Base):
    __tablename__ = "audit_logs"

    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    user_id = Column(UUID(as_uuid=True), ForeignKey("users.id"))
    pharmacy_id = Column(UUID(as_uuid=True), ForeignKey("pharmacies.id"))
    action = Column(String(100), nullable=False)
    entity_type = Column(String(100), nullable=False)
    entity_id = Column(String(255))
    ip_address = Column(String(45))
    changes = Column(JSON)
    created_at = Column(DateTime(timezone=True), server_default=func.now())

    __table_args__ = (
        Index("idx_audit_user", "user_id"),
        Index("idx_audit_pharmacy", "pharmacy_id"),
        Index("idx_audit_action", "action"),
    )
