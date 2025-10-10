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
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"bytes": 1000000,
			"green": true,
			"gco2e": 0.45,
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

	if metrics.CarbonPerVisit != 0.45 {
		t.Errorf("GetMetrics() CarbonPerVisit = %v, want %v", metrics.CarbonPerVisit, 0.45)
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
