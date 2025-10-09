package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pamelamiranda/eco-link/internal/adapters"
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
	fmt.Println("eco-link app iniciado!")

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Initialize clients
	carbon := &carbonClient{
		client:  httpClient,
		baseURL: "https://api.websitecarbon.com",
	}
	green := &greenClient{
		client:  httpClient,
		baseURL: "https://api.green-host.com",
	}

	// Initialize analysis service
	service := adapters.NewService(carbon, green)

	// Example: Analyze a website
	report, err := service.Analyze("https://example.com")
	if err != nil {
		fmt.Printf("Error analyzing website: %v\n", err)
		return
	}

	// Print analysis results
	fmt.Printf("\nAnalysis Results for %s:\n", report.Site.URL)
	fmt.Printf("Carbon per visit: %.2fg CO2\n", report.Metrics.CarbonPerVisit)
	fmt.Printf("Bytes per visit: %d bytes\n", report.Metrics.BytesPerVisit)
	fmt.Printf("Green hosting: %v\n", report.Metrics.GreenHost)
	fmt.Printf("Eco-friendly: %v\n", report.Metrics.IsEcoFriendly())
}
