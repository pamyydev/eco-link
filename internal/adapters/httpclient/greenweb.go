package httpclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type GreenWebResponse struct {
	IsGreen bool `json:"isGreen"`
}

type GreenWebAdapter struct {
	client  *http.Client
	baseURL string
}

func NewGreenWebAdapter(timeout time.Duration) *GreenWebAdapter {
	return &GreenWebAdapter{
		client: &http.Client{
			Timeout: timeout,
		},
		baseURL: "https://api.greenweb.org",
	}
}

func (g *GreenWebAdapter) CheckGreenHost(url string) (bool, error) {
	resp, err := g.client.Get(fmt.Sprintf("%s/check?url=%s", g.baseURL, url))
	if err != nil {
		return false, fmt.Errorf("falha ao chamar Green Web API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("API retornou status code inválido: %d", resp.StatusCode)
	}

	var result GreenWebResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	return result.IsGreen, nil
}