package adapters

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pamelamiranda/eco-link/internal/domain"
)

// MetricsResponse represents the API response format
type MetricsResponse struct {
	Site           string  `json:"site"`
	CarbonPerVisit float64 `json:"carbonPerVisit"`
	GreenHost      bool    `json:"greenHost"`
}

// Service handles website sustainability analysis
type Service struct {
	carbonClient WebsiteCarbonClient
	greenClient  GreenWebClient
}

// NewService creates a new analysis service
func NewService(carbonClient WebsiteCarbonClient, greenClient GreenWebClient) *Service {
	return &Service{
		carbonClient: carbonClient,
		greenClient:  greenClient,
	}
}

// fetchMetrics gets carbon metrics from the API
func (s *Service) fetchMetrics(url string) (*MetricsResponse, error) {
	resp, err := s.client.Get(fmt.Sprintf("%s/analyze?url=%s", s.apiURL, url))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch metrics: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var metrics MetricsResponse
	if err := json.NewDecoder(resp.Body).Decode(&metrics); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &metrics, nil
}

// calculateBytes measures the website size in bytes
func (s *Service) calculateBytes(url string) (int, error) {
	resp, err := s.client.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch site: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response: %w", err)
	}

	return len(body), nil
}

// checkGreenHost verifies if the host uses renewable energy
func (s *Service) checkGreenHost(url string) (bool, error) {
	resp, err := s.client.Get(fmt.Sprintf("%s/check?url=%s", s.greenHost, url))
	if err != nil {
		return false, fmt.Errorf("failed to check green host: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		IsGreen bool `json:"isGreen"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, fmt.Errorf("failed to decode green host response: %w", err)
	}

	return result.IsGreen, nil
}

// Analyze fetches and processes sustainability metrics for a given URL
func (s *Service) Analyze(url string) (domain.Report, error) {
	site, err := domain.NewSite(url)
	if err != nil {
		return domain.Report{}, fmt.Errorf("invalid site URL: %w", err)
	}

	metrics, err := s.carbonClient.GetMetrics(url)
	if err != nil {
		return domain.Report{}, fmt.Errorf("failed to get carbon metrics: %w", err)
	}

	isGreen, err := s.greenClient.CheckGreenHost(url)
	if err != nil {
		return domain.Report{}, fmt.Errorf("failed to check green host: %w", err)
	}

	metrics.GreenHost = isGreen
	report := domain.NewReport(*site, metrics)

	return report, nil
}