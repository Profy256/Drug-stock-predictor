-- ============================================================
-- MED PREDICT SYSTEM — PostgreSQL Schema
-- Python Backend Migration
-- ============================================================

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- ─── ROLES & USERS ──────────────────────────────────────────
CREATE TYPE user_role AS ENUM ('data_entrant', 'admin', 'dho');

CREATE TABLE pharmacies (
  id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name          VARCHAR(255) NOT NULL,
  region        VARCHAR(100),
  district      VARCHAR(100),
  lat           DECIMAL(10,7),
  lng           DECIMAL(10,7),
  contact_phone VARCHAR(30),
  whatsapp_number VARCHAR(30),
  is_active     BOOLEAN DEFAULT true,
  created_at    TIMESTAMPTZ DEFAULT NOW(),
  updated_at    TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE users (
  id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  pharmacy_id   UUID REFERENCES pharmacies(id) ON DELETE CASCADE,
  name          VARCHAR(255) NOT NULL,
  email         VARCHAR(255) UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  role          user_role NOT NULL,
  is_active     BOOLEAN DEFAULT true,
  created_at    TIMESTAMPTZ DEFAULT NOW(),
  updated_at    TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_pharmacy ON users(pharmacy_id);

-- ─── DYNAMIC FORM FIELDS CONFIG ─────────────────────────────
CREATE TABLE patient_form_fields (
  id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  pharmacy_id   UUID REFERENCES pharmacies(id) ON DELETE CASCADE,
  field_key     VARCHAR(100) NOT NULL,
  label         VARCHAR(200) NOT NULL,
  field_type    VARCHAR(50)  NOT NULL DEFAULT 'text',
  options       TEXT[],
  is_required   BOOLEAN DEFAULT false,
  is_active     BOOLEAN DEFAULT true,
  sort_order    INT DEFAULT 0,
  created_at    TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(pharmacy_id, field_key)
);

-- ─── MEDICINES / STOCK ───────────────────────────────────────
CREATE TABLE medicines (
  id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  pharmacy_id       UUID REFERENCES pharmacies(id) ON DELETE CASCADE,
  name              VARCHAR(255) NOT NULL,
  generic_name      VARCHAR(255),
  category          VARCHAR(100),
  unit              VARCHAR(50) DEFAULT 'boxes',
  quantity_total    INT NOT NULL DEFAULT 0,
  quantity_remaining INT NOT NULL DEFAULT 0,
  expiry_date       DATE NOT NULL,
  batch_number      VARCHAR(100),
  supplier          VARCHAR(255),
  unit_cost         DECIMAL(12,2),
  reorder_level     INT DEFAULT 10,
  notification_days INT DEFAULT 14,
  created_by        UUID REFERENCES users(id),
  created_at        TIMESTAMPTZ DEFAULT NOW(),
  updated_at        TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_medicines_pharmacy    ON medicines(pharmacy_id);
CREATE INDEX idx_medicines_expiry      ON medicines(expiry_date);
CREATE INDEX idx_medicines_name_trgm   ON medicines USING gin(name gin_trgm_ops);

-- ─── BATCHES (Daily Data Submissions) ────────────────────────
CREATE TABLE batches (
  id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  pharmacy_id       UUID REFERENCES pharmacies(id) ON DELETE CASCADE,
  submitted_by      UUID REFERENCES users(id),
  status            VARCHAR(50) NOT NULL DEFAULT 'pending',
  rejection_reason  TEXT,
  approved_by       UUID REFERENCES users(id),
  record_count      INT DEFAULT 0,
  created_at        TIMESTAMPTZ DEFAULT NOW(),
  updated_at        TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_batches_pharmacy ON batches(pharmacy_id);
CREATE INDEX idx_batches_status   ON batches(status);

-- ─── PENDING RECORDS (Before Approval) ───────────────────────
CREATE TABLE pending_records (
  id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  batch_id          UUID REFERENCES batches(id) ON DELETE CASCADE,
  patient_hash      VARCHAR(255) NOT NULL,
  medicine_id       UUID REFERENCES medicines(id),
  quantity_dispensed INT NOT NULL,
  diagnosis         TEXT,
  patient_data      JSONB DEFAULT '{}',
  created_at        TIMESTAMPTZ DEFAULT NOW()
);

-- ─── APPROVED VISITS (After Batch Approval) ──────────────────
CREATE TABLE approved_visits (
  id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  pharmacy_id       UUID REFERENCES pharmacies(id) ON DELETE CASCADE,
  medicine_id       UUID REFERENCES medicines(id),
  quantity_dispensed INT NOT NULL,
  diagnosis         TEXT,
  patient_data      JSONB DEFAULT '{}',
  visit_date        TIMESTAMPTZ,
  approved_at       TIMESTAMPTZ DEFAULT NOW()
);

-- ─── NOTIFICATIONS ──────────────────────────────────────────
CREATE TABLE notification_logs (
  id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  pharmacy_id UUID REFERENCES pharmacies(id) ON DELETE CASCADE,
  type       VARCHAR(100) NOT NULL,
  channel    VARCHAR(50) NOT NULL,
  recipient  VARCHAR(255) NOT NULL,
  message    TEXT NOT NULL,
  status     VARCHAR(50) DEFAULT 'pending',
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- ─── AUDIT LOG ──────────────────────────────────────────────
CREATE TABLE audit_logs (
  id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id       UUID REFERENCES users(id),
  pharmacy_id   UUID REFERENCES pharmacies(id),
  action        VARCHAR(100) NOT NULL,
  entity_type   VARCHAR(100) NOT NULL,
  entity_id     VARCHAR(255),
  ip_address    VARCHAR(45),
  changes       JSONB,
  created_at    TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_audit_user     ON audit_logs(user_id);
CREATE INDEX idx_audit_pharmacy ON audit_logs(pharmacy_id);
CREATE INDEX idx_audit_action   ON audit_logs(action);
