package models

import (
	"encoding/json"
	"time"

	"github.com/lib/pq"
)

// UserRole represents the role of a user
type UserRole string

const (
	RoleDataEntrant UserRole = "data_entrant"
	RoleAdmin       UserRole = "admin"
	RoleDHO         UserRole = "dho"
)

// ============================================================
// Pharmacy
// ============================================================

type Pharmacy struct {
	ID           string    `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Region       string    `json:"region" db:"region"`
	District     string    `json:"district" db:"district"`
	Lat          float64   `json:"lat" db:"lat"`
	Lng          float64   `json:"lng" db:"lng"`
	ContactPhone string    `json:"contact_phone" db:"contact_phone"`
	WhatsAppNum  string    `json:"whatsapp_number" db:"whatsapp_number"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// ============================================================
// User
// ============================================================

type User struct {
	ID           string    `json:"id" db:"id"`
	PharmacyID   string    `json:"pharmacy_id" db:"pharmacy_id"`
	Name         string    `json:"name" db:"name"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Role         UserRole  `json:"role" db:"role"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// ============================================================
// Patient Form Fields (dynamic)
// ============================================================

type PatientFormField struct {
	ID         string         `json:"id" db:"id"`
	PharmacyID string         `json:"pharmacy_id" db:"pharmacy_id"`
	FieldKey   string         `json:"field_key" db:"field_key"`
	Label      string         `json:"label" db:"label"`
	FieldType  string         `json:"field_type" db:"field_type"` // text, number, date, select
	Options    pq.StringArray `json:"options" db:"options"`
	IsRequired bool           `json:"is_required" db:"is_required"`
	IsActive   bool           `json:"is_active" db:"is_active"`
	SortOrder  int            `json:"sort_order" db:"sort_order"`
	CreatedAt  time.Time      `json:"created_at" db:"created_at"`
}

// ============================================================
// Medicine / Stock
// ============================================================

type Medicine struct {
	ID                string    `json:"id" db:"id"`
	PharmacyID        string    `json:"pharmacy_id" db:"pharmacy_id"`
	Name              string    `json:"name" db:"name"`
	GenericName       string    `json:"generic_name" db:"generic_name"`
	Category          string    `json:"category" db:"category"`
	Unit              string    `json:"unit" db:"unit"` // boxes, vials, strips
	QuantityTotal     int       `json:"quantity_total" db:"quantity_total"`
	QuantityRemaining int       `json:"quantity_remaining" db:"quantity_remaining"`
	ExpiryDate        time.Time `json:"expiry_date" db:"expiry_date"`
	BatchNumber       string    `json:"batch_number" db:"batch_number"`
	Supplier          string    `json:"supplier" db:"supplier"`
	UnitCost          float64   `json:"unit_cost" db:"unit_cost"`
	ReorderLevel      int       `json:"reorder_level" db:"reorder_level"`
	NotificationDays  int       `json:"notification_days" db:"notification_days"`
	CreatedBy         string    `json:"created_by" db:"created_by"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`

	// Computed fields
	Status string `json:"status" db:"-"` // ok, expiring, low, expired
}

// ============================================================
// Batch (Daily Data Submission)
// ============================================================

type Batch struct {
	ID              string    `json:"id" db:"id"`
	PharmacyID      string    `json:"pharmacy_id" db:"pharmacy_id"`
	SubmittedBy     string    `json:"submitted_by" db:"submitted_by"`
	Status          string    `json:"status" db:"status"` // pending, approved, rejected
	RejectionReason string    `json:"rejection_reason" db:"rejection_reason"`
	ApprovedBy      *string   `json:"approved_by" db:"approved_by"`
	RecordCount     int       `json:"record_count" db:"record_count"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// PendingRecord is a single patient visit record in pending state
type PendingRecord struct {
	ID                string          `json:"id" db:"id"`
	BatchID           string          `json:"batch_id" db:"batch_id"`
	PatientHash       string          `json:"patient_hash" db:"patient_hash"` // SHA256 of phone/ID
	MedicineID        string          `json:"medicine_id" db:"medicine_id"`
	QuantityDispensed int             `json:"quantity_dispensed" db:"quantity_dispensed"`
	Diagnosis         string          `json:"diagnosis" db:"diagnosis"`
	PatientData       json.RawMessage `json:"patient_data" db:"patient_data"` // JSONB fields
	CreatedAt         time.Time       `json:"created_at" db:"created_at"`
}

// ApprovedVisit is anonymized patient visit after batch approval
type ApprovedVisit struct {
	ID                string          `json:"id" db:"id"`
	PharmacyID        string          `json:"pharmacy_id" db:"pharmacy_id"`
	MedicineID        string          `json:"medicine_id" db:"medicine_id"`
	QuantityDispensed int             `json:"quantity_dispensed" db:"quantity_dispensed"`
	Diagnosis         string          `json:"diagnosis" db:"diagnosis"`
	PatientData       json.RawMessage `json:"patient_data" db:"patient_data"` // JSONB
	VisitDate         time.Time       `json:"visit_date" db:"visit_date"`
	ApprovedAt        time.Time       `json:"approved_at" db:"approved_at"`
}

// ============================================================
// Notifications
// ============================================================

type NotificationLog struct {
	ID         string    `json:"id" db:"id"`
	PharmacyID string    `json:"pharmacy_id" db:"pharmacy_id"`
	Type       string    `json:"type" db:"type"`       // expiry_alert, low_stock, etc
	Channel    string    `json:"channel" db:"channel"` // whatsapp, email
	Recipient  string    `json:"recipient" db:"recipient"`
	Message    string    `json:"message" db:"message"`
	Status     string    `json:"status" db:"status"` // sent, failed
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// ============================================================
// Analytics Cache
// ============================================================

type AnalyticsCache struct {
	ID         string          `json:"id" db:"id"`
	PharmacyID string          `json:"pharmacy_id" db:"pharmacy_id"`
	CacheType  string          `json:"cache_type" db:"cache_type"` // trends, stockout_risk, etc
	Data       json.RawMessage `json:"data" db:"data"`
	ExpiresAt  time.Time       `json:"expires_at" db:"expires_at"`
	CreatedAt  time.Time       `json:"created_at" db:"created_at"`
}

// ============================================================
// Audit Log
// ============================================================

type AuditLog struct {
	ID         string    `json:"id" db:"id"`
	UserID     string    `json:"user_id" db:"user_id"`
	PharmacyID string    `json:"pharmacy_id" db:"pharmacy_id"`
	Action     string    `json:"action" db:"action"`
	EntityType string    `json:"entity_type" db:"entity_type"`
	EntityID   string    `json:"entity_id" db:"entity_id"`
	Details    string    `json:"details" db:"details"` // JSONB
	IPAddress  string    `json:"ip_address" db:"ip_address"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// ============================================================
// Request/Response DTOs
// ============================================================

// LoginRequest for authentication
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse returns JWT and user info
type LoginResponse struct {
	Token      string   `json:"token"`
	UserID     string   `json:"user_id"`
	PharmacyID string   `json:"pharmacy_id"`
	Name       string   `json:"name"`
	Email      string   `json:"email"`
	Role       UserRole `json:"role"`
}

// RegisterRequest for creating new user
type RegisterRequest struct {
	PharmacyID string   `json:"pharmacy_id" binding:"required"`
	Name       string   `json:"name" binding:"required"`
	Email      string   `json:"email" binding:"required,email"`
	Password   string   `json:"password" binding:"required,min=8"`
	Role       UserRole `json:"role" binding:"required"`
}

// AddStockRequest for stock creation
type AddStockRequest struct {
	Name             string    `json:"name" binding:"required"`
	GenericName      string    `json:"generic_name"`
	Category         string    `json:"category"`
	Unit             string    `json:"unit" binding:"required"`
	QuantityTotal    int       `json:"quantity_total" binding:"required,gt=0"`
	ExpiryDate       time.Time `json:"expiry_date" binding:"required"`
	BatchNumber      string    `json:"batch_number"`
	Supplier         string    `json:"supplier"`
	UnitCost         float64   `json:"unit_cost"`
	ReorderLevel     int       `json:"reorder_level"`
	NotificationDays int       `json:"notification_days"`
}

// AdjustStockRequest for quantity updates
type AdjustStockRequest struct {
	QuantityAdjustment int    `json:"quantity_adjustment" binding:"required"`
	Reason             string `json:"reason" binding:"required"`
}

// SubmitBatchRequest for data submission
type SubmitBatchRequest struct {
	Records []struct {
		PatientHash       string          `json:"patient_hash" binding:"required"`
		MedicineID        string          `json:"medicine_id" binding:"required"`
		QuantityDispensed int             `json:"quantity_dispensed" binding:"required,gt=0"`
		Diagnosis         string          `json:"diagnosis" binding:"required"`
		PatientData       json.RawMessage `json:"patient_data"`
	} `json:"records" binding:"required,min=1"`
}

// MedicineStatus constants
const (
	StatusOK       = "ok"
	StatusExpiring = "expiring"
	StatusExpired  = "expired"
	StatusLowStock = "low_stock"
)

// Batch statuses
const (
	BatchStatusPending  = "pending"
	BatchStatusApproved = "approved"
	BatchStatusRejected = "rejected"
)
