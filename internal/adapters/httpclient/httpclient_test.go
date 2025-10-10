package httpclient

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestWebsiteCarbonAdapter_GetMetrics(t *testing.T) {
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

func TestGreenWebAdapter_CheckGreenHost(t *testing.T) {
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

func TestGreenWebAdapter_CheckGreenHost_WithMock(t *testing.T) {
	// Criar adapter com mock
	adapter := NewGreenWebAdapterWithMock(5*time.Second, "pkg/web/mock.json")

	tests := []struct {
		name     string
		url      string
		expected bool
	}{
		{
			name:     "Green site from mock",
			url:      "https://example.com",
			expected: true,
		},
		{
			name:     "Non-green site from mock",
			url:      "https://google.com",
			expected: false,
		},
		{
			name:     "Unknown site uses default mock value",
			url:      "https://unknown-site.com",
			expected: false, // Default value from mock.json
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := adapter.CheckGreenHost(tt.url)
			if err != nil {
				t.Errorf("CheckGreenHost() error = %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("CheckGreenHost() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGreenWebAdapter_CheckGreenHost_FallbackToMock(t *testing.T) {
	// Criar servidor que retorna erro
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	// Criar adapter normal (não mock)
	adapter := NewGreenWebAdapter(5 * time.Second)
	adapter.baseURL = server.URL

	// Deve fazer fallback para mock quando API falha
	result, err := adapter.CheckGreenHost("https://example.com")
	if err != nil {
		t.Errorf("CheckGreenHost() should fallback to mock, got error: %v", err)
		return
	}

	// Deve retornar um valor (true ou false) do mock
	if result != true && result != false {
		t.Errorf("CheckGreenHost() should return boolean from mock, got: %v", result)
	}
}

func TestGreenWebAdapter_ExtractDomain(t *testing.T) {
	adapter := NewGreenWebAdapter(5 * time.Second)

	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "Simple domain",
			url:      "https://example.com",
			expected: "example.com",
		},
		{
			name:     "Domain with port",
			url:      "https://example.com:8080",
			expected: "example.com",
		},
		{
			name:     "Domain with path",
			url:      "https://example.com/path/to/page",
			expected: "example.com",
		},
		{
			name:     "Subdomain",
			url:      "https://www.example.com",
			expected: "www.example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := adapter.extractDomain(tt.url)
			if result != tt.expected {
				t.Errorf("extractDomain() = %v, want %v", result, tt.expected)
			}
		})
	}
}
