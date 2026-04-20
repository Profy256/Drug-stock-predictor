package db

import (
	"database/sql"
	"errors"
	"time"

	"med-predict-backend/internal/models"
)

// ============================================================
// User Queries
// ============================================================

func (d *Database) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := d.DB.QueryRow(
		`SELECT id, pharmacy_id, name, email, password_hash, role, is_active, created_at, updated_at
		 FROM users WHERE email = $1 AND is_active = true`,
		email,
	).Scan(&user.ID, &user.PharmacyID, &user.Name, &user.Email, &user.PasswordHash,
		&user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (d *Database) GetUserByID(userID string) (*models.User, error) {
	user := &models.User{}
	err := d.DB.QueryRow(
		`SELECT id, pharmacy_id, name, email, password_hash, role, is_active, created_at, updated_at
		 FROM users WHERE id = $1`,
		userID,
	).Scan(&user.ID, &user.PharmacyID, &user.Name, &user.Email, &user.PasswordHash,
		&user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (d *Database) CreateUser(user *models.User) error {
	err := d.DB.QueryRow(
		`INSERT INTO users (id, pharmacy_id, name, email, password_hash, role, is_active, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		 RETURNING id, created_at, updated_at`,
		user.ID, user.PharmacyID, user.Name, user.Email, user.PasswordHash,
		user.Role, user.IsActive, time.Now(), time.Now(),
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	return err
}

func (d *Database) GetPharmacyUsers(pharmacyID string) ([]models.User, error) {
	rows, err := d.DB.Query(
		`SELECT id, pharmacy_id, name, email, password_hash, role, is_active, created_at, updated_at
		 FROM users WHERE pharmacy_id = $1 ORDER BY created_at DESC`,
		pharmacyID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.PharmacyID, &user.Name, &user.Email, &user.PasswordHash,
			&user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

// ============================================================
// Pharmacy Queries
// ============================================================

func (d *Database) GetPharmacy(pharmacyID string) (*models.Pharmacy, error) {
	pharmacy := &models.Pharmacy{}
	err := d.DB.QueryRow(
		`SELECT id, name, region, district, lat, lng, contact_phone, whatsapp_number, is_active, created_at, updated_at
		 FROM pharmacies WHERE id = $1`,
		pharmacyID,
	).Scan(&pharmacy.ID, &pharmacy.Name, &pharmacy.Region, &pharmacy.District, &pharmacy.Lat, &pharmacy.Lng,
		&pharmacy.ContactPhone, &pharmacy.WhatsAppNum, &pharmacy.IsActive, &pharmacy.CreatedAt, &pharmacy.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("pharmacy not found")
		}
		return nil, err
	}
	return pharmacy, nil
}

func (d *Database) CreatePharmacy(pharmacy *models.Pharmacy) error {
	pharmacy.CreatedAt = time.Now()
	pharmacy.UpdatedAt = time.Now()
	err := d.DB.QueryRow(
		`INSERT INTO pharmacies (id, name, region, district, lat, lng, contact_phone, whatsapp_number, is_active, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		 RETURNING id, created_at, updated_at`,
		pharmacy.ID, pharmacy.Name, pharmacy.Region, pharmacy.District, pharmacy.Lat, pharmacy.Lng,
		pharmacy.ContactPhone, pharmacy.WhatsAppNum, pharmacy.IsActive, pharmacy.CreatedAt, pharmacy.UpdatedAt,
	).Scan(&pharmacy.ID, &pharmacy.CreatedAt, &pharmacy.UpdatedAt)
	return err
}

func (d *Database) GetAllPharmacies() ([]models.Pharmacy, error) {
	rows, err := d.DB.Query(
		`SELECT id, name, region, district, lat, lng, contact_phone, whatsapp_number, is_active, created_at, updated_at
		 FROM pharmacies WHERE is_active = true ORDER BY name`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pharmacies []models.Pharmacy
	for rows.Next() {
		var p models.Pharmacy
		err := rows.Scan(&p.ID, &p.Name, &p.Region, &p.District, &p.Lat, &p.Lng,
			&p.ContactPhone, &p.WhatsAppNum, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		pharmacies = append(pharmacies, p)
	}
	return pharmacies, rows.Err()
}

// ============================================================
// Medicine Queries
// ============================================================

func (d *Database) CreateMedicine(med *models.Medicine) error {
	med.CreatedAt = time.Now()
	med.UpdatedAt = time.Now()
	err := d.DB.QueryRow(
		`INSERT INTO medicines (id, pharmacy_id, name, generic_name, category, unit, quantity_total, 
		 quantity_remaining, expiry_date, batch_number, supplier, unit_cost, reorder_level, notification_days, created_by, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
		 RETURNING id, created_at, updated_at`,
		med.ID, med.PharmacyID, med.Name, med.GenericName, med.Category, med.Unit,
		med.QuantityTotal, med.QuantityRemaining, med.ExpiryDate, med.BatchNumber,
		med.Supplier, med.UnitCost, med.ReorderLevel, med.NotificationDays,
		med.CreatedBy, med.CreatedAt, med.UpdatedAt,
	).Scan(&med.ID, &med.CreatedAt, &med.UpdatedAt)
	return err
}

func (d *Database) GetMedicine(medicineID string) (*models.Medicine, error) {
	med := &models.Medicine{}
	err := d.DB.QueryRow(
		`SELECT id, pharmacy_id, name, generic_name, category, unit, quantity_total, quantity_remaining,
		 expiry_date, batch_number, supplier, unit_cost, reorder_level, notification_days, created_by, created_at, updated_at
		 FROM medicines WHERE id = $1`,
		medicineID,
	).Scan(&med.ID, &med.PharmacyID, &med.Name, &med.GenericName, &med.Category, &med.Unit,
		&med.QuantityTotal, &med.QuantityRemaining, &med.ExpiryDate, &med.BatchNumber,
		&med.Supplier, &med.UnitCost, &med.ReorderLevel, &med.NotificationDays,
		&med.CreatedBy, &med.CreatedAt, &med.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("medicine not found")
		}
		return nil, err
	}
	return med, nil
}

func (d *Database) GetPharmacyMedicines(pharmacyID string) ([]models.Medicine, error) {
	rows, err := d.DB.Query(
		`SELECT id, pharmacy_id, name, generic_name, category, unit, quantity_total, quantity_remaining,
		 expiry_date, batch_number, supplier, unit_cost, reorder_level, notification_days, created_by, created_at, updated_at
		 FROM medicines WHERE pharmacy_id = $1 ORDER BY name`,
		pharmacyID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var medicines []models.Medicine
	for rows.Next() {
		var med models.Medicine
		err := rows.Scan(&med.ID, &med.PharmacyID, &med.Name, &med.GenericName, &med.Category, &med.Unit,
			&med.QuantityTotal, &med.QuantityRemaining, &med.ExpiryDate, &med.BatchNumber,
			&med.Supplier, &med.UnitCost, &med.ReorderLevel, &med.NotificationDays,
			&med.CreatedBy, &med.CreatedAt, &med.UpdatedAt)
		if err != nil {
			return nil, err
		}
		medicines = append(medicines, med)
	}
	return medicines, rows.Err()
}

func (d *Database) SearchMedicines(pharmacyID, query string) ([]models.Medicine, error) {
	rows, err := d.DB.Query(
		`SELECT id, pharmacy_id, name, generic_name, category, unit, quantity_total, quantity_remaining,
		 expiry_date, batch_number, supplier, unit_cost, reorder_level, notification_days, created_by, created_at, updated_at
		 FROM medicines WHERE pharmacy_id = $1 AND (name ILIKE $2 OR generic_name ILIKE $2)
		 ORDER BY name LIMIT 20`,
		pharmacyID, "%"+query+"%",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var medicines []models.Medicine
	for rows.Next() {
		var med models.Medicine
		err := rows.Scan(&med.ID, &med.PharmacyID, &med.Name, &med.GenericName, &med.Category, &med.Unit,
			&med.QuantityTotal, &med.QuantityRemaining, &med.ExpiryDate, &med.BatchNumber,
			&med.Supplier, &med.UnitCost, &med.ReorderLevel, &med.NotificationDays,
			&med.CreatedBy, &med.CreatedAt, &med.UpdatedAt)
		if err != nil {
			return nil, err
		}
		medicines = append(medicines, med)
	}
	return medicines, rows.Err()
}

func (d *Database) UpdateMedicineQuantity(medicineID string, quantityAdjustment int) error {
	result, err := d.DB.Exec(
		`UPDATE medicines SET quantity_remaining = quantity_remaining + $1, updated_at = $2
		 WHERE id = $3`,
		quantityAdjustment, time.Now(), medicineID,
	)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("medicine not found")
	}
	return nil
}

func (d *Database) GetExpiringMedicines(pharmacyID string, daysThreshold int) ([]models.Medicine, error) {
	rows, err := d.DB.Query(
		`SELECT id, pharmacy_id, name, generic_name, category, unit, quantity_total, quantity_remaining,
		 expiry_date, batch_number, supplier, unit_cost, reorder_level, notification_days, created_by, created_at, updated_at
		 FROM medicines WHERE pharmacy_id = $1 
		 AND expiry_date <= NOW() + INTERVAL '1 day' * $2 
		 AND expiry_date > NOW()
		 ORDER BY expiry_date ASC`,
		pharmacyID, daysThreshold,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var medicines []models.Medicine
	for rows.Next() {
		var med models.Medicine
		err := rows.Scan(&med.ID, &med.PharmacyID, &med.Name, &med.GenericName, &med.Category, &med.Unit,
			&med.QuantityTotal, &med.QuantityRemaining, &med.ExpiryDate, &med.BatchNumber,
			&med.Supplier, &med.UnitCost, &med.ReorderLevel, &med.NotificationDays,
			&med.CreatedBy, &med.CreatedAt, &med.UpdatedAt)
		if err != nil {
			return nil, err
		}
		medicines = append(medicines, med)
	}
	return medicines, rows.Err()
}

// ============================================================
// Batch Queries
// ============================================================

func (d *Database) CreateBatch(batch *models.Batch) error {
	batch.CreatedAt = time.Now()
	batch.UpdatedAt = time.Now()
	err := d.DB.QueryRow(
		`INSERT INTO batches (id, pharmacy_id, submitted_by, status, record_count, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, created_at, updated_at`,
		batch.ID, batch.PharmacyID, batch.SubmittedBy, batch.Status, batch.RecordCount, batch.CreatedAt, batch.UpdatedAt,
	).Scan(&batch.ID, &batch.CreatedAt, &batch.UpdatedAt)
	return err
}

func (d *Database) GetBatch(batchID string) (*models.Batch, error) {
	batch := &models.Batch{}
	err := d.DB.QueryRow(
		`SELECT id, pharmacy_id, submitted_by, status, rejection_reason, approved_by, record_count, created_at, updated_at
		 FROM batches WHERE id = $1`,
		batchID,
	).Scan(&batch.ID, &batch.PharmacyID, &batch.SubmittedBy, &batch.Status, &batch.RejectionReason,
		&batch.ApprovedBy, &batch.RecordCount, &batch.CreatedAt, &batch.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("batch not found")
		}
		return nil, err
	}
	return batch, nil
}

// GetPharmacyBatches returns batches with pagination (most recent first)
func (d *Database) GetPharmacyBatches(pharmacyID string, limit, offset int) ([]models.Batch, error) {
	rows, err := d.DB.Query(
		`SELECT id, pharmacy_id, submitted_by, status, rejection_reason, approved_by, record_count, created_at, updated_at
		 FROM batches WHERE pharmacy_id = $1
		 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		pharmacyID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var batches []models.Batch
	for rows.Next() {
		var b models.Batch
		err := rows.Scan(&b.ID, &b.PharmacyID, &b.SubmittedBy, &b.Status, &b.RejectionReason,
			&b.ApprovedBy, &b.RecordCount, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, err
		}
		batches = append(batches, b)
	}
	return batches, rows.Err()
}

// UpdateBatchStatus updates batch status and approval info
func (d *Database) UpdateBatchStatus(batchID, status, reason, approvedBy string) error {
	_, err := d.DB.Exec(
		`UPDATE batches SET status = $1, rejection_reason = $2, approved_by = $3, updated_at = $4
		 WHERE id = $5`,
		status, reason, approvedBy, time.Now(), batchID,
	)
	return err
}

// ============================================================
// Pending Record Queries
// ============================================================

func (d *Database) CreatePendingRecord(record *models.PendingRecord) error {
	record.CreatedAt = time.Now()
	err := d.DB.QueryRow(
		`INSERT INTO pending_records (id, batch_id, patient_hash, medicine_id, quantity_dispensed, diagnosis, patient_data, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING id, created_at`,
		record.ID, record.BatchID, record.PatientHash, record.MedicineID,
		record.QuantityDispensed, record.Diagnosis, record.PatientData, record.CreatedAt,
	).Scan(&record.ID, &record.CreatedAt)
	return err
}

func (d *Database) GetBatchRecords(batchID string) ([]models.PendingRecord, error) {
	rows, err := d.DB.Query(
		`SELECT id, batch_id, patient_hash, medicine_id, quantity_dispensed, diagnosis, patient_data, created_at
		 FROM pending_records WHERE batch_id = $1`,
		batchID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []models.PendingRecord
	for rows.Next() {
		var r models.PendingRecord
		err := rows.Scan(&r.ID, &r.BatchID, &r.PatientHash, &r.MedicineID,
			&r.QuantityDispensed, &r.Diagnosis, &r.PatientData, &r.CreatedAt)
		if err != nil {
			return nil, err
		}
		records = append(records, r)
	}
	return records, rows.Err()
}

func (d *Database) DeletePendingRecord(recordID string) error {
	result, err := d.DB.Exec(`DELETE FROM pending_records WHERE id = $1`, recordID)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("record not found")
	}
	return nil
}

// ============================================================
// Approved Visit Queries
// ============================================================

func (d *Database) CreateApprovedVisit(visit *models.ApprovedVisit) error {
	visit.ApprovedAt = time.Now()
	err := d.DB.QueryRow(
		`INSERT INTO approved_visits (id, pharmacy_id, medicine_id, quantity_dispensed, diagnosis, patient_data, visit_date, approved_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING id, approved_at`,
		visit.ID, visit.PharmacyID, visit.MedicineID, visit.QuantityDispensed,
		visit.Diagnosis, visit.PatientData, visit.VisitDate, visit.ApprovedAt,
	).Scan(&visit.ID, &visit.ApprovedAt)
	return err
}

// GetPharmacyApprovedVisits returns approved visits with optional date filtering
func (d *Database) GetPharmacyApprovedVisits(pharmacyID string, startDate, endDate time.Time) ([]models.ApprovedVisit, error) {
	rows, err := d.DB.Query(
		`SELECT id, pharmacy_id, medicine_id, quantity_dispensed, diagnosis, patient_data, visit_date, approved_at
		 FROM approved_visits WHERE pharmacy_id = $1 AND visit_date >= $2 AND visit_date < $3
		 ORDER BY visit_date DESC`,
		pharmacyID, startDate, endDate,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var visits []models.ApprovedVisit
	for rows.Next() {
		var v models.ApprovedVisit
		err := rows.Scan(&v.ID, &v.PharmacyID, &v.MedicineID, &v.QuantityDispensed,
			&v.Diagnosis, &v.PatientData, &v.VisitDate, &v.ApprovedAt)
		if err != nil {
			return nil, err
		}
		visits = append(visits, v)
	}
	return visits, rows.Err()
}

// ============================================================
// Audit Log Queries
// ============================================================

func (d *Database) LogAuditEvent(log *models.AuditLog) error {
	log.CreatedAt = time.Now()
	_, err := d.DB.Exec(
		`INSERT INTO audit_logs (id, user_id, pharmacy_id, action, entity_type, entity_id, details, ip_address, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		log.ID, log.UserID, log.PharmacyID, log.Action, log.EntityType, log.EntityID, log.Details, log.IPAddress, log.CreatedAt,
	)
	return err
}

func (d *Database) GetAuditLogs(pharmacyID string, limit, offset int) ([]models.AuditLog, error) {
	rows, err := d.DB.Query(
		`SELECT id, user_id, pharmacy_id, action, entity_type, entity_id, details, ip_address, created_at
		 FROM audit_logs WHERE pharmacy_id = $1
		 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		pharmacyID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.AuditLog
	for rows.Next() {
		var log models.AuditLog
		err := rows.Scan(&log.ID, &log.UserID, &log.PharmacyID, &log.Action, &log.EntityType,
			&log.EntityID, &log.Details, &log.IPAddress, &log.CreatedAt)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, rows.Err()
}

// ============================================================
// Notification Log Queries
// ============================================================

func (d *Database) LogNotification(log *models.NotificationLog) error {
	log.CreatedAt = time.Now()
	err := d.DB.QueryRow(
		`INSERT INTO notification_logs (id, pharmacy_id, type, channel, recipient, message, status, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING id`,
		log.ID, log.PharmacyID, log.Type, log.Channel, log.Recipient, log.Message, log.Status, log.CreatedAt,
	).Scan(&log.ID)
	return err
}

// ============================================================
// Patient Form Fields Queries
// ============================================================

func (d *Database) GetPharmacyFormFields(pharmacyID string) ([]models.PatientFormField, error) {
	rows, err := d.DB.Query(
		`SELECT id, pharmacy_id, field_key, label, field_type, options, is_required, is_active, sort_order, created_at
		 FROM patient_form_fields WHERE pharmacy_id = $1 AND is_active = true
		 ORDER BY sort_order ASC`,
		pharmacyID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fields []models.PatientFormField
	for rows.Next() {
		var f models.PatientFormField
		err := rows.Scan(&f.ID, &f.PharmacyID, &f.FieldKey, &f.Label, &f.FieldType,
			&f.Options, &f.IsRequired, &f.IsActive, &f.SortOrder, &f.CreatedAt)
		if err != nil {
			return nil, err
		}
		fields = append(fields, f)
	}
	return fields, rows.Err()
}

func (d *Database) CreateFormField(field *models.PatientFormField) error {
	field.CreatedAt = time.Now()
	err := d.DB.QueryRow(
		`INSERT INTO patient_form_fields (id, pharmacy_id, field_key, label, field_type, options, is_required, is_active, sort_order, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		 RETURNING id, created_at`,
		field.ID, field.PharmacyID, field.FieldKey, field.Label, field.FieldType,
		field.Options, field.IsRequired, field.IsActive, field.SortOrder, field.CreatedAt,
	).Scan(&field.ID, &field.CreatedAt)
	return err
}
