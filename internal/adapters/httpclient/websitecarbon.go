package httpclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pamelamiranda/eco-link/internal/domain"
)

type WebsiteCarbonResponse struct {
	Bytes      int     `json:"bytes"`
	Green      bool    `json:"green"`
	CO2        float64 `json:"gco2e"`
	Rating     string  `json:"rating"`
	CleanerThan float64 `json:"cleanerThan"`
}

type WebsiteCarbonAdapter struct {
	client  *http.Client
	baseURL string
}

func NewWebsiteCarbonAdapter(timeout time.Duration) *WebsiteCarbonAdapter {
	return &WebsiteCarbonAdapter{
		client: &http.Client{
			Timeout: timeout,
		},
		baseURL: "https://api.websitecarbon.com",
	}
}

func (w *WebsiteCarbonAdapter) GetMetrics(url string) (domain.Metrics, error) {
	resp, err := w.client.Get(fmt.Sprintf("%s/data?bytes=1000000&green=1", w.baseURL))
	if err != nil {
		return domain.Metrics{}, fmt.Errorf("falha ao chamar Website Carbon API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return domain.Metrics{}, fmt.Errorf("API retornou status code inválido: %d", resp.StatusCode)
	}

	var result WebsiteCarbonResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return domain.Metrics{}, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	return domain.Metrics{
		CarbonPerVisit: result.CO2,
		BytesPerVisit:  result.Bytes,
		GreenHost:      result.Green,
	}, nil
}