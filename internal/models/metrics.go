package models

type BusinessMetrics struct {
	Date          string  `json:"date"`
	Channel       string  `json:"channel"`
	CampaignID    string  `json:"campaign_id"`
	Clicks        int     `json:"clicks"`
	Impressions   int     `json:"impressions"`
	Cost          float64 `json:"cost"`
	Leads         int     `json:"leads"`
	Opportunities int     `json:"opportunities"`
	ClosedWon     int     `json:"closed_won"`
	Revenue       float64 `json:"revenue"`
	CPC           float64 `json:"cpc"`
	CPA           float64 `json:"cpa"`
	CVRLeadToOpp  float64 `json:"cvr_lead_to_opp"`
	CVROppToWon   float64 `json:"cvr_opp_to_won"`
	ROAS          float64 `json:"roas"`
}

type MetricsFilter struct {
	From        string
	To          string
	Channel     string
	UTMCampaign string
	Limit       int
	Offset      int
}