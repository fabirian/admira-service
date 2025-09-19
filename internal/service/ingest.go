package service

import (
	"context"
	"time"

	"admira-service/internal/etl"
	"admira-service/pkg/ads"
	"admira-service/pkg/crm"
	"admira-service/pkg/metrics"
)

type IngestService struct {
	adsService       *ads.Service
	crmService       *crm.Service
	metricsCalculator *metrics.Calculator
	metricsRepo      *metrics.Repository
	etlProcessor     *etl.Processor
	transformer      *etl.Transformer
}

func NewIngestService(adsService *ads.Service, crmService *crm.Service,
	metricsCalculator *metrics.Calculator, metricsRepo *metrics.Repository,
	etlProcessor *etl.Processor, transformer *etl.Transformer) *IngestService {
	return &IngestService{
		adsService:       adsService,
		crmService:       crmService,
		metricsCalculator: metricsCalculator,
		metricsRepo:      metricsRepo,
		etlProcessor:     etlProcessor,
		transformer:      transformer,
	}
}

func (s *IngestService) RunIngestion(ctx context.Context, since time.Time) (int, error) {
	// Obtener datos de Ads
	adsData, err := s.adsService.GetPerformance(ctx, since)
	if err != nil {
		return 0, err
	}

	// Obtener datos de CRM
	crmData, err := s.crmService.GetOpportunities(ctx, since)
	if err != nil {
		return 0, err
	}

	// Normalizar datos
	normalizedAds := s.transformer.NormalizeAdsData(adsData)
	normalizedCRM := s.transformer.NormalizeCRMData(crmData)

	// Calcular métricas
	calculatedMetrics := s.metricsCalculator.CalculateMetrics(normalizedAds, normalizedCRM)

	// Guardar métricas
	for _, metric := range calculatedMetrics {
		s.metricsRepo.Save(metric)
	}

	return len(calculatedMetrics), nil
}

func (s *IngestService) RunIngestionWithRetry(ctx context.Context, since time.Time, maxRetries int) (int, error) {
	var lastErr error
	
	for attempt := 1; attempt <= maxRetries; attempt++ {
		processed, err := s.RunIngestion(ctx, since)
		if err == nil {
			return processed, nil
		}
		
		lastErr = err
		
		// Esperar antes del próximo intento (backoff exponencial)
		backoff := time.Duration(attempt*attempt) * time.Second
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		case <-time.After(backoff):
			// Continuar al siguiente intento
		}
	}
	
	return 0, lastErr
}