package api

import (
	"github.com/gin-gonic/gin"

	"admira-service/pkg/ads"
	"admira-service/pkg/crm"
	"admira-service/pkg/metrics"
)

func SetupRouter(adsService *ads.Service, crmService *crm.Service, 
	metricsCalculator *metrics.Calculator, metricsRepo *metrics.Repository, cfg interface{}) *gin.Engine {
	
	router := gin.New()
	
	// Middlewares
	router.Use(LoggerMiddleware())
	router.Use(RecoveryMiddleware())
	router.Use(CORSMiddleware())
	
	// Handler
	handler := NewHandler(adsService, crmService, metricsCalculator, metricsRepo)
	
	// Routes
	api := router.Group("/api/v1")
	{
		api.POST("/ingest/run", handler.RunIngest)
		api.GET("/metrics/channel", handler.GetChannelMetrics)
		api.GET("/metrics/funnel", handler.GetFunnelMetrics)
		api.POST("/export/run", handler.RunExport) // Opcional
	}
	
	// Health checks
	router.GET("/healthz", handler.HealthCheck)
	router.GET("/readyz", handler.ReadyCheck)
	
	return router
}