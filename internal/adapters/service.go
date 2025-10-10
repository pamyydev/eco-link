package adapters

import (
	"fmt"
	"time"

	"github.com/pamelamiranda/eco-link/internal/adapters/httpclient"
	"github.com/pamelamiranda/eco-link/internal/domain"
	"github.com/pamelamiranda/eco-link/internal/ports"
)

// Service handles website sustainability analysis
type Service struct {
	carbonClient ports.WebsiteCarbonClient
	greenClient  ports.GreenWebClient
	maxRetries   int
	retryDelay   time.Duration
	logger       *httpclient.Logger
}

// NewService creates a new analysis service
func NewService(carbonClient ports.WebsiteCarbonClient, greenClient ports.GreenWebClient) *Service {
	return &Service{
		carbonClient: carbonClient,
		greenClient:  greenClient,
		maxRetries:   3,
		retryDelay:   time.Second * 2,
		logger:       httpclient.NewLogger(),
	}
}

// NewServiceWithRetry creates a new analysis service with custom retry configuration
func NewServiceWithRetry(carbonClient ports.WebsiteCarbonClient, greenClient ports.GreenWebClient, maxRetries int, retryDelay time.Duration) *Service {
	return &Service{
		carbonClient: carbonClient,
		greenClient:  greenClient,
		maxRetries:   maxRetries,
		retryDelay:   retryDelay,
		logger:       httpclient.NewLogger(),
	}
}

// isRetryableError checks if an error is retryable based on HTTP status codes
func isRetryableError(err error) bool {
	// Verifica se é um erro de timeout ou conexão
	if err == nil {
		return false
	}

	// Para erros de rede, sempre tenta novamente
	return true
}

// retryWithBackoff executes a function with exponential backoff retry logic
func (s *Service) retryWithBackoff(operation func() error, operationName, url string) error {
	var lastErr error

	for attempt := 0; attempt <= s.maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff: 2^attempt * retryDelay
			delay := time.Duration(1<<uint(attempt-1)) * s.retryDelay
			s.logger.LogRetry(operationName, url, attempt, s.maxRetries, lastErr)
			time.Sleep(delay)
		}

		err := operation()
		if err == nil {
			return nil
		}

		lastErr = err

		// Se não é um erro retryable, para imediatamente
		if !isRetryableError(err) {
			break
		}

		// Se é a última tentativa, não continua
		if attempt == s.maxRetries {
			break
		}
	}

	s.logger.LogRetryExhausted(operationName, url, s.maxRetries, lastErr)
	return fmt.Errorf("%s falhou após %d tentativas: %w", operationName, s.maxRetries+1, lastErr)
}

// Analyze fetches and processes sustainability metrics for a given URL
func (s *Service) Analyze(url string) (domain.Report, error) {
	startTime := time.Now()
	s.logger.LogServiceOperation("website analysis", url, startTime)

	// Validação de entrada
	if url == "" {
		err := fmt.Errorf("URL não pode estar vazia")
		s.logger.LogServiceCompletion("website analysis", url, false, time.Since(startTime))
		return domain.Report{}, err
	}

	site, err := domain.NewSite(url)
	if err != nil {
		s.logger.LogServiceCompletion("website analysis", url, false, time.Since(startTime))
		return domain.Report{}, fmt.Errorf("URL inválida '%s': %w", url, err)
	}

	// Busca métricas de carbono com retry
	var metrics domain.Metrics
	err = s.retryWithBackoff(func() error {
		var err error
		metrics, err = s.carbonClient.GetMetrics(url)
		return err
	}, "obtenção de métricas de carbono", url)

	if err != nil {
		s.logger.LogServiceCompletion("website analysis", url, false, time.Since(startTime))
		return domain.Report{}, fmt.Errorf("falha ao obter métricas de carbono para URL '%s': %w", url, err)
	}

	// Verifica se é green host com retry
	var isGreen bool
	err = s.retryWithBackoff(func() error {
		var err error
		isGreen, err = s.greenClient.CheckGreenHost(url)
		return err
	}, "verificação de green host", url)

	if err != nil {
		s.logger.LogServiceCompletion("website analysis", url, false, time.Since(startTime))
		return domain.Report{}, fmt.Errorf("falha ao verificar green host para URL '%s': %w", url, err)
	}

	// Atualiza as métricas com a informação do green host
	metrics.GreenHost = isGreen
	report := domain.NewReport(*site, metrics)

	s.logger.LogServiceCompletion("website analysis", url, true, time.Since(startTime))
	return report, nil
}
