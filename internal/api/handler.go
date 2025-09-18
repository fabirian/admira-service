package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"admira-service/internal/models"
	"admira-service/pkg/ads"
	"admira-service/pkg/crm"
	"admira-service/pkg/metrics"
)

type Handler struct {
	adsService      *ads.Service
	crmService      *crm.Service
	metricsCalculator *metrics.Calculator
	metricsRepo     *metrics.Repository
}

func NewHandler(adsService *ads.Service, crmService *crm.Service, 
	metricsCalculator *metrics.Calculator, metricsRepo *metrics.Repository) *Handler {
	return &Handler{
		adsService:      adsService,
		crmService:      crmService,
		metricsCalculator: metricsCalculator,
		metricsRepo:     metricsRepo,
	}
}

func (h *Handler) RunIngest(c *gin.Context) {
	since := c.Query("since")
	
	var sinceTime time.Time
	if since != "" {
		parsedTime, err := time.Parse("2006-01-02", since)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
			return
		}
		sinceTime = parsedTime
	}
	
	// Obtener datos de Ads y CRM
	adsData, err := h.adsService.GetPerformance(c.Request.Context(), sinceTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ads data"})
		return
	}
	
	crmData, err := h.crmService.GetOpportunities(c.Request.Context(), sinceTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch crm data"})
		return
	}
	
	// Calcular métricas
	calculatedMetrics := h.metricsCalculator.CalculateMetrics(adsData, crmData)
	
	// Guardar métricas
	for _, metric := range calculatedMetrics {
		h.metricsRepo.Save(metric)
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Ingestion completed",
		"metrics_processed": len(calculatedMetrics),
	})
}

func (h *Handler) GetChannelMetrics(c *gin.Context) {
	filter := models.MetricsFilter{
		From:    c.Query("from"),
		To:      c.Query("to"),
		Channel: c.Query("channel"),
	}
	
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))
	filter.Limit = limit
	filter.Offset = offset
	
	metrics := h.metricsRepo.FindByFilter(filter)
	c.JSON(http.StatusOK, metrics)
}

func (h *Handler) GetFunnelMetrics(c *gin.Context) {
	filter := models.MetricsFilter{
		From:        c.Query("from"),
		To:          c.Query("to"),
		UTMCampaign: c.Query("utm_campaign"),
	}
	
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))
	filter.Limit = limit
	filter.Offset = offset
	
	metrics := h.metricsRepo.FindByFilter(filter)
	c.JSON(http.StatusOK, metrics)
}

func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

func (h *Handler) ReadyCheck(c *gin.Context) {
	// Verificar conexiones a dependencias
	if err := h.adsService.Ping(c.Request.Context()); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy"})
		return
	}
	
	if err := h.crmService.Ping(c.Request.Context()); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}