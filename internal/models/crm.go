package models

import "time"

type CRMOpportunity struct {
	OpportunityID string    `json:"opportunity_id"`
	ContactEmail  string    `json:"contact_email"`
	Stage         string    `json:"stage"`
	Amount        float64   `json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
	UTMCampaign   string    `json:"utm_campaign"`
	UTMSource     string    `json:"utm_source"`
	UTMMedium     string    `json:"utm_medium"`
	IngestedAt    time.Time `json:"ingested_at"`
}

type CRMResponse struct {
	External struct {
		CRM struct {
			Opportunities []CRMOpportunity `json:"opportunities"`
		} `json:"crm"`
	} `json:"external"`
}