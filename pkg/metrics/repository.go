package metrics

import (
	"sync"
	"time"

	"admira-service/internal/models"
)

type Repository struct {
	metrics []models.BusinessMetrics
	mu      sync.RWMutex
}

func NewRepository() *Repository {
	return &Repository{
		metrics: make([]models.BusinessMetrics, 0),
	}
}

func (r *Repository) Save(metric models.BusinessMetrics) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.metrics = append(r.metrics, metric)
}

func (r *Repository) FindByFilter(filter models.MetricsFilter) []models.BusinessMetrics {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []models.BusinessMetrics
	var fromTime, toTime time.Time
	var err error

	if filter.From != "" {
		fromTime, err = time.Parse("2006-01-02", filter.From)
		if err != nil {
			return result
		}
	}

	if filter.To != "" {
		toTime, err = time.Parse("2006-01-02", filter.To)
		if err != nil {
			return result
		}
		// Añadir un día para incluir la fecha completa
		toTime = toTime.Add(24 * time.Hour)
	}

	count := 0
	for _, metric := range r.metrics {
		metricTime, err := time.Parse("2006-01-02", metric.Date)
		if err != nil {
			continue
		}

		// Aplicar filtros
		if filter.From != "" && metricTime.Before(fromTime) {
			continue
		}
		if filter.To != "" && metricTime.After(toTime) {
			continue
		}
		if filter.Channel != "" && metric.Channel != filter.Channel {
			continue
		}
		if filter.UTMCampaign != "" && metric.CampaignID != filter.UTMCampaign {
			continue
		}

		// Paginación
		if filter.Offset > 0 && count < filter.Offset {
			count++
			continue
		}

		result = append(result, metric)

		// Limitar resultados
		if filter.Limit > 0 && len(result) >= filter.Limit {
			break
		}
	}

	return result
}

func (r *Repository) GetAll() []models.BusinessMetrics {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.metrics
}