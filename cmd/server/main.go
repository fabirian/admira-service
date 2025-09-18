package main

import (
	"log"
	"os"

	"admira-service/internal/api"
	"admira-service/internal/config"
	"admira-service/pkg/ads"
	"admira-service/pkg/crm"
	"admira-service/pkg/metrics"
)

func main() {
	// Cargar configuración
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	// Inicializar clients
	adsClient := ads.NewClient(cfg.AdsAPIURL)
	crmClient := crm.NewClient(cfg.CrmAPIURL)

	// Inicializar servicios
	adsService := ads.NewService(adsClient)
	crmService := crm.NewService(crmClient)
	metricsRepo := metrics.NewRepository()
	metricsCalculator := metrics.NewCalculator(metricsRepo)

	// Configurar router
	router := api.SetupRouter(adsService, crmService, metricsCalculator, metricsRepo, cfg)

	// Iniciar servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}