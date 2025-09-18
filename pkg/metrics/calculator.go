package metrics

import "admira-service/internal/models"

type Calculator struct {
	repo *Repository
}

func NewCalculator(repo *Repository) *Calculator {
	return &Calculator{repo: repo}
}

func (c *Calculator) CalculateMetrics(ads []models.AdsPerformance, opps []models.CRMOpportunity) []models.BusinessMetrics {
	// Agrupar por UTM parameters
	groupedData := c.groupDataByUTM(ads, opps)
	
	var metrics []models.BusinessMetrics
	
	for key, data := range groupedData {
		metric := models.BusinessMetrics{
			Date:        key.Date,
			Channel:     key.Channel,
			CampaignID:  key.CampaignID,
			Clicks:      data.Clicks,
			Impressions: data.Impressions,
			Cost:        data.Cost,
			Leads:       data.Leads,
			Opportunities: data.Opportunities,
			ClosedWon:   data.ClosedWon,
			Revenue:     data.Revenue,
		}
		
		// Calcular métricas
		metric.CPC = c.calculateCPC(metric.Cost, metric.Clicks)
		metric.CPA = c.calculateCPA(metric.Cost, metric.Leads)
		metric.CVRLeadToOpp = c.calculateCVR(metric.Opportunities, metric.Leads)
		metric.CVROppToWon = c.calculateCVR(metric.ClosedWon, metric.Opportunities)
		metric.ROAS = c.calculateROAS(metric.Revenue, metric.Cost)
		
		metrics = append(metrics, metric)
	}
	
	return metrics
}

func (c *Calculator) calculateCPC(cost float64, clicks int) float64 {
	if clicks == 0 {
		return 0
	}
	return cost / float64(clicks)
}

func (c *Calculator) calculateCPA(cost float64, leads int) float64 {
	if leads == 0 {
		return 0
	}
	return cost / float64(leads)
}

func (c *Calculator) calculateCVR(numerator, denominator int) float64 {
	if denominator == 0 {
		return 0
	}
	return float64(numerator) / float64(denominator)
}

func (c *Calculator) calculateROAS(revenue, cost float64) float64 {
	if cost == 0 {
		return 0
	}
	return revenue / cost
}

// Struct auxiliar para agrupamiento
type utmKey struct {
	Date       string
	Channel    string
	CampaignID string
	UTMCampaign string
	UTMSource  string
	UTMMedium  string
}

type groupedData struct {
	Clicks       int
	Impressions  int
	Cost         float64
	Leads        int
	Opportunities int
	ClosedWon    int
	Revenue      float64
}

func (c *Calculator) groupDataByUTM(ads []models.AdsPerformance, opps []models.CRMOpportunity) map[utmKey]groupedData {
	grouped := make(map[utmKey]groupedData)
	
	// Procesar datos de Ads
	for _, ad := range ads {
		key := utmKey{
			Date:        ad.Date,
			Channel:     ad.Channel,
			CampaignID:  ad.CampaignID,
			UTMCampaign: ad.UTMCampaign,
			UTMSource:   ad.UTMSource,
			UTMMedium:   ad.UTMMedium,
		}
		
		data := grouped[key]
		data.Clicks += ad.Clicks
		data.Impressions += ad.Impressions
		data.Cost += ad.Cost
		// Considerar cada click como un lead potencial
		data.Leads += ad.Clicks
		grouped[key] = data
	}
	
	// Procesar datos de CRM
	for _, opp := range opps {
		key := utmKey{
			Date:        opp.CreatedAt.Format("2006-01-02"),
			UTMCampaign: opp.UTMCampaign,
			UTMSource:   opp.UTMSource,
			UTMMedium:   opp.UTMMedium,
		}
		
		data := grouped[key]
		data.Opportunities++
		if opp.Stage == "closed_won" {
			data.ClosedWon++
			data.Revenue += opp.Amount
		}
		grouped[key] = data
	}
	
	return grouped
}