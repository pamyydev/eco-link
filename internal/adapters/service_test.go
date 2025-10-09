package adapters

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockCarbonClient implements WebsiteCarbonClient for testing
type mockCarbonClient struct{}

func (m *mockCarbonClient) GetMetrics(url string) (domain.Metrics, error) {
	return domain.Metrics{
		CarbonPerVisit: 0.45,
		BytesPerVisit:  1000,
	}, nil
}

// mockGreenClient implements GreenWebClient for testing
type mockGreenClient struct{}

func (m *mockGreenClient) CheckGreenHost(url string) (bool, error) {
	return true, nil
}

func TestService_Analyze(t *testing.T) {
	service := NewService(&mockCarbonClient{}, &mockGreenClient{})

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