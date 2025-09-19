package service

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"admira-service/internal/models"
	"admira-service/pkg/metrics"
)

type ExportService struct {
	metricsRepo *metrics.Repository
	sinkURL     string
	sinkSecret  string
}

func NewExportService(metricsRepo *metrics.Repository, sinkURL, sinkSecret string) *ExportService {
	return &ExportService{
		metricsRepo: metricsRepo,
		sinkURL:     sinkURL,
		sinkSecret:  sinkSecret,
	}
}

func (s *ExportService) ExportData(ctx context.Context, date string) error {
	if s.sinkURL == "" {
		return fmt.Errorf("sink URL not configured")
	}

	// Obtener métricas para la fecha específica
	filter := models.MetricsFilter{
		From: date,
		To:   date,
	}
	metrics := s.metricsRepo.FindByFilter(filter)

	// Exportar cada métrica
	for _, metric := range metrics {
		if err := s.exportMetric(ctx, metric); err != nil {
			return err
		}
	}

	return nil
}

func (s *ExportService) exportMetric(ctx context.Context, metric models.BusinessMetrics) error {
	// Convertir a JSON
	jsonData, err := json.Marshal(metric)
	if err != nil {
		return fmt.Errorf("marshaling metric: %w", err)
	}

	// Calcular HMAC signature
	signature := s.calculateSignature(jsonData)

	// Crear request
	req, err := http.NewRequestWithContext(ctx, "POST", s.sinkURL, bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Signature", signature)

	// Enviar request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (s *ExportService) calculateSignature(data []byte) string {
	h := hmac.New(sha256.New, []byte(s.sinkSecret))
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}