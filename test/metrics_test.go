package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"admira-service/internal/models"
	"admira-service/pkg/metrics"
)

// Función auxiliar para parsear fecha
func parseTime(dateStr string) time.Time {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		// Si falla, devolver tiempo actual
		return time.Now()
	}
	return t
}

func TestMetricsCalculation(t *testing.T) {
	repo := metrics.NewRepository()
	calculator := metrics.NewCalculator(repo)

	testDate := "2025-08-01"
	
	adsData := []models.AdsPerformance{
		{
			Date:        testDate,
			CampaignID:  "C-1001",
			Channel:     "google_ads",
			Clicks:      100,
			Impressions: 5000,
			Cost:        50.0,
			UTMCampaign: "test_campaign",
			UTMSource:   "google",
			UTMMedium:   "cpc",
			IngestedAt:  time.Now(),
		},
	}

	crmData := []models.CRMOpportunity{
		{
			OpportunityID: "O-9001",
			ContactEmail:  "test@example.com",
			Stage:         "closed_won",
			Amount:        1000.0,
			CreatedAt:     parseTime(testDate),
			UTMCampaign:   "test_campaign",
			UTMSource:     "google",
			UTMMedium:     "cpc",
			IngestedAt:    time.Now(),
		},
	}

	result := calculator.CalculateMetrics(adsData, crmData)

	assert.Len(t, result, 1, "Should have exactly 1 metric result")
	assert.Equal(t, 0.5, result[0].CPC)           // 50 / 100 = 0.5
	assert.Equal(t, 0.5, result[0].CPA)           // 50 / 100 = 0.5
	assert.Equal(t, 0.01, result[0].CVRLeadToOpp) // 1 / 100 = 0.01 ✅ CORRECTO
	assert.Equal(t, 1.0, result[0].CVROppToWon)   // 1 / 1 = 1.0
	assert.Equal(t, 20.0, result[0].ROAS)         // 1000 / 50 = 20.0
}

func TestRepositoryOperations(t *testing.T) {
	repo := metrics.NewRepository()

	metric := models.BusinessMetrics{
		Date:          "2025-08-01",
		Channel:       "google_ads",
		CampaignID:    "C-1001",
		Clicks:        100,
		Impressions:   5000,
		Cost:          50.0,
		Leads:         10,
		Opportunities: 5,
		ClosedWon:     2,
		Revenue:       1000.0,
		CPC:           0.5,
		CPA:           5.0,
		CVRLeadToOpp:  0.5,
		CVROppToWon:   0.4,
		ROAS:          20.0,
	}

	// Test Save
	repo.Save(metric)

	// Test FindByFilter
	filter := models.MetricsFilter{
		From: "2025-08-01",
		To:   "2025-08-01",
	}
	result := repo.FindByFilter(filter)

	assert.Len(t, result, 1)
	assert.Equal(t, "google_ads", result[0].Channel)
	assert.Equal(t, "C-1001", result[0].CampaignID)
}

func TestDivisionByZero(t *testing.T) {
	repo := metrics.NewRepository()
	calculator := metrics.NewCalculator(repo)

	// Test CPC con clicks = 0
	cpc := calculator.CalculateCPC(100.0, 0)
	assert.Equal(t, 0.0, cpc)

	// Test CPA con leads = 0
	cpa := calculator.CalculateCPA(100.0, 0)
	assert.Equal(t, 0.0, cpa)

	// Test CVR con denominator = 0
	cvr := calculator.CalculateCVR(10, 0)
	assert.Equal(t, 0.0, cvr)

	// Test ROAS con cost = 0
	roas := calculator.CalculateROAS(1000.0, 0)
	assert.Equal(t, 0.0, roas)
}