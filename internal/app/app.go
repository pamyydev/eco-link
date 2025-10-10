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
)

type carbonClient struct {
	client  *http.Client
	baseURL string
}

func (c *carbonClient) GetMetrics(url string) (domain.Metrics, error) {
	resp, err := c.client.Get(fmt.Sprintf("%s/analyze?url=%s", c.baseURL, url))
	if err != nil {
		return domain.Metrics{}, err
	}
	defer resp.Body.Close()

	var result struct {
		CarbonPerVisit float64 `json:"carbonPerVisit"`
		BytesPerVisit  int     `json:"bytesPerVisit"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return domain.Metrics{}, err
	}

	return domain.Metrics{
		CarbonPerVisit: result.CarbonPerVisit,
		BytesPerVisit:  result.BytesPerVisit,
	}, nil
}

type greenClient struct {
	client  *http.Client
	baseURL string
}

func (c *greenClient) CheckGreenHost(url string) (bool, error) {
	resp, err := c.client.Get(fmt.Sprintf("%s/check?url=%s", c.baseURL, url))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var result struct {
		IsGreen bool `json:"isGreen"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}

	return result.IsGreen, nil
}

func Run() {
	logger := log.New(os.Stdout, "[eco-link] ", log.LstdFlags)
	logger.Println("Iniciando aplicação...")

	// Get configs from env or use defaults
	timeout := getEnvDuration("ANALYSIS_TIMEOUT", 10*time.Second)
	targetURL := getEnvString("TARGET_URL", "https://example.com")

	// Initialize clients
	carbonClient := httpclient.NewWebsiteCarbonAdapter(timeout)
	greenClient := httpclient.NewGreenWebAdapter(timeout)

	// Initialize analysis service
	service := adapters.NewService(carbonClient, greenClient)

	// Analyze website
	logger.Printf("Analisando website: %s\n", targetURL)
	report, err := service.Analyze(targetURL)
	if err != nil {
		logger.Fatalf("Erro ao analisar website: %v\n", err)
	}

	// Print results
	logger.Printf("\nResultados da análise para %s:\n", report.Site.URL)
	logger.Printf("Carbono por visita: %.2fg CO2\n", report.Metrics.CarbonPerVisit)
	logger.Printf("Bytes por visita: %d bytes\n", report.Metrics.BytesPerVisit)
	logger.Printf("Hosting verde: %v\n", report.Metrics.GreenHost)
	logger.Printf("Eco-friendly: %v\n", report.Metrics.IsEcoFriendly())
}

// getEnvDuration retorna uma duração da variável de ambiente ou valor padrão
func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// getEnvString retorna uma string da variável de ambiente ou valor padrão
func getEnvString(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
