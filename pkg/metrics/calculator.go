package metrics

import (
	"time"

	"admira-service/internal/models"
)

type Calculator struct {
	repo *Repository
}

func NewCalculator(repo *Repository) *Calculator {
	return &Calculator{repo: repo}
}

func (c *Calculator) CalculateMetrics(ads []models.AdsPerformance, opps []models.CRMOpportunity) []models.BusinessMetrics {
	// Agrupar por UTM parameters ignorando la fecha inicialmente
	groupedData := c.groupDataByUTM(ads, opps)
	
	var metrics []models.BusinessMetrics
	
	for key, data := range groupedData {
		// Usar la fecha de los datos de Ads si está disponible, sino de CRM
		date := key.Date
		if date == "" {
			// Si no hay fecha de Ads, usar la fecha más común de CRM o today
			date = time.Now().Format("2006-01-02")
		}
		
		metric := models.BusinessMetrics{
			Date:          date,
			Channel:       key.Channel,
			CampaignID:    key.CampaignID,
			Clicks:        data.Clicks,
			Impressions:   data.Impressions,
			Cost:          data.Cost,
			Leads:         data.Leads,
			Opportunities: data.Opportunities,
			ClosedWon:     data.ClosedWon,
			Revenue:       data.Revenue,
		}
		
		// Calcular métricas
		metric.CPC = c.CalculateCPC(metric.Cost, metric.Clicks)
		metric.CPA = c.CalculateCPA(metric.Cost, metric.Leads)
		metric.CVRLeadToOpp = c.CalculateCVR(metric.Opportunities, metric.Leads)
		metric.CVROppToWon = c.CalculateCVR(metric.ClosedWon, metric.Opportunities)
		metric.ROAS = c.CalculateROAS(metric.Revenue, metric.Cost)
		
		metrics = append(metrics, metric)
	}
	
	return metrics
}

func (c *Calculator) CalculateCPC(cost float64, clicks int) float64 {
	if clicks == 0 {
		return 0
	}
	return cost / float64(clicks)
}

func (c *Calculator) CalculateCPA(cost float64, leads int) float64 {
	if leads == 0 {
		return 0
	}
	return cost / float64(leads)
}

func (c *Calculator) CalculateCVR(numerator, denominator int) float64 {
	if denominator == 0 {
		return 0
	}
	return float64(numerator) / float64(denominator)
}

func (c *Calculator) CalculateROAS(revenue, cost float64) float64 {
	if cost == 0 {
		return 0
	}
	return revenue / cost
}

// Struct auxiliar para agrupamiento - SIN fecha para agrupar por UTM solamente
type utmKey struct {
	Channel     string
	CampaignID  string
	UTMCampaign string
	UTMSource   string
	UTMMedium   string
	Date        string // Mantenemos fecha separada para tracking
}

type groupedData struct {
	Clicks        int
	Impressions   int
	Cost          float64
	Leads         int
	Opportunities int
	ClosedWon     int
	Revenue       float64
}

func (c *Calculator) groupDataByUTM(ads []models.AdsPerformance, opps []models.CRMOpportunity) map[utmKey]groupedData {
	grouped := make(map[utmKey]groupedData)
	
	// Procesar datos de Ads
	for _, ad := range ads {
		key := utmKey{
			Channel:     ad.Channel,
			CampaignID:  ad.CampaignID,
			UTMCampaign: ad.UTMCampaign,
			UTMSource:   ad.UTMSource,
			UTMMedium:   ad.UTMMedium,
			Date:        ad.Date, // Guardamos la fecha pero no la usamos para agrupar
		}
		
		// Crear clave sin fecha para agrupamiento
		groupKey := utmKey{
			Channel:     key.Channel,
			CampaignID:  key.CampaignID,
			UTMCampaign: key.UTMCampaign,
			UTMSource:   key.UTMSource,
			UTMMedium:   key.UTMMedium,
		}
		
		data := grouped[groupKey]
		data.Clicks += ad.Clicks
		data.Impressions += ad.Impressions
		data.Cost += ad.Cost
		// Considerar cada click como un lead potencial
		data.Leads += ad.Clicks
		grouped[groupKey] = data
	}
	
	// Procesar datos de CRM
	for _, opp := range opps {
		key := utmKey{
			UTMCampaign: opp.UTMCampaign,
			UTMSource:   opp.UTMSource,
			UTMMedium:   opp.UTMMedium,
		}
		
		// Buscar clave existente que coincida con los UTM parameters
		var foundKey utmKey
		found := false
		
		for existingKey := range grouped {
			if existingKey.UTMCampaign == key.UTMCampaign &&
				existingKey.UTMSource == key.UTMSource &&
				existingKey.UTMMedium == key.UTMMedium {
				foundKey = existingKey
				found = true
				break
			}
		}
		
		if !found {
			// Si no encontramos campaña de Ads, crear una nueva entrada
			foundKey = utmKey{
				UTMCampaign: key.UTMCampaign,
				UTMSource:   key.UTMSource,
				UTMMedium:   key.UTMMedium,
			}
		}
		
		data := grouped[foundKey]
		data.Opportunities++
		if opp.Stage == "closed_won" {
			data.ClosedWon++
			data.Revenue += opp.Amount
		}
		grouped[foundKey] = data
	}
	
	return grouped
}