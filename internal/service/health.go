package service

import (
	"context"

	"admira-service/pkg/ads"
	"admira-service/pkg/crm"
)

type HealthService struct {
	adsService *ads.Service
	crmService *crm.Service
}

func NewHealthService(adsService *ads.Service, crmService *crm.Service) *HealthService {
	return &HealthService{
		adsService: adsService,
		crmService: crmService,
	}
}

func (s *HealthService) CheckHealth(ctx context.Context) map[string]string {
	status := make(map[string]string)
	
	// Verificar Ads service
	if err := s.adsService.Ping(ctx); err != nil {
		status["ads_service"] = "unhealthy"
	} else {
		status["ads_service"] = "healthy"
	}
	
	// Verificar CRM service
	if err := s.crmService.Ping(ctx); err != nil {
		status["crm_service"] = "unhealthy"
	} else {
		status["crm_service"] = "healthy"
	}
	
	status["overall"] = "healthy"
	if status["ads_service"] == "unhealthy" || status["crm_service"] == "unhealthy" {
		status["overall"] = "degraded"
	}
	
	return status
}