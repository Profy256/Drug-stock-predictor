# MED PREDICT SYSTEM
### AI-Powered Hospital Pharmacy Stock Management

A full-stack web application for hospital pharmacies featuring role-based access control, an AI analytics dashboard, predictive stockout alerts, and automated WhatsApp expiry notifications.

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                     REACT FRONTEND                          │
│  Role Select → DE Portal | Admin Portal | DHO Portal        │
│  localStorage: role, token, user                            │
│  sessionStorage: DE session records (wiped on submit)       │
└──────────────────────┬──────────────────────────────────────┘
                       │ REST API (HTTPS / JSON)
┌──────────────────────▼──────────────────────────────────────┐
│                  NODE.JS / EXPRESS API                       │
│  JWT Auth · RBAC · Rate limiting · Audit logging            │
│  Routes: auth, stock, batches, analytics, admin, dho        │
│  Cron: expiry check (07:00) · cache refresh (02:00)         │
└───────────┬──────────────────────────┬──────────────────────┘
            │                          │
┌───────────▼──────────┐  ┌───────────▼──────────────────────┐
│  PostgreSQL (Cloud)   │  │  External Services               │
│  · medicines          │  │  · WhatsApp (Twilio / Meta)      │
│  · data_batches       │  │  · AI LLM (Anthropic / OpenAI)   │
│  · pending_visits     │  │  · Email (Mailgun / SMTP)        │
│  · approved_visits    │  └──────────────────────────────────┘
│  · audit_logs         │
│  JSONB: patient_data  │
└───────────────────────┘
```

---

## Role System & localStorage Persistence

| Role | Key | Portal Color | Capabilities |
|------|-----|-------------|--------------|
| `data_entrant` | DE | Teal `#00BFA6` | Stock entry, patient forms, patient search |
| `admin` | MGR | Indigo `#7C6FFF` | Approval queue, AI dashboard, user mgmt |
| `dho` | DHO | Amber `#FFB547` | Regional map, stockout risk, pharmacy registration |

**How persistence works:**
1. On first visit → role selection screen
2. Chosen role saved to `localStorage['med_predict_role']`
3. On return → skips role selection, shows correct login portal
4. After login → `token` and `user` saved to localStorage
5. Session auto-expires after 24h (JWT expiry timer)
6. "Switch portal" → clears all keys, returns to role selection

---

## Quick Start

### Prerequisites
- Node.js 20+
- PostgreSQL 15+
- Docker & Docker Compose (optional but recommended)

### Option A — Docker (recommended)

```bash
git clone <repo>
cd med-predict

# Copy and edit environment variables
cp backend/.env.example backend/.env
# Fill in: JWT_SECRET, CREDENTIAL_SALT, DB_PASSWORD, and API keys

# Start everything
docker-compose up -d

# Apply database schema (runs automatically via initdb)
# Verify at: http://localhost:5173
```

### Option B — Manual setup

```bash
# 1. Database
createdb med_predict
psql med_predict < backend/db/schema.sql

# 2. Backend
cd backend
cp .env.example .env          # fill in your values
npm install
npm run dev                   # runs on :4000

# 3. Frontend (new terminal)
cd frontend
npm install
npm run dev                   # runs on :5173
```

---

## API Reference

### Authentication

| Method | Endpoint | Auth | Body | Response |
|--------|----------|------|------|----------|
| POST | `/api/v1/auth/login` | — | `{ email, password }` | `{ token, user }` |
| POST | `/api/v1/auth/register` | Admin/DHO | `{ name, email, password, role }` | `{ user }` |
| GET | `/api/v1/auth/me` | Any | — | `{ user }` |

### Stock Management

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/v1/stock` | Any | List all medicines (with status badges) |
| GET | `/api/v1/stock/search?q=` | Any | Typeahead for dispensing form |
| GET | `/api/v1/stock/expiring` | Admin/DHO | Medicines within alert threshold |
| POST | `/api/v1/stock` | DE/Admin | Add incoming stock batch |
| PUT | `/api/v1/stock/:id` | Admin | Manual quantity adjustment (with reason) |

### Batch Approval Queue (Core Workflow)

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/api/v1/batches` | DE | Submit daily records for approval |
| GET | `/api/v1/batches?status=pending` | Admin | View approval queue |
| GET | `/api/v1/batches/:id` | Admin | Batch detail with patient records |
| **POST** | **`/api/v1/batches/:id/approve`** | **Admin** | **Approve: strips PII, commits data, wipes pending** |
| POST | `/api/v1/batches/:id/reject` | Admin | Reject with reason |

### Analytics & AI

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/v1/analytics/trends?from=&to=` | Admin | Top medicines, diseases, visit trend |
| GET | `/api/v1/analytics/ai-summary?from=&to=` | Admin | AI-generated briefing (Claude/GPT) |
| GET | `/api/v1/analytics/stockout-risk` | Admin/DHO | Statistical stockout predictions |
| GET | `/api/v1/analytics/regional` | DHO | Cross-pharmacy overview |

### Admin Configuration

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET/POST | `/api/v1/admin/form-fields` | Admin | Manage dynamic patient form fields |
| PUT/DELETE | `/api/v1/admin/form-fields/:id` | Admin | Update or remove a field |
| GET/PUT | `/api/v1/admin/users/:id` | Admin | List and manage DE accounts |
| GET | `/api/v1/admin/audit-logs` | Admin | Immutable activity log |

### DHO (Top Administrator)

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/v1/dho/pharmacies` | DHO | All registered pharmacies with stats |
| POST | `/api/v1/dho/pharmacies` | DHO | Register new pharmacy + seed admin account |
| GET | `/api/v1/dho/regional-map` | DHO | GPS map data with risk level overlay |

---

## Data Privacy Architecture

```
DE enters data (with patient PII)
         ↓
pending_visit_records  ← Full patient_data JSONB (name, ID, etc.)
         ↓
Manager reviews & APPROVES
         ↓
approved_visits  ← Only: hashed_credential + medicines_given + diagnosis
                          (SHA-256 hash of credential — one-way, unlinkable)
         ↓
pending_visit_records  ← DELETED (PII permanently wiped)
```

Patient PII **never** reaches the central analytics database.
The DE's sessionStorage is cleared immediately after submission.

---

## AI Analytics Logic

The analytics engine in `services/analyticsService.js` works in three layers:

**Layer 1 — Raw SQL aggregation** (always runs, no API key needed)
- Counts medicine dispensing frequency using `jsonb_array_elements(medicines_given)`
- Groups diagnoses, calculates period-over-period change
- Computes stock status snapshot

**Layer 2 — AI summary** (requires API key)
- Sends structured JSON payload to Claude (primary) or GPT-4o-mini (fallback)
- Prompt instructs: trend analysis, procurement priorities, stock alerts, forecast
- Result cached for 12h in `analytics_cache` table

**Layer 3 — Stockout risk prediction** (statistical, no AI key needed)
- Calculates `avg_daily_use` over last 30 days per medicine per pharmacy
- Divides `quantity_remaining` by `avg_daily_use` → `estimated_days_remaining`
- Classifies: `critical` (<7d), `warning` (<14d), `expiring`, `low_stock`, `ok`

---

## WhatsApp Notification Flow

```
Daily cron (07:00 EAT)
    ↓
Query: medicines WHERE expiry_date <= NOW + notification_days
   AND no 'whatsapp/sent' notification_log entry for today
    ↓
Build message with 🔴🟡🟢 urgency indicators
    ↓
sendViaTwilio() or sendViaMeta()
    ↓
[SUCCESS] → Log 'sent' in notification_log
[FAILURE] → Log 'failed' → retry queue (every 6h) → email fallback
```

Admin can configure:
- Global `notification_days` threshold (default: 14)
- Per-medicine overrides in the `medicines.notification_days` column
- WhatsApp number via Admin Dashboard → Notifications tab

---

## Database Key Design Decisions

**JSONB for patient data** (`pending_visit_records.patient_data`)
- Enables Admin to add/remove form fields without schema migrations
- GIN index for fast JSON key search
- Entire column deleted when batch is approved (PII wipe)

**JSONB for medicines dispensed** (`approved_visits.medicines_given`)
- Array of `{ medicine_id, name, qty }` — survives medicine name changes
- `jsonb_array_elements()` unwinding enables medicine frequency aggregation
- GIN indexed for query performance

**Separate pending vs approved tables**
- Hard architectural separation prevents PII leaking into analytics queries
- Makes the "wipe on approval" operation a single `DELETE` statement

---

## Environment Variables Reference

| Variable | Required | Description |
|----------|----------|-------------|
| `JWT_SECRET` | ✅ | Min 64 chars random string |
| `CREDENTIAL_SALT` | ✅ | Salt for patient credential hashing |
| `DB_*` | ✅ | PostgreSQL connection |
| `ANTHROPIC_API_KEY` | Recommended | For AI summaries (Claude) |
| `OPENAI_API_KEY` | Optional | Fallback AI provider |
| `TWILIO_*` | For alerts | WhatsApp via Twilio |
| `META_*` | For alerts | WhatsApp via Meta Cloud API |
| `TZ` | Recommended | Timezone for cron (e.g. `Africa/Kampala`) |

---

## Project Structure

```
med-predict/
├── backend/
│   ├── db/
│   │   ├── schema.sql          PostgreSQL schema (JSONB patient data)
│   │   └── pool.js             Connection pool + transaction helper
│   ├── middleware/
│   │   └── auth.js             JWT verify + requireRole() factory
│   ├── routes/
│   │   ├── auth.js             Login, register, /me
│   │   ├── stock.js            Medicine CRUD + expiry list
│   │   ├── batches.js          ★ Approval queue (core workflow)
│   │   ├── patients.js         DE patient credential search
│   │   ├── analytics.js        Trends + AI summary + stockout risk
│   │   ├── admin.js            Form fields, users, audit logs
│   │   ├── dho.js              Pharmacy registration + regional map
│   │   └── notifications.js    WhatsApp config + test send
│   ├── services/
│   │   ├── analyticsService.js ★ AI engine + SQL aggregations
│   │   ├── cronJobs.js         ★ Expiry checker + WhatsApp sender
│   │   ├── auditService.js     Immutable audit log writer
│   │   └── logger.js           Winston logger
│   ├── server.js               Express app entry point
│   ├── .env.example            All environment variables documented
│   └── Dockerfile
├── frontend/
│   ├── src/
│   │   ├── App.jsx             ★ Complete UI (all 3 portals)
│   │   ├── main.jsx            React DOM entry
│   │   ├── hooks/
│   │   │   └── useRolePersistence.js  ★ Auth context + localStorage logic
│   │   └── services/
│   │       └── api.js          Typed API client (all endpoints)
│   ├── index.html
│   ├── vite.config.js
│   └── package.json
└── docker-compose.yml          Full stack local dev environment
```

---

## Security Checklist

- [x] Passwords: bcrypt, cost factor 12
- [x] JWT: 24h expiry, auto-logout timer in frontend
- [x] RBAC: enforced at both API and frontend routing layer
- [x] Patient PII: never reaches cloud DB, wiped on approval
- [x] Credential hashing: SHA-256 + secret salt (one-way)
- [x] Rate limiting: 10 req/15min on auth, 200 req/min on API
- [x] CORS: restricted to configured frontend URL
- [x] Helmet: security headers
- [x] Input validation: express-validator on all write endpoints
- [x] Audit log: every login, batch action, stock change recorded
- [x] SQL injection: parameterized queries only (pg driver)
- [ ] TODO: Add 2FA for Admin/DHO (TOTP recommended)
- [ ] TODO: Add HTTPS termination at reverse proxy (nginx/Caddy)
