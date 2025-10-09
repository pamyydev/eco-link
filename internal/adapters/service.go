package adapters

import (
	"fmt"
	"github.com/pamelamiranda/eco-link/internal/domain"
	"github.com/pamelamiranda/eco-link/internal/ports"
)

// Service handles website sustainability analysis
type Service struct {
	carbonClient ports.WebsiteCarbonClient
	greenClient  ports.GreenWebClient
}

// NewService creates a new analysis service
func NewService(carbonClient ports.WebsiteCarbonClient, greenClient ports.GreenWebClient) *Service {
	return &Service{
		carbonClient: carbonClient,
		greenClient:  greenClient,
	}
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