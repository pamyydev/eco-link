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