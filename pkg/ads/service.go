package ads

import (
	"context"
	"time"

	"admira-service/internal/models"
)

type Service struct {
	client *Client
}

func NewService(client *Client) *Service {
	return &Service{client: client}
}

func (s *Service) GetPerformance(ctx context.Context, since time.Time) ([]models.AdsPerformance, error) {
	response, err := s.client.GetPerformance(ctx)
	if err != nil {
		return nil, err
	}

	// Filtrar por fecha since si se proporciona
	var filtered []models.AdsPerformance
	for _, perf := range response.External.Ads.Performance {
		perfDate, err := time.Parse("2006-01-02", perf.Date)
		if err != nil {
			continue // Saltar fechas inválidas
		}

		if since.IsZero() || !perfDate.Before(since) {
			perf.IngestedAt = time.Now()
			filtered = append(filtered, perf)
		}
	}

	return filtered, nil
}

func (s *Service) Ping(ctx context.Context) error {
	_, err := s.client.GetPerformance(ctx)
	return err
}