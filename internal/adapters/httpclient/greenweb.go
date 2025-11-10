package httpclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type GreenWebResponse struct {
	IsGreen bool `json:"isGreen"`
}

// MockData representa os dados mock para desenvolvimento offline
type MockData struct {
	Carbon   float64 `json:"carbon"`
	Green    bool    `json:"green"`
	Bytes    int     `json:"bytes"`
	Examples struct {
		GreenSites []struct {
			Domain string  `json:"domain"`
			Carbon float64 `json:"carbon"`
			Green  bool    `json:"green"`
			Bytes  int     `json:"bytes"`
		} `json:"green_sites"`
		NonGreenSites []struct {
			Domain string  `json:"domain"`
			Carbon float64 `json:"carbon"`
			Green  bool    `json:"green"`
			Bytes  int     `json:"bytes"`
		} `json:"non_green_sites"`
	} `json:"examples"`
}

type GreenWebAdapter struct {
	client   *http.Client
	baseURL  string
	mockPath string
	useMock  bool
}

func NewGreenWebAdapter(timeout time.Duration) *GreenWebAdapter {
	// Detectar se deve usar mock baseado nas variáveis de ambiente
	useMock := os.Getenv("ANALYSIS_ENV") == "dev" || os.Getenv("ENV") == "dev" || os.Getenv("GREENWEB_MOCK") == "true"

	mockPath := filepath.Join("pkg", "web", "mock.json")

	return &GreenWebAdapter{
		client: &http.Client{
			Timeout: timeout,
		},
		baseURL:  "https://api.greenweb.org",
		mockPath: mockPath,
		useMock:  useMock,
	}
}

// NewGreenWebAdapterWithMock cria um adapter configurado para usar mock
func NewGreenWebAdapterWithMock(timeout time.Duration, mockPath string) *GreenWebAdapter {
	return &GreenWebAdapter{
		client: &http.Client{
			Timeout: timeout,
		},
		baseURL:  "https://api.greenweb.org",
		mockPath: mockPath,
		useMock:  true,
	}
}

// SetBaseURL permite sobrescrever a baseURL do adapter (útil em exemplos/testes)
func (g *GreenWebAdapter) SetBaseURL(u string) {
	g.baseURL = u
}

// CheckGreenHost verifica se o domínio usa energia renovável
// Retorna true se o host é green, false caso contrário
func (g *GreenWebAdapter) CheckGreenHost(targetURL string) (bool, error) {
	// Se configurado para usar mock, usar dados locais
	if g.useMock {
		return g.checkGreenHostWithMock(targetURL)
	}

	// Tentar chamar a API real primeiro
	isGreen, err := g.checkGreenHostWithAPI(targetURL)
	if err != nil {
		// Se a API falhar, tentar usar mock como fallback
		fmt.Printf("API Green Web falhou, usando mock como fallback: %v\n", err)
		return g.checkGreenHostWithMock(targetURL)
	}

	return isGreen, nil
}

// checkGreenHostWithAPI chama a API real da Green Web Foundation
func (g *GreenWebAdapter) checkGreenHostWithAPI(targetURL string) (bool, error) {
	// Validação de entrada
	if targetURL == "" {
		return false, fmt.Errorf("URL não pode estar vazia")
	}

	// Validação da URL
	if _, err := url.ParseRequestURI(targetURL); err != nil {
		return false, fmt.Errorf("URL inválida: %w", err)
	}

	// Construir URL da API com encoding adequado
	apiURL := fmt.Sprintf("%s/check?url=%s", g.baseURL, url.QueryEscape(targetURL))

	resp, err := g.client.Get(apiURL)
	if err != nil {
		return false, fmt.Errorf("falha ao chamar Green Web API para URL '%s': %w", targetURL, err)
	}
	defer resp.Body.Close()

	// Verificar status da resposta
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("API Green Web retornou status %d para URL '%s'", resp.StatusCode, targetURL)
	}

	// Decodificar resposta JSON
	var result GreenWebResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, fmt.Errorf("erro ao decodificar resposta JSON da API Green Web para URL '%s': %w", targetURL, err)
	}

	return result.IsGreen, nil
}

// checkGreenHostWithMock usa dados mock para desenvolvimento offline
func (g *GreenWebAdapter) checkGreenHostWithMock(targetURL string) (bool, error) {
	// Carregar dados mock do arquivo JSON
	mockData, err := g.loadMockData()
	if err != nil {
		return false, fmt.Errorf("erro ao carregar dados mock: %w", err)
	}

	// Extrair domínio da URL
	domain := g.extractDomain(targetURL)

	// Verificar se o domínio está na lista de sites green
	for _, site := range mockData.Examples.GreenSites {
		if strings.Contains(domain, site.Domain) {
			return true, nil
		}
	}

	// Verificar se o domínio está na lista de sites não-green
	for _, site := range mockData.Examples.NonGreenSites {
		if strings.Contains(domain, site.Domain) {
			return false, nil
		}
	}

	// Se não encontrou o domínio específico, usar valor padrão do mock
	return mockData.Green, nil
}

// loadMockData carrega os dados mock do arquivo JSON
func (g *GreenWebAdapter) loadMockData() (*MockData, error) {
	// Tentar carregar do caminho relativo primeiro
	data, err := os.ReadFile(g.mockPath)
	if err != nil {
		// Se falhar, tentar resolver subindo diretórios (útil quando o teste é executado
		// a partir do diretório do pacote: internal/adapters/httpclient).
		// Tentamos até 6 níveis acima.
		var tryErr error
		cur := "."
		for i := 0; i < 6; i++ {
			tryPath := filepath.Join(cur, g.mockPath)
			data, tryErr = os.ReadFile(tryPath)
			if tryErr == nil {
				err = nil
				break
			}
			cur = filepath.Join(cur, "..")
		}

		// Ainda não achou, tentar caminho absoluto final
		if err != nil {
			absPath, _ := filepath.Abs(g.mockPath)
			data, err = os.ReadFile(absPath)
			if err != nil {
				return nil, fmt.Errorf("não foi possível carregar arquivo mock '%s': %w", g.mockPath, err)
			}
		}
	}

	var mockData MockData
	if err := json.Unmarshal(data, &mockData); err != nil {
		return nil, fmt.Errorf("erro ao decodificar arquivo mock: %w", err)
	}

	return &mockData, nil
}

// extractDomain extrai o domínio de uma URL
func (g *GreenWebAdapter) extractDomain(targetURL string) string {
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return targetURL
	}

	host := parsedURL.Host
	// Remover porta se existir
	if idx := strings.Index(host, ":"); idx != -1 {
		host = host[:idx]
	}

	return host
}
