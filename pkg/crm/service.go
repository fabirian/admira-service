package crm

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

func (s *Service) GetOpportunities(ctx context.Context, since time.Time) ([]models.CRMOpportunity, error) {
	response, err := s.client.GetOpportunities(ctx)
	if err != nil {
		return nil, err
	}

	// Filtrar por fecha since si se proporciona
	var filtered []models.CRMOpportunity
	for _, opp := range response.External.CRM.Opportunities {
		if since.IsZero() || !opp.CreatedAt.Before(since) {
			opp.IngestedAt = time.Now()
			filtered = append(filtered, opp)
		}
	}

	return filtered, nil
}

func (s *Service) Ping(ctx context.Context) error {
	_, err := s.client.GetOpportunities(ctx)
	return err
}