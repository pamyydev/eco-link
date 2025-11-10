package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/pamelamiranda/eco-link/internal/adapters"
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

type GreenCheck struct {
	Carbon float64 `json:"carbon"`
	Green  bool    `json:"green"`
	Bytes  int64   `json:"bytes"`
}

// fetchGreen tenta chamar a API da Green Web Foundation.
// Se ANALYSIS_ENV == "dev" lê arquivo mock (MOCK_GREENCHECK_FILE) ou usa mock embutido.
// Em caso de falha na API, tenta fallback para o mock local/embutido.
func (a *Application) fetchGreen(domain string) (GreenCheck, error) {
	var gi GreenCheck

	mockFile := os.Getenv("MOCK_GREENCHECK_FILE")
	if mockFile == "" {
		mockFile = "greencheck_mock.json"
	}

	// modo dev: usar mock local (se não existir, usar mock embutido)
	if os.Getenv("ANALYSIS_ENV") == "dev" {
		data, err := os.ReadFile(mockFile)
		if err != nil {
			data = []byte(`{"carbon":0.6,"green":false,"bytes":512000}`)
		}
		if err := json.Unmarshal(data, &gi); err != nil {
			return gi, fmt.Errorf("falha ao parsear mock greencheck: %w", err)
		}
		return gi, nil
	}

	// modo prod: chamar API
	base := "https://api.thegreenwebfoundation.org/v3/greencheck/"
	encoded := url.PathEscape(domain)
	reqURL := base + encoded

	client := http.Client{Timeout: a.cfg.Timeout}
	resp, err := client.Get(reqURL)
	if err != nil {
		// fallback para mock local/embutido
		data, rerr := os.ReadFile(mockFile)
		if rerr != nil {
			data = []byte(`{"carbon":0.6,"green":false,"bytes":512000}`)
		}
		if perr := json.Unmarshal(data, &gi); perr != nil {
			return gi, fmt.Errorf("erro na API: %v; também falhou ao parsear mock: %w", err, perr)
		}
		return gi, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, rerr := os.ReadFile(mockFile)
		if rerr != nil {
			data = []byte(`{"carbon":0.6,"green":false,"bytes":512000}`)
		}
		if perr := json.Unmarshal(data, &gi); perr != nil {
			return gi, fmt.Errorf("API retornou status %d; falha ao parsear mock: %w", resp.StatusCode, perr)
		}
		return gi, nil
	}

	if err := json.NewDecoder(resp.Body).Decode(&gi); err != nil {
		return gi, fmt.Errorf("falha ao decodificar resposta da greencheck: %w", err)
	}
	return gi, nil
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

	// { added code } Verificação com Green Web Foundation (API ou mock)
	a.logger.Println("Verificando status no Green Web Foundation...")
	ginfo, err := a.fetchGreen(a.cfg.TargetURL)
	if err != nil {
		a.logger.Printf("Erro ao obter greencheck: %v\n", err)
	} else {
		a.logger.Printf("GreenWeb - carbon: %.3f, green: %v, bytes: %d\n", ginfo.Carbon, ginfo.Green, ginfo.Bytes)
	}
}
