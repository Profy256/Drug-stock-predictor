package models

import "time"

// UserRole enum
type UserRole string

const (
	RoleDataEntrant UserRole = "data_entrant"
	RoleAdmin       UserRole = "admin"
	RoleDHO         UserRole = "dho"
)

// ============================================================
// Authentication Models
// ============================================================

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	PharmacyID string `json:"pharmacy_id" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=6"`
}

type LoginResponse struct {
	Token      string   `json:"token"`
	UserID     string   `json:"user_id"`
	PharmacyID string   `json:"pharmacy_id"`
	Name       string   `json:"name"`
	Email      string   `json:"email"`
	Role       UserRole `json:"role"`
}

type TokenData struct {
	UserID     string   `json:"user_id"`
	PharmacyID string   `json:"pharmacy_id"`
	Role       UserRole `json:"role"`
}

// ============================================================
// Pharmacy Models
// ============================================================

type PharmacyCreate struct {
	Name           string   `json:"name" binding:"required"`
	Region         string   `json:"region" binding:"required"`
	District       string   `json:"district" binding:"required"`
	Lat            *float64 `json:"lat"`
	Lng            *float64 `json:"lng"`
	ContactPhone   *string  `json:"contact_phone"`
	WhatsAppNumber *string  `json:"whatsapp_number"`
}

type PharmacyUpdate struct {
	Name           *string  `json:"name"`
	Region         *string  `json:"region"`
	District       *string  `json:"district"`
	Lat            *float64 `json:"lat"`
	Lng            *float64 `json:"lng"`
	ContactPhone   *string  `json:"contact_phone"`
	WhatsAppNumber *string  `json:"whatsapp_number"`
	IsActive       *bool    `json:"is_active"`
}

type PharmacyResponse struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Region         string    `json:"region"`
	District       string    `json:"district"`
	Lat            *float64  `json:"lat"`
	Lng            *float64  `json:"lng"`
	ContactPhone   *string   `json:"contact_phone"`
	WhatsAppNumber *string   `json:"whatsapp_number"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ============================================================
// User Models
// ============================================================

type UserCreate struct {
	PharmacyID string   `json:"pharmacy_id" binding:"required"`
	Name       string   `json:"name" binding:"required"`
	Email      string   `json:"email" binding:"required,email"`
	Password   string   `json:"password" binding:"required,min=6"`
	Role       UserRole `json:"role" binding:"required"`
}

type UserResponse struct {
	ID         string    `json:"id"`
	PharmacyID string    `json:"pharmacy_id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Role       UserRole  `json:"role"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ============================================================
// Medicine/Stock Models
// ============================================================

type MedicineCreate struct {
	PharmacyID     string  `json:"pharmacy_id" binding:"required"`
	Name           string  `json:"name" binding:"required"`
	GenericName    string  `json:"generic_name"`
	Dosage         string  `json:"dosage"`
	Unit           string  `json:"unit"`
	QuantityOnHand int     `json:"quantity_on_hand" binding:"required,min=0"`
	ReorderLevel   int     `json:"reorder_level" binding:"required,min=0"`
	UnitCost       float64 `json:"unit_cost" binding:"required,min=0"`
	SellingPrice   float64 `json:"selling_price" binding:"required,min=0"`
	ExpiryDate     string  `json:"expiry_date"`
	ManufacturerID string  `json:"manufacturer_id"`
	BatchNumber    string  `json:"batch_number"`
}

type MedicineUpdate struct {
	Name           *string  `json:"name"`
	GenericName    *string  `json:"generic_name"`
	Dosage         *string  `json:"dosage"`
	Unit           *string  `json:"unit"`
	QuantityOnHand *int     `json:"quantity_on_hand"`
	ReorderLevel   *int     `json:"reorder_level"`
	UnitCost       *float64 `json:"unit_cost"`
	SellingPrice   *float64 `json:"selling_price"`
	ExpiryDate     *string  `json:"expiry_date"`
	ManufacturerID *string  `json:"manufacturer_id"`
	BatchNumber    *string  `json:"batch_number"`
	IsActive       *bool    `json:"is_active"`
}

type MedicineResponse struct {
	ID             string    `json:"id"`
	PharmacyID     string    `json:"pharmacy_id"`
	Name           string    `json:"name"`
	GenericName    string    `json:"generic_name"`
	Dosage         string    `json:"dosage"`
	Unit           string    `json:"unit"`
	QuantityOnHand int       `json:"quantity_on_hand"`
	ReorderLevel   int       `json:"reorder_level"`
	UnitCost       float64   `json:"unit_cost"`
	SellingPrice   float64   `json:"selling_price"`
	ExpiryDate     *string   `json:"expiry_date"`
	ManufacturerID *string   `json:"manufacturer_id"`
	BatchNumber    *string   `json:"batch_number"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ============================================================
// Batch Models
// ============================================================

type BatchCreate struct {
	PharmacyID      string      `json:"pharmacy_id" binding:"required"`
	Name            string      `json:"name" binding:"required"`
	Description     string      `json:"description"`
	Status          BatchStatus `json:"status"`
	SubmittedByID   string      `json:"submitted_by_id"`
	ApprovedByID    *string     `json:"approved_by_id"`
	RejectionReason *string     `json:"rejection_reason"`
}

type BatchUpdate struct {
	Name            *string      `json:"name"`
	Description     *string      `json:"description"`
	Status          *BatchStatus `json:"status"`
	ApprovedByID    *string      `json:"approved_by_id"`
	RejectionReason *string      `json:"rejection_reason"`
}

type BatchStatus string

const (
	BatchPending  BatchStatus = "pending"
	BatchApproved BatchStatus = "approved"
	BatchRejected BatchStatus = "rejected"
)

type BatchResponse struct {
	ID              string      `json:"id"`
	PharmacyID      string      `json:"pharmacy_id"`
	Name            string      `json:"name"`
	Description     string      `json:"description"`
	Status          BatchStatus `json:"status"`
	SubmittedByID   string      `json:"submitted_by_id"`
	ApprovedByID    *string     `json:"approved_by_id"`
	RejectionReason *string     `json:"rejection_reason"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

// ============================================================
// Analytics Models
// ============================================================

type PredictionResponse struct {
	MedicineID   string    `json:"medicine_id"`
	MedicineName string    `json:"medicine_name"`
	Current      int       `json:"current"`
	Predicted    float64   `json:"predicted"`
	Trend        string    `json:"trend"`
	Confidence   float64   `json:"confidence"`
	Timestamp    time.Time `json:"timestamp"`
}

type AlertResponse struct {
	ID         string    `json:"id"`
	PharmacyID string    `json:"pharmacy_id"`
	Type       string    `json:"type"`
	Message    string    `json:"message"`
	IsRead     bool      `json:"is_read"`
	CreatedAt  time.Time `json:"created_at"`
}

// ============================================================
// Patient/Records Models
// ============================================================

type PatientFormField struct {
	Name        string   `json:"name" binding:"required"`
	FieldType   string   `json:"field_type" binding:"required"`
	Required    bool     `json:"required"`
	Options     []string `json:"options"`
	Placeholder string   `json:"placeholder"`
}

type PendingRecord struct {
	ID          string                 `json:"id"`
	PharmacyID  string                 `json:"pharmacy_id"`
	PatientData map[string]interface{} `json:"patient_data"`
	CreatedAt   time.Time              `json:"created_at"`
}

type ApprovedVisit struct {
	ID          string                 `json:"id"`
	PharmacyID  string                 `json:"pharmacy_id"`
	PatientData map[string]interface{} `json:"patient_data"`
	ApprovedBy  string                 `json:"approved_by"`
	ApprovedAt  time.Time              `json:"approved_at"`
}

// ============================================================
// Admin Models
// ============================================================

type AuditLog struct {
	ID           string                 `json:"id"`
	UserID       string                 `json:"user_id"`
	PharmacyID   string                 `json:"pharmacy_id"`
	Action       string                 `json:"action"`
	ResourceType string                 `json:"resource_type"`
	ResourceID   string                 `json:"resource_id"`
	Changes      map[string]interface{} `json:"changes"`
	CreatedAt    time.Time              `json:"created_at"`
}
