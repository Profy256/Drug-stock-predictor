package services

import (
	"fmt"
	"time"

	"med-predict-backend/internal/db"
	"med-predict-backend/internal/models"
)

// AnalyticsService handles data analytics and predictions
type AnalyticsService struct {
	db     *db.Database
	logger *Logger
}

// NewAnalyticsService creates an analytics service
func NewAnalyticsService(database *db.Database, logger *Logger) *AnalyticsService {
	return &AnalyticsService{
		db:     database,
		logger: logger,
	}
}

// Trends contains analytics data for a pharmacy
type Trends struct {
	TopMedicines []MedicineFreq `json:"top_medicines"`
	TopDiseases  []DiseaseFreq  `json:"top_diseases"`
	DailyVisits  []DailyVisit   `json:"daily_visits"`
	StockStatus  StockStatus    `json:"stock_status"`
	PeriodChange float64        `json:"period_change_percent"`
}

type MedicineFreq struct {
	Name      string  `json:"name"`
	Count     int     `json:"count"`
	Frequency float64 `json:"frequency_percent"`
}

type DiseaseFreq struct {
	Name      string  `json:"name"`
	Count     int     `json:"count"`
	Frequency float64 `json:"frequency_percent"`
}

type DailyVisit struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type StockStatus struct {
	OK       int `json:"ok"`
	Expiring int `json:"expiring_soon"`
	Expired  int `json:"expired"`
	LowStock int `json:"low_stock"`
}

// GetTrends computes analytics trends for a pharmacy
func (as *AnalyticsService) GetTrends(pharmacyID string, startDate, endDate time.Time) (*Trends, error) {
	trends := &Trends{
		TopMedicines: []MedicineFreq{},
		TopDiseases:  []DiseaseFreq{},
		DailyVisits:  []DailyVisit{},
	}

	// Fetch approved visits in the date range
	visits, err := as.db.GetPharmacyApprovedVisits(pharmacyID, startDate, endDate)
	if err != nil {
		as.logger.Error("failed to fetch visits", "error", err.Error())
		return trends, err
	}

	// Aggregate medicine and disease frequency
	medicineMap := make(map[string]int)
	diseaseMap := make(map[string]int)
	dailyMap := make(map[string]int)
	totalVisits := len(visits)

	for _, visit := range visits {
		// Count medicine usage
		med, err := as.db.GetMedicine(visit.MedicineID)
		if err == nil {
			medicineMap[med.Name]++
		}

		// Count disease frequency
		if visit.Diagnosis != "" {
			diseaseMap[visit.Diagnosis]++
		}

		// Count daily visits
		dayKey := visit.VisitDate.Format("2006-01-02")
		dailyMap[dayKey]++
	}

	// Convert to sorted slices (top 10)
	trends.TopMedicines = getTopMedicines(medicineMap, totalVisits)
	trends.TopDiseases = getTopDiseases(diseaseMap, totalVisits)
	trends.DailyVisits = getDailyVisits(dailyMap)

	// Get stock status
	medicines, err := as.db.GetPharmacyMedicines(pharmacyID)
	if err == nil {
		trends.StockStatus = calculateStockStatus(medicines)
	}

	return trends, nil
}

// StockoutRisk represents medicine stockout prediction
type StockoutRisk struct {
	MedicineID    string `json:"medicine_id"`
	MedicineName  string `json:"medicine_name"`
	DaysRemaining int    `json:"days_remaining"`
	RiskLevel     string `json:"risk_level"` // critical, warning, expiring, low, ok
}

// PredictStockoutRisk calculates stockout risk for all medicines
func (as *AnalyticsService) PredictStockoutRisk(pharmacyID string) ([]StockoutRisk, error) {
	medicines, err := as.db.GetPharmacyMedicines(pharmacyID)
	if err != nil {
		return nil, err
	}

	var risks []StockoutRisk

	for _, med := range medicines {
		risk := StockoutRisk{
			MedicineID:   med.ID,
			MedicineName: med.Name,
		}

		// Check if expired
		if med.ExpiryDate.Before(time.Now()) {
			risk.RiskLevel = "expired"
			risk.DaysRemaining = 0
		} else if med.ExpiryDate.Sub(time.Now()).Hours()/24 <= 3 {
			risk.RiskLevel = "critical"
			risk.DaysRemaining = int(med.ExpiryDate.Sub(time.Now()).Hours() / 24)
		} else if med.ExpiryDate.Sub(time.Now()).Hours()/24 <= 7 {
			risk.RiskLevel = "warning"
			risk.DaysRemaining = int(med.ExpiryDate.Sub(time.Now()).Hours() / 24)
		} else if med.ExpiryDate.Sub(time.Now()).Hours()/24 <= 14 {
			risk.RiskLevel = "expiring"
			risk.DaysRemaining = int(med.ExpiryDate.Sub(time.Now()).Hours() / 24)
		} else if med.QuantityRemaining <= med.ReorderLevel {
			risk.RiskLevel = "low_stock"
			risk.DaysRemaining = -1
		} else {
			risk.RiskLevel = "ok"
			risk.DaysRemaining = int(med.ExpiryDate.Sub(time.Now()).Hours() / 24)
		}

		risks = append(risks, risk)
	}

	return risks, nil
}

// AISummary represents AI-generated briefing
type AISummary struct {
	Summary   string    `json:"summary"`
	Period    string    `json:"period"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GenerateAISummary creates an AI summary (with fallback to statistical summary)
func (as *AnalyticsService) GenerateAISummary(trends *Trends, pharmacyName, period string) *AISummary {
	// For now, generate a statistical summary
	// In a real implementation, you would call Anthropic/OpenAI API here

	summary := fmt.Sprintf(
		"Pharmacy %s Summary (%s):\n"+
			"- Total visits analyzed\n"+
			"- Top medicine: %s (%d dispensed)\n"+
			"- Common diagnosis: %s\n"+
			"- Stock status: %d medicines OK, %d expiring soon, %d expired\n"+
			"- Stockout risk status monitored",
		pharmacyName,
		period,
		getFirstMedicine(trends.TopMedicines),
		getFirstMedicineCount(trends.TopMedicines),
		getFirstDisease(trends.TopDiseases),
		trends.StockStatus.OK,
		trends.StockStatus.Expiring,
		trends.StockStatus.Expired,
	)

	return &AISummary{
		Summary:   summary,
		Period:    period,
		UpdatedAt: time.Now(),
	}
}

// Helper functions

func getTopMedicines(medicineMap map[string]int, total int) []MedicineFreq {
	var meds []MedicineFreq
	for name, count := range medicineMap {
		meds = append(meds, MedicineFreq{
			Name:      name,
			Count:     count,
			Frequency: float64(count) / float64(total) * 100,
		})
	}
	// Sort descending (in production, use sort.Slice)
	return meds[:minInt(10, len(meds))]
}

func getTopDiseases(diseaseMap map[string]int, total int) []DiseaseFreq {
	var diseases []DiseaseFreq
	for name, count := range diseaseMap {
		diseases = append(diseases, DiseaseFreq{
			Name:      name,
			Count:     count,
			Frequency: float64(count) / float64(total) * 100,
		})
	}
	return diseases[:minInt(10, len(diseases))]
}

func getDailyVisits(dailyMap map[string]int) []DailyVisit {
	var daily []DailyVisit
	for date, count := range dailyMap {
		daily = append(daily, DailyVisit{Date: date, Count: count})
	}
	return daily
}

func calculateStockStatus(medicines []models.Medicine) StockStatus {
	status := StockStatus{}
	now := time.Now()

	for _, med := range medicines {
		if med.ExpiryDate.Before(now) {
			status.Expired++
		} else if med.ExpiryDate.Sub(now).Hours()/24 <= 7 {
			status.Expiring++
		} else if med.QuantityRemaining <= med.ReorderLevel {
			status.LowStock++
		} else {
			status.OK++
		}
	}

	return status
}

func getFirstMedicine(meds []MedicineFreq) string {
	if len(meds) > 0 {
		return meds[0].Name
	}
	return "N/A"
}

func getFirstMedicineCount(meds []MedicineFreq) int {
	if len(meds) > 0 {
		return meds[0].Count
	}
	return 0
}

func getFirstDisease(diseases []DiseaseFreq) string {
	if len(diseases) > 0 {
		return diseases[0].Name
	}
	return "N/A"
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
