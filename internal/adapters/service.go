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
// Integra verificação de energia renovável via Green Web Foundation
func (s *Service) Analyze(url string) (domain.Report, error) {
	// Validar URL do site
	site, err := domain.NewSite(url)
	if err != nil {
		return domain.Report{}, fmt.Errorf("URL inválida: %w", err)
	}

	// Obter métricas de carbono
	metrics, err := s.carbonClient.GetMetrics(url)
	if err != nil {
		return domain.Report{}, fmt.Errorf("falha ao obter métricas de carbono: %w", err)
	}

	// Verificar se o host usa energia renovável via Green Web Foundation
	isGreen, err := s.greenClient.CheckGreenHost(url)
	if err != nil {
		return domain.Report{}, fmt.Errorf("falha ao verificar energia renovável: %w", err)
	}

	// Integrar verificação green ao relatório final
	metrics.GreenHost = isGreen
	report := domain.NewReport(*site, metrics)

	return report, nil
}
