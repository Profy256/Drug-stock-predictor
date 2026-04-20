-- ============================================================
-- MED PREDICT SYSTEM — PostgreSQL Schema
-- Go Backend Migration
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
  diagnosis         VARCHAR(255),
  patient_data      JSONB,
  created_at        TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_pending_batch ON pending_records(batch_id);

-- ─── APPROVED VISITS (Anonymized, After Approval) ────────────
CREATE TABLE approved_visits (
  id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  pharmacy_id       UUID REFERENCES pharmacies(id) ON DELETE CASCADE,
  medicine_id       UUID REFERENCES medicines(id),
  quantity_dispensed INT NOT NULL,
  diagnosis         VARCHAR(255),
  patient_data      JSONB,
  visit_date        TIMESTAMPTZ NOT NULL,
  approved_at       TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_approved_pharmacy ON approved_visits(pharmacy_id);
CREATE INDEX idx_approved_medicine ON approved_visits(medicine_id);
CREATE INDEX idx_approved_date     ON approved_visits(visit_date);

-- ─── ANALYTICS CACHE ────────────────────────────────────────
CREATE TABLE analytics_cache (
  id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  pharmacy_id   UUID REFERENCES pharmacies(id) ON DELETE CASCADE,
  cache_type    VARCHAR(100) NOT NULL,
  data          JSONB,
  expires_at    TIMESTAMPTZ NOT NULL,
  created_at    TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_analytics_pharmacy ON analytics_cache(pharmacy_id);
CREATE INDEX idx_analytics_expires  ON analytics_cache(expires_at);

-- ─── AUDIT LOGS (Compliance & Tracking) ─────────────────────
CREATE TABLE audit_logs (
  id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id       UUID REFERENCES users(id),
  pharmacy_id   UUID REFERENCES pharmacies(id) ON DELETE CASCADE,
  action        VARCHAR(100) NOT NULL,
  entity_type   VARCHAR(50),
  entity_id     UUID,
  details       JSONB,
  ip_address    VARCHAR(45),
  created_at    TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_audit_user      ON audit_logs(user_id);
CREATE INDEX idx_audit_pharmacy  ON audit_logs(pharmacy_id);
CREATE INDEX idx_audit_action    ON audit_logs(action);
CREATE INDEX idx_audit_created   ON audit_logs(created_at);

-- ─── NOTIFICATION LOGS ──────────────────────────────────────
CREATE TABLE notification_logs (
  id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  pharmacy_id   UUID REFERENCES pharmacies(id) ON DELETE CASCADE,
  type          VARCHAR(100) NOT NULL,
  channel       VARCHAR(50),
  recipient     VARCHAR(255),
  message       TEXT,
  status        VARCHAR(50),
  created_at    TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_notification_pharmacy ON notification_logs(pharmacy_id);
CREATE INDEX idx_notification_created  ON notification_logs(created_at);
