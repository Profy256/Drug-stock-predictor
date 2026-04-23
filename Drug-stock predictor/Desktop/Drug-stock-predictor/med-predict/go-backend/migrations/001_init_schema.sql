-- Database schema for Med Predict
-- PostgreSQL

-- Create UUID extension if not exists
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================================
-- Pharmacies Table
-- ============================================================
CREATE TABLE IF NOT EXISTS pharmacies (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    region VARCHAR(255),
    district VARCHAR(255),
    lat DECIMAL(10, 8),
    lng DECIMAL(11, 8),
    contact_phone VARCHAR(20),
    whatsapp_number VARCHAR(20),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================
-- Users Table
-- ============================================================
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(255) PRIMARY KEY,
    pharmacy_id VARCHAR(255) NOT NULL REFERENCES pharmacies(id),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'data_entrant',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_pharmacy_id ON users(pharmacy_id);

-- ============================================================
-- Medicines Table
-- ============================================================
CREATE TABLE IF NOT EXISTS medicines (
    id VARCHAR(255) PRIMARY KEY,
    pharmacy_id VARCHAR(255) NOT NULL REFERENCES pharmacies(id),
    name VARCHAR(255) NOT NULL,
    generic_name VARCHAR(255),
    dosage VARCHAR(50),
    unit VARCHAR(50),
    quantity_on_hand INTEGER DEFAULT 0,
    reorder_level INTEGER DEFAULT 0,
    unit_cost DECIMAL(10, 2),
    selling_price DECIMAL(10, 2),
    expiry_date DATE,
    manufacturer_id VARCHAR(255),
    batch_number VARCHAR(100),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_medicines_pharmacy_id ON medicines(pharmacy_id);
CREATE INDEX IF NOT EXISTS idx_medicines_name ON medicines(name);

-- ============================================================
-- Batches Table
-- ============================================================
CREATE TABLE IF NOT EXISTS batches (
    id VARCHAR(255) PRIMARY KEY,
    pharmacy_id VARCHAR(255) NOT NULL REFERENCES pharmacies(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) DEFAULT 'pending',
    submitted_by_id VARCHAR(255) REFERENCES users(id),
    approved_by_id VARCHAR(255) REFERENCES users(id),
    rejection_reason TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_batches_pharmacy_id ON batches(pharmacy_id);
CREATE INDEX IF NOT EXISTS idx_batches_status ON batches(status);

-- ============================================================
-- Pending Records Table
-- ============================================================
CREATE TABLE IF NOT EXISTS pending_records (
    id VARCHAR(255) PRIMARY KEY,
    pharmacy_id VARCHAR(255) NOT NULL REFERENCES pharmacies(id),
    patient_data JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_pending_records_pharmacy_id ON pending_records(pharmacy_id);

-- ============================================================
-- Approved Visits Table
-- ============================================================
CREATE TABLE IF NOT EXISTS approved_visits (
    id VARCHAR(255) PRIMARY KEY,
    pharmacy_id VARCHAR(255) NOT NULL REFERENCES pharmacies(id),
    patient_data JSONB,
    approved_by VARCHAR(255) REFERENCES users(id),
    approved_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_approved_visits_pharmacy_id ON approved_visits(pharmacy_id);

-- ============================================================
-- Audit Logs Table
-- ============================================================
CREATE TABLE IF NOT EXISTS audit_logs (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) REFERENCES users(id),
    pharmacy_id VARCHAR(255) REFERENCES pharmacies(id),
    action VARCHAR(255),
    resource_type VARCHAR(100),
    resource_id VARCHAR(255),
    changes JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_pharmacy_id ON audit_logs(pharmacy_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at);

-- ============================================================
-- Insert sample data
-- ============================================================

-- Sample pharmacy
INSERT INTO pharmacies (id, name, region, district, contact_phone, is_active)
VALUES ('pharm_1', 'Demo Pharmacy', 'Region 1', 'District 1', '+1234567890', true)
ON CONFLICT DO NOTHING;

-- Sample user
INSERT INTO users (id, pharmacy_id, name, email, password_hash, role, is_active)
VALUES (
    'user_1',
    'pharm_1',
    'Demo User',
    'demo@example.com',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36P4/1Pq',
    'admin',
    true
)
ON CONFLICT DO NOTHING;
