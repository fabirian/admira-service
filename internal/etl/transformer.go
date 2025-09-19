package etl

import (
	"strings"
	"time"

	"admira-service/internal/models"
)

type Transformer struct{}

func NewTransformer() *Transformer {
	return &Transformer{}
}

func (t *Transformer) NormalizeAdsData(adsData []models.AdsPerformance) []models.AdsPerformance {
	var normalized []models.AdsPerformance
	
	for _, ad := range adsData {
		// Normalizar campos nulos o vacíos
		if ad.UTMCampaign == "" {
			ad.UTMCampaign = "unknown_campaign"
		}
		if ad.UTMSource == "" {
			ad.UTMSource = "unknown_source"
		}
		if ad.UTMMedium == "" {
			ad.UTMMedium = "unknown_medium"
		}
		
		// Normalizar channel a minúsculas
		ad.Channel = strings.ToLower(ad.Channel)
		
		// Validar y normalizar fecha
		if ad.Date == "" {
			ad.Date = time.Now().Format("2006-01-02")
		}
		
		normalized = append(normalized, ad)
	}
	
	return normalized
}

func (t *Transformer) NormalizeCRMData(crmData []models.CRMOpportunity) []models.CRMOpportunity {
	var normalized []models.CRMOpportunity
	
	for _, opp := range crmData {
		// Normalizar campos nulos o vacíos
		if opp.UTMCampaign == "" {
			opp.UTMCampaign = "unknown_campaign"
		}
		if opp.UTMSource == "" {
			opp.UTMSource = "unknown_source"
		}
		if opp.UTMMedium == "" {
			opp.UTMMedium = "unknown_medium"
		}
		
		// Normalizar stage
		opp.Stage = strings.ToLower(strings.ReplaceAll(opp.Stage, " ", "_"))
		
		normalized = append(normalized, opp)
	}
	
	return normalized
}

func (t *Transformer) GenerateUTMKey(utmCampaign, utmSource, utmMedium string) string {
	return utmCampaign + "|" + utmSource + "|" + utmMedium
}

func (t *Transformer) ParseUTMKey(key string) (string, string, string) {
	parts := strings.Split(key, "|")
	if len(parts) != 3 {
		return "unknown", "unknown", "unknown"
	}
	return parts[0], parts[1], parts[2]
}