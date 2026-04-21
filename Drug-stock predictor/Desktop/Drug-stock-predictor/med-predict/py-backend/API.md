# API Documentation

## Overview

Med Predict Backend is a RESTful API for pharmaceutical inventory management. All endpoints require JWT authentication (except login/register).

**Base URL:** `http://localhost:8000/api/v1`

## Authentication

### Login
```
POST /auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}

Response (200 OK):
{
  "token": "eyJ0eXAiOiJKV1QiLCJhbGc...",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "pharmacy_id": "550e8400-e29b-41d4-a716-446655440001",
  "name": "John Doe",
  "email": "user@example.com",
  "role": "admin"
}
```

### Register
```
POST /auth/register
Content-Type: application/json

{
  "pharmacy_id": "550e8400-e29b-41d4-a716-446655440001",
  "name": "Jane Doe",
  "email": "jane@example.com",
  "password": "password123"
}

Response (200 OK):
{
  "id": "550e8400-e29b-41d4-a716-446655440002",
  "pharmacy_id": "550e8400-e29b-41d4-a716-446655440001",
  "name": "Jane Doe",
  "email": "jane@example.com",
  "role": "data_entrant",
  "is_active": true,
  "created_at": "2026-04-21T10:30:00Z",
  "updated_at": "2026-04-21T10:30:00Z"
}
```

### Get Current User
```
GET /auth/me
Authorization: Bearer <TOKEN>

Response (200 OK):
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "pharmacy_id": "550e8400-e29b-41d4-a716-446655440001",
  "name": "John Doe",
  "email": "user@example.com",
  "role": "admin"
}
```

## Stock Management

### List Medicines
```
GET /stock/medicines
Authorization: Bearer <TOKEN>

Query Parameters:
- None

Response (200 OK):
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440003",
    "pharmacy_id": "550e8400-e29b-41d4-a716-446655440001",
    "name": "Paracetamol",
    "generic_name": "Acetaminophen",
    "category": "Analgesic",
    "unit": "boxes",
    "quantity_total": 1000,
    "quantity_remaining": 500,
    "expiry_date": "2026-12-31",
    "batch_number": "BATCH001",
    "supplier": "PharmaCorp",
    "unit_cost": 5.50,
    "reorder_level": 100,
    "notification_days": 14,
    "status": "ok",
    "created_by": "550e8400-e29b-41d4-a716-446655440000",
    "created_at": "2026-04-21T10:30:00Z",
    "updated_at": "2026-04-21T10:30:00Z"
  }
]
```

### Create Medicine
```
POST /stock/medicines
Authorization: Bearer <TOKEN>
Content-Type: application/json

{
  "name": "Paracetamol",
  "generic_name": "Acetaminophen",
  "category": "Analgesic",
  "unit": "boxes",
  "quantity_total": 1000,
  "quantity_remaining": 500,
  "expiry_date": "2026-12-31",
  "batch_number": "BATCH001",
  "supplier": "PharmaCorp",
  "unit_cost": 5.50,
  "reorder_level": 100,
  "notification_days": 14
}

Response (200 OK): [Medicine object]
```

### Get Medicine
```
GET /stock/medicines/{medicine_id}
Authorization: Bearer <TOKEN>

Response (200 OK): [Medicine object]
```

### Update Medicine
```
PUT /stock/medicines/{medicine_id}
Authorization: Bearer <TOKEN>
Content-Type: application/json

{
  "quantity_remaining": 400,
  "status": "low"
}

Response (200 OK): [Medicine object]
```

### Delete Medicine
```
DELETE /stock/medicines/{medicine_id}
Authorization: Bearer <TOKEN>

Response (204 No Content)
```

## Analytics

### Stockout Predictions
```
GET /analytics/stockout-predictions
Authorization: Bearer <TOKEN>

Response (200 OK):
{
  "predictions": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440003",
      "name": "Paracetamol",
      "quantity_remaining": 50,
      "reorder_level": 100,
      "status": "low"
    }
  ]
}
```

### Trends
```
GET /analytics/trends
Authorization: Bearer <TOKEN>

Response (200 OK):
{
  "trends": {
    "550e8400-e29b-41d4-a716-446655440003": 250,
    "550e8400-e29b-41d4-a716-446655440004": 180
  }
}
```

### Expiry Alerts
```
GET /analytics/expiry-alerts
Authorization: Bearer <TOKEN>

Response (200 OK):
{
  "alerts": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440003",
      "name": "Paracetamol",
      "expiry_date": "2026-04-30",
      "days_to_expiry": 9,
      "quantity_remaining": 100
    }
  ]
}
```

## Batch Management

### Create Batch
```
POST /batches
Authorization: Bearer <TOKEN>
Content-Type: application/json

{
  "pharmacy_id": "550e8400-e29b-41d4-a716-446655440001"
}

Response (200 OK): [Batch object]
```

### List Batches
```
GET /batches
Authorization: Bearer <TOKEN>

Response (200 OK): [Batch array]
```

### Get Batch
```
GET /batches/{batch_id}
Authorization: Bearer <TOKEN>

Response (200 OK): [Batch object]
```

## Patient Management

### List Form Fields
```
GET /patient/form-fields
Authorization: Bearer <TOKEN>

Response (200 OK): [FormField array]
```

### Create Form Field
```
POST /patient/form-fields
Authorization: Bearer <TOKEN>
Content-Type: application/json

{
  "field_key": "age",
  "label": "Age",
  "field_type": "number",
  "options": null,
  "is_required": true,
  "sort_order": 1
}

Response (200 OK): [FormField object]
```

## Records Management

### Create Pending Record
```
POST /records/pending
Authorization: Bearer <TOKEN>
Content-Type: application/json

{
  "batch_id": "550e8400-e29b-41d4-a716-446655440010",
  "patient_hash": "sha256_hash_of_patient_id",
  "medicine_id": "550e8400-e29b-41d4-a716-446655440003",
  "quantity_dispensed": 10,
  "diagnosis": "Headache",
  "patient_data": {
    "age": 35,
    "gender": "M"
  }
}

Response (200 OK): [PendingRecord object]
```

### List Pending Records
```
GET /records/pending/{batch_id}
Authorization: Bearer <TOKEN>

Response (200 OK): [PendingRecord array]
```

### List Approved Visits
```
GET /records/approved
Authorization: Bearer <TOKEN>

Response (200 OK): [ApprovedVisit array]
```

## Admin Operations

### Approve Batch (Admin Only)
```
POST /admin/batches/{batch_id}/approve
Authorization: Bearer <TOKEN>

Response (200 OK):
{
  "status": "approved",
  "records_moved": 25
}
```

### Reject Batch (Admin Only)
```
POST /admin/batches/{batch_id}/reject
Authorization: Bearer <TOKEN>
Content-Type: application/json

{
  "reason": "Incomplete data"
}

Response (200 OK):
{
  "status": "rejected",
  "reason": "Incomplete data"
}
```

### List All Users (Admin Only)
```
GET /admin/users
Authorization: Bearer <TOKEN>

Response (200 OK): [User array]
```

## DHO Operations

### List DHO Batches (DHO Only)
```
GET /dho/batches
Authorization: Bearer <TOKEN>

Response (200 OK): [Batch array]
```

### Get DHO Batch Details (DHO Only)
```
GET /dho/batches/{batch_id}/details
Authorization: Bearer <TOKEN>

Response (200 OK): [Batch object]
```

## Error Responses

### 400 Bad Request
```json
{
  "detail": "Invalid request data"
}
```

### 401 Unauthorized
```json
{
  "detail": "Invalid authentication credentials"
}
```

### 403 Forbidden
```json
{
  "detail": "Only admins can approve batches"
}
```

### 404 Not Found
```json
{
  "detail": "Medicine not found"
}
```

### 500 Internal Server Error
```json
{
  "detail": "Internal server error"
}
```

## Status Codes

| Code | Meaning |
|------|---------|
| 200 | Success |
| 201 | Created |
| 204 | No Content |
| 400 | Bad Request |
| 401 | Unauthorized |
| 403 | Forbidden |
| 404 | Not Found |
| 500 | Internal Server Error |

## Medicine Status Values

- `ok` - Sufficient quantity, not expiring soon
- `low` - Quantity below reorder level
- `expiring` - Expiry date within notification days
- `expired` - Past expiry date

## User Roles

- `data_entrant` - Can create and submit batches
- `admin` - Can approve/reject batches, manage users
- `dho` - Can review all batches from all pharmacies

## Examples

### Complete Login Flow

```bash
# 1. Login
curl -X POST "http://localhost:8000/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'

# Response includes token
# TOKEN=eyJ0eXAiOiJKV1QiLCJhbGc...

# 2. Get Current User
curl -X GET "http://localhost:8000/api/v1/auth/me" \
  -H "Authorization: Bearer $TOKEN"

# 3. List Medicines
curl -X GET "http://localhost:8000/api/v1/stock/medicines" \
  -H "Authorization: Bearer $TOKEN"
```

## API Documentation Interface

Once the server is running, access interactive API documentation:

- **Swagger UI**: http://localhost:8000/docs
- **ReDoc**: http://localhost:8000/redoc
- **OpenAPI JSON**: http://localhost:8000/openapi.json

## Rate Limiting

- Login endpoint: 60 requests per minute per IP
- API endpoints: 60 requests per minute per user

## Testing

Use the Swagger UI at `/docs` to test all endpoints interactively.
