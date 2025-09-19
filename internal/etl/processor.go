package etl

import (
	"context"
	"time"

	"admira-service/internal/models"
	"admira-service/pkg/ads"
	"admira-service/pkg/crm"
	"admira-service/pkg/metrics"
)

type Processor struct {
	adsService       *ads.Service
	crmService       *crm.Service
	metricsCalculator *metrics.Calculator
	metricsRepo      *metrics.Repository
}

func NewProcessor(adsService *ads.Service, crmService *crm.Service,
	metricsCalculator *metrics.Calculator, metricsRepo *metrics.Repository) *Processor {
	return &Processor{
		adsService:       adsService,
		crmService:       crmService,
		metricsCalculator: metricsCalculator,
		metricsRepo:      metricsRepo,
	}
}

func (p *Processor) ProcessData(ctx context.Context, since time.Time) (int, error) {
	// Obtener datos de Ads con reintentos
	adsData, err := p.retryGetAds(ctx, since)
	if err != nil {
		return 0, err
	}

	// Obtener datos de CRM con reintentos
	crmData, err := p.retryGetCRM(ctx, since)
	if err != nil {
		return 0, err
	}

	// Calcular métricas
	calculatedMetrics := p.metricsCalculator.CalculateMetrics(adsData, crmData)

	// Guardar métricas
	for _, metric := range calculatedMetrics {
		p.metricsRepo.Save(metric)
	}

	return len(calculatedMetrics), nil
}

func (p *Processor) retryGetAds(ctx context.Context, since time.Time) ([]models.AdsPerformance, error) {
	var err error
	var result []models.AdsPerformance
	
	for attempt := 1; attempt <= 3; attempt++ {
		result, err = p.adsService.GetPerformance(ctx, since)
		if err == nil {
			return result, nil
		}
		
		// Esperar antes del próximo intento (backoff exponencial)
		backoff := time.Duration(attempt*attempt) * time.Second
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(backoff):
			// Continuar al siguiente intento
		}
	}
	
	return nil, err
}

func (p *Processor) retryGetCRM(ctx context.Context, since time.Time) ([]models.CRMOpportunity, error) {
	var err error
	var result []models.CRMOpportunity
	
	for attempt := 1; attempt <= 3; attempt++ {
		result, err = p.crmService.GetOpportunities(ctx, since)
		if err == nil {
			return result, nil
		}
		
		// Esperar antes del próximo intento (backoff exponencial)
		backoff := time.Duration(attempt*attempt) * time.Second
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(backoff):
			// Continuar al siguiente intento
		}
	}
	
	return nil, err
}