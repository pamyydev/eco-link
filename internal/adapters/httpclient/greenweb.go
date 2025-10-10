package httpclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type GreenWebResponse struct {
	IsGreen bool `json:"isGreen"`
}

type GreenWebAdapter struct {
	client  *http.Client
	baseURL string
	logger  *Logger
}

func NewGreenWebAdapter(timeout time.Duration) *GreenWebAdapter {
	return &GreenWebAdapter{
		client: &http.Client{
			Timeout: timeout,
		},
		baseURL: "https://api.greenweb.org",
		logger:  NewLogger(),
	}
}

func (g *GreenWebAdapter) CheckGreenHost(urlStr string) (bool, error) {
	startTime := time.Now()
	g.logger.LogAPICall("Green Web", urlStr, startTime)

	// Validação de entrada
	if urlStr == "" {
		err := fmt.Errorf("URL não pode estar vazia")
		g.logger.LogValidationError("URL", urlStr, err)
		return false, err
	}

	// Validação da URL
	if _, err := url.ParseRequestURI(urlStr); err != nil {
		g.logger.LogValidationError("URL", urlStr, err)
		return false, fmt.Errorf("URL inválida: %w", err)
	}

	// Construção da URL da API com encoding adequado
	apiURL := fmt.Sprintf("%s/check?url=%s", g.baseURL, url.QueryEscape(urlStr))

	resp, err := g.client.Get(apiURL)
	if err != nil {
		g.logger.LogAPIFailure("Green Web", urlStr, err, time.Since(startTime))
		return false, fmt.Errorf("falha ao chamar Green Web API para URL '%s': %w", urlStr, err)
	}
	defer resp.Body.Close()

	g.logger.LogResponse(urlStr, resp.StatusCode, time.Since(startTime))

	// Verificação de status codes específicos
	switch resp.StatusCode {
	case http.StatusOK:
		// Continua o processamento
	case http.StatusBadRequest:
		err := fmt.Errorf("requisição inválida para URL '%s' (status 400)", urlStr)
		g.logger.LogAPIFailure("Green Web", urlStr, err, time.Since(startTime))
		return false, err
	case http.StatusNotFound:
		err := fmt.Errorf("URL '%s' não encontrada na API Green Web (status 404)", urlStr)
		g.logger.LogAPIFailure("Green Web", urlStr, err, time.Since(startTime))
		return false, err
	case http.StatusTooManyRequests:
		err := fmt.Errorf("limite de requisições excedido para URL '%s' (status 429)", urlStr)
		g.logger.LogAPIFailure("Green Web", urlStr, err, time.Since(startTime))
		return false, err
	case http.StatusInternalServerError:
		err := fmt.Errorf("erro interno do servidor Green Web para URL '%s' (status 500)", urlStr)
		g.logger.LogAPIFailure("Green Web", urlStr, err, time.Since(startTime))
		return false, err
	default:
		err := fmt.Errorf("API Green Web retornou status inesperado %d para URL '%s'", resp.StatusCode, urlStr)
		g.logger.LogAPIFailure("Green Web", urlStr, err, time.Since(startTime))
		return false, err
	}

	var result GreenWebResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		g.logger.LogAPIFailure("Green Web", urlStr, err, time.Since(startTime))
		return false, fmt.Errorf("erro ao decodificar resposta JSON da API Green Web para URL '%s': %w", urlStr, err)
	}

	g.logger.LogDataValidation("Green Web", result)
	g.logger.LogAPISuccess("Green Web", urlStr, time.Since(startTime))

	return result.IsGreen, nil
}
