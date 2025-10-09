package ports

import "github.com/pamelamiranda/eco-link/internal/domain"

// WebsiteCarbonClient defines contract for carbon metrics API
type WebsiteCarbonClient interface {
	GetMetrics(url string) (domain.Metrics, error)
}

// GreenWebClient defines contract for green hosting verification
type GreenWebClient interface {
	CheckGreenHost(url string) (bool, error)
}