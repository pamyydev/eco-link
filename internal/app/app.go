package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pamelamiranda/eco-link/internal/adapters"
	"github.com/pamelamiranda/eco-link/internal/adapters/httpclient"
	"github.com/pamelamiranda/eco-link/internal/domain"
	"github.com/pamelamiranda/eco-link/pkg/config"
)


type Config struct {
	Timeout   time.Duration
	TargetURL string
}

func NewConfig() Config {
	return Config{
		Timeout:   config.GetEnvDuration("ANALYSIS_TIMEOUT", 10*time.Second),
		TargetURL: config.GetEnvString("TARGET_URL", "https://example.com"),
	}
}

type Application struct {
	logger  *log.Logger
	cfg     Config
	service *adapters.Service
}

func NewApplication(logger *log.Logger, cfg Config, service *adapters.Service) *Application {
	return &Application{
		logger:  logger,
		cfg:     cfg,
		service: service,
	}
}

func (a *Application) Run() {
	a.logger.Println("Iniciando aplicação...")

	// Analyze website
	a.logger.Printf("Analisando website: %s\n", a.cfg.TargetURL)
	report, err := a.service.Analyze(a.cfg.TargetURL)
	if err != nil {
		a.logger.Fatalf("Erro ao analisar website: %v\n", err)
	}

	// Print results
	a.logger.Printf("\nResultados da análise para %s:\n", report.Site.URL)
	a.logger.Printf("Carbono por visita: %.2fg CO2\n", report.Metrics.CarbonPerVisit)
	a.logger.Printf("Bytes por visita: %d bytes\n", report.Metrics.BytesPerVisit)
	a.logger.Printf("Hosting verde: %v\n", report.Metrics.GreenHost)
	a.logger.Printf("Eco-friendly: %v\n", report.Metrics.IsEcoFriendly())
}

