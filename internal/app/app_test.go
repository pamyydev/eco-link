package app

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCarbonClient_GetMetrics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"carbonPerVisit": 0.5, "bytesPerVisit": 1000}`))
	}))
	defer server.Close()

	client := &carbonClient{
		client:  http.DefaultClient,
		baseURL: server.URL,
	}

	metrics, err := client.GetMetrics("https://example.com")
	if err != nil {
		t.Errorf("GetMetrics() error = %v", err)
		return
	}

	if metrics.CarbonPerVisit != 0.5 {
		t.Errorf("GetMetrics() carbonPerVisit = %v, want %v", metrics.CarbonPerVisit, 0.5)
	}
}

func TestGreenClient_CheckGreenHost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"isGreen": true}`))
	}))
	defer server.Close()

	client := &greenClient{
		client:  http.DefaultClient,
		baseURL: server.URL,
	}

	isGreen, err := client.CheckGreenHost("https://example.com")
	if err != nil {
		t.Errorf("CheckGreenHost() error = %v", err)
		return
	}

	if !isGreen {
		t.Error("CheckGreenHost() expected true")
	}
}