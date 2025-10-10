package httpclient

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestWebsiteCarbonAdapter_GetMetrics_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verificar se a URL contém o parâmetro correto
		if !strings.Contains(r.URL.RawQuery, "url=") {
			t.Errorf("Expected URL parameter in request, got: %s", r.URL.RawQuery)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"bytes": 1000000,
			"green": true,
			"cO2": 0.42,
			"rating": "A",
			"cleanerThan": 0.9
		}`))
	}))
	defer server.Close()

	adapter := NewWebsiteCarbonAdapter(5 * time.Second)
	adapter.baseURL = server.URL

	metrics, err := adapter.GetMetrics("https://example.com")
	if err != nil {
		t.Errorf("GetMetrics() error = %v", err)
		return
	}

	if metrics.CarbonPerVisit != 0.42 {
		t.Errorf("GetMetrics() CarbonPerVisit = %v, want %v", metrics.CarbonPerVisit, 0.42)
	}
}

func TestWebsiteCarbonAdapter_GetMetrics_EmptyURL(t *testing.T) {
	adapter := NewWebsiteCarbonAdapter(5 * time.Second)

	_, err := adapter.GetMetrics("")
	if err == nil {
		t.Error("Expected error for empty URL")
	}

	if !strings.Contains(err.Error(), "URL não pode estar vazia") {
		t.Errorf("Expected error about empty URL, got: %v", err)
	}
}

func TestWebsiteCarbonAdapter_GetMetrics_InvalidURL(t *testing.T) {
	adapter := NewWebsiteCarbonAdapter(5 * time.Second)

	_, err := adapter.GetMetrics("invalid-url")
	if err == nil {
		t.Error("Expected error for invalid URL")
	}

	if !strings.Contains(err.Error(), "URL inválida") {
		t.Errorf("Expected error about invalid URL, got: %v", err)
	}
}

func TestWebsiteCarbonAdapter_GetMetrics_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
	}))
	defer server.Close()

	adapter := NewWebsiteCarbonAdapter(5 * time.Second)
	adapter.baseURL = server.URL

	_, err := adapter.GetMetrics("https://example.com")
	if err == nil {
		t.Error("Expected error for HTTP 400")
	}

	if !strings.Contains(err.Error(), "requisição inválida") {
		t.Errorf("Expected error about bad request, got: %v", err)
	}
}

func TestWebsiteCarbonAdapter_GetMetrics_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`invalid json`))
	}))
	defer server.Close()

	adapter := NewWebsiteCarbonAdapter(5 * time.Second)
	adapter.baseURL = server.URL

	_, err := adapter.GetMetrics("https://example.com")
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}

	if !strings.Contains(err.Error(), "erro ao decodificar resposta JSON") {
		t.Errorf("Expected error about JSON decoding, got: %v", err)
	}
}

func TestWebsiteCarbonAdapter_GetMetrics_NegativeValues(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"bytes": -1000,
			"green": true,
			"gco2e": -0.45,
			"rating": "A",
			"cleanerThan": 0.9
		}`))
	}))
	defer server.Close()

	adapter := NewWebsiteCarbonAdapter(5 * time.Second)
	adapter.baseURL = server.URL

	_, err := adapter.GetMetrics("https://example.com")
	if err == nil {
		t.Error("Expected error for negative values")
	}

	if !strings.Contains(err.Error(), "valor de CO2 inválido") && !strings.Contains(err.Error(), "valor de bytes inválido") {
		t.Errorf("Expected error about invalid values, got: %v", err)
	}
}

func TestGreenWebAdapter_CheckGreenHost_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"isGreen": true}`))
	}))
	defer server.Close()

	adapter := NewGreenWebAdapter(5 * time.Second)
	adapter.baseURL = server.URL

	isGreen, err := adapter.CheckGreenHost("https://example.com")
	if err != nil {
		t.Errorf("CheckGreenHost() error = %v", err)
		return
	}

	if !isGreen {
		t.Error("CheckGreenHost() expected true")
	}
}

func TestGreenWebAdapter_CheckGreenHost_EmptyURL(t *testing.T) {
	adapter := NewGreenWebAdapter(5 * time.Second)

	_, err := adapter.CheckGreenHost("")
	if err == nil {
		t.Error("Expected error for empty URL")
	}

	if !strings.Contains(err.Error(), "URL não pode estar vazia") {
		t.Errorf("Expected error about empty URL, got: %v", err)
	}
}

func TestGreenWebAdapter_CheckGreenHost_InvalidURL(t *testing.T) {
	adapter := NewGreenWebAdapter(5 * time.Second)

	_, err := adapter.CheckGreenHost("invalid-url")
	if err == nil {
		t.Error("Expected error for invalid URL")
	}

	if !strings.Contains(err.Error(), "URL inválida") {
		t.Errorf("Expected error about invalid URL, got: %v", err)
	}
}

func TestGreenWebAdapter_CheckGreenHost_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
	}))
	defer server.Close()

	adapter := NewGreenWebAdapter(5 * time.Second)
	adapter.baseURL = server.URL

	_, err := adapter.CheckGreenHost("https://example.com")
	if err == nil {
		t.Error("Expected error for HTTP 404")
	}

	if !strings.Contains(err.Error(), "não encontrada na API Green Web") {
		t.Errorf("Expected error about not found, got: %v", err)
	}
}

func TestGreenWebAdapter_CheckGreenHost_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`invalid json`))
	}))
	defer server.Close()

	adapter := NewGreenWebAdapter(5 * time.Second)
	adapter.baseURL = server.URL

	_, err := adapter.CheckGreenHost("https://example.com")
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}

	if !strings.Contains(err.Error(), "erro ao decodificar resposta JSON") {
		t.Errorf("Expected error about JSON decoding, got: %v", err)
	}
}

func TestWebsiteCarbonAdapter_CalculateEstimatedCO2(t *testing.T) {
	adapter := NewWebsiteCarbonAdapter(5 * time.Second)

	tests := []struct {
		name          string
		pageSizeBytes int
		expectedCO2   float64
	}{
		{
			name:          "Small page (100KB)",
			pageSizeBytes: 102400,
			expectedCO2:   0.05, // 100KB * 0.0005 = 0.05g CO₂
		},
		{
			name:          "Medium page (1MB)",
			pageSizeBytes: 1048576,
			expectedCO2:   0.512, // 1024KB * 0.0005 = 0.512g CO₂
		},
		{
			name:          "Large page (5MB)",
			pageSizeBytes: 5242880,
			expectedCO2:   2.56, // 5120KB * 0.0005 = 2.56g CO₂
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := adapter.CalculateEstimatedCO2(tt.pageSizeBytes)
			if result != tt.expectedCO2 {
				t.Errorf("CalculateEstimatedCO2() = %v, want %v", result, tt.expectedCO2)
			}
		})
	}
}

func TestWebsiteCarbonAdapter_GetMetricsWithEstimation(t *testing.T) {
	t.Run("API returns CO2 data", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{
				"bytes": 1000000,
				"green": true,
				"cO2": 0.42,
				"rating": "A",
				"cleanerThan": 0.9
			}`))
		}))
		defer server.Close()

		adapter := NewWebsiteCarbonAdapter(5 * time.Second)
		adapter.baseURL = server.URL

		metrics, err := adapter.GetMetricsWithEstimation("https://example.com")
		if err != nil {
			t.Errorf("GetMetricsWithEstimation() error = %v", err)
			return
		}

		// Deve usar o valor da API, não a estimativa
		if metrics.CarbonPerVisit != 0.42 {
			t.Errorf("GetMetricsWithEstimation() CarbonPerVisit = %v, want %v", metrics.CarbonPerVisit, 0.42)
		}
	})

	t.Run("API returns zero CO2, should estimate", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{
				"bytes": 102400,
				"green": true,
				"cO2": 0,
				"rating": "A",
				"cleanerThan": 0.9
			}`))
		}))
		defer server.Close()

		adapter := NewWebsiteCarbonAdapter(5 * time.Second)
		adapter.baseURL = server.URL

		metrics, err := adapter.GetMetricsWithEstimation("https://example.com")
		if err != nil {
			t.Errorf("GetMetricsWithEstimation() error = %v", err)
			return
		}

		// Deve calcular estimativa: 100KB * 0.0005 = 0.05g CO₂
		expectedCO2 := 0.05
		if metrics.CarbonPerVisit != expectedCO2 {
			t.Errorf("GetMetricsWithEstimation() CarbonPerVisit = %v, want %v", metrics.CarbonPerVisit, expectedCO2)
		}
	})
}
