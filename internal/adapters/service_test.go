package adapters

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupMockServers(t *testing.T) (*httptest.Server, *httptest.Server, *httptest.Server) {
	// Mock server for the website being analyzed
	siteServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Mock website content"))
	}))

	// Mock server for metrics API
	metricsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"site":"%s","carbonPerVisit":0.45,"greenHost":true}`, r.URL.Query().Get("url"))
	}))

	// Mock server for green host check
	greenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"isGreen":true}`))
	}))

	return siteServer, metricsServer, greenServer
}

func TestService_Analyze(t *testing.T) {
	siteServer, metricsServer, greenServer := setupMockServers(t)
	defer siteServer.Close()
	defer metricsServer.Close()
	defer greenServer.Close()

	service := NewService(metricsServer.URL, greenServer.URL)

	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "valid url",
			url:     siteServer.URL,
			wantErr: false,
		},
		{
			name:    "invalid url",
			url:     "not-a-url",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report, err := service.Analyze(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Analyze() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if report.Site.URL != tt.url {
					t.Errorf("Analyze() got URL = %v, want %v", report.Site.URL, tt.url)
				}
				if !report.Metrics.GreenHost {
					t.Error("Analyze() expected GreenHost to be true")
				}
			}
		})
	}
}