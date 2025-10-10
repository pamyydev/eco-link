package httpclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/pamelamiranda/eco-link/internal/domain"
)

type WebsiteCarbonResponse struct {
	Bytes       int     `json:"bytes"`
	Green       bool    `json:"green"`
	CO2         float64 `json:"gco2e"`
	Rating      string  `json:"rating"`
	CleanerThan float64 `json:"cleanerThan"`
}

type WebsiteCarbonAdapter struct {
	client  *http.Client
	baseURL string
	logger  *Logger
}

func NewWebsiteCarbonAdapter(timeout time.Duration) *WebsiteCarbonAdapter {
	return &WebsiteCarbonAdapter{
		client: &http.Client{
			Timeout: timeout,
		},
		baseURL: "https://api.websitecarbon.com",
		logger:  NewLogger(),
	}
}

func (w *WebsiteCarbonAdapter) GetMetrics(urlStr string) (domain.Metrics, error) {
	startTime := time.Now()
	w.logger.LogAPICall("Website Carbon", urlStr, startTime)

	// Validação de entrada
	if urlStr == "" {
		err := fmt.Errorf("URL não pode estar vazia")
		w.logger.LogValidationError("URL", urlStr, err)
		return domain.Metrics{}, err
	}

	// Validação da URL
	if _, err := url.ParseRequestURI(urlStr); err != nil {
		w.logger.LogValidationError("URL", urlStr, err)
		return domain.Metrics{}, fmt.Errorf("URL inválida: %w", err)
	}

	// Construção da URL da API com parâmetros
	apiURL := fmt.Sprintf("%s/data?bytes=1000000&green=1", w.baseURL)

	resp, err := w.client.Get(apiURL)
	if err != nil {
		w.logger.LogAPIFailure("Website Carbon", urlStr, err, time.Since(startTime))
		return domain.Metrics{}, fmt.Errorf("falha ao chamar Website Carbon API para URL '%s': %w", urlStr, err)
	}
	defer resp.Body.Close()

	w.logger.LogResponse(urlStr, resp.StatusCode, time.Since(startTime))

	// Verificação de status codes específicos
	switch resp.StatusCode {
	case http.StatusOK:
		// Continua o processamento
	case http.StatusBadRequest:
		err := fmt.Errorf("requisição inválida para URL '%s' (status 400)", urlStr)
		w.logger.LogAPIFailure("Website Carbon", urlStr, err, time.Since(startTime))
		return domain.Metrics{}, err
	case http.StatusNotFound:
		err := fmt.Errorf("URL '%s' não encontrada na API (status 404)", urlStr)
		w.logger.LogAPIFailure("Website Carbon", urlStr, err, time.Since(startTime))
		return domain.Metrics{}, err
	case http.StatusTooManyRequests:
		err := fmt.Errorf("limite de requisições excedido para URL '%s' (status 429)", urlStr)
		w.logger.LogAPIFailure("Website Carbon", urlStr, err, time.Since(startTime))
		return domain.Metrics{}, err
	case http.StatusInternalServerError:
		err := fmt.Errorf("erro interno do servidor Website Carbon para URL '%s' (status 500)", urlStr)
		w.logger.LogAPIFailure("Website Carbon", urlStr, err, time.Since(startTime))
		return domain.Metrics{}, err
	default:
		err := fmt.Errorf("API Website Carbon retornou status inesperado %d para URL '%s'", resp.StatusCode, urlStr)
		w.logger.LogAPIFailure("Website Carbon", urlStr, err, time.Since(startTime))
		return domain.Metrics{}, err
	}

	var result WebsiteCarbonResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		w.logger.LogAPIFailure("Website Carbon", urlStr, err, time.Since(startTime))
		return domain.Metrics{}, fmt.Errorf("erro ao decodificar resposta JSON da API Website Carbon para URL '%s': %w", urlStr, err)
	}

	w.logger.LogDataValidation("Website Carbon", result)

	// Validação dos dados recebidos
	if result.CO2 < 0 {
		w.logger.LogInvalidData("Website Carbon", "CO2", result.CO2, "negative value")
		return domain.Metrics{}, fmt.Errorf("valor de CO2 inválido recebido da API: %f", result.CO2)
	}

	if result.Bytes < 0 {
		w.logger.LogInvalidData("Website Carbon", "Bytes", result.Bytes, "negative value")
		return domain.Metrics{}, fmt.Errorf("valor de bytes inválido recebido da API: %d", result.Bytes)
	}

	w.logger.LogAPISuccess("Website Carbon", urlStr, time.Since(startTime))

	return domain.Metrics{
		CarbonPerVisit: result.CO2,
		BytesPerVisit:  result.Bytes,
		GreenHost:      result.Green,
	}, nil
}
