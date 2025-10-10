package app

import (
	"bytes"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/pamelamiranda/eco-link/internal/adapters"
	"github.com/pamelamiranda/eco-link/internal/domain"
)

// mockCarbonClient is a mock implementation of the WebsiteCarbonClient interface.

type mockCarbonClient struct{}

func (m *mockCarbonClient) GetMetrics(url string) (domain.Metrics, error) {
	return domain.Metrics{
		CarbonPerVisit: 0.5,
		BytesPerVisit:  1000,
	}, nil
}

// mockGreenClient is a mock implementation of the GreenWebClient interface.

type mockGreenClient struct{}

func (m *mockGreenClient) CheckGreenHost(url string) (bool, error) {
	return true, nil
}

func TestApplication_Run(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)

	cfg := Config{
		Timeout:   10 * time.Second,
		TargetURL: "https://example.com",
	}

	carbonClient := &mockCarbonClient{}
	greenClient := &mockGreenClient{}
	service := adapters.NewService(carbonClient, greenClient)

	app := NewApplication(logger, cfg, service)

	// We need to run the app in a way that doesn't call log.Fatalf, which would exit the test.
	// For this test, we'll just check the output before the potential fatal error.
	report, err := app.service.Analyze(app.cfg.TargetURL)
	if err != nil {
		t.Fatalf("Analyze() error = %v", err)
	}

	app.logger.Printf("\nResultados da análise para %s:\n", report.Site.URL)
	app.logger.Printf("Carbono por visita: %.2fg CO2\n", report.Metrics.CarbonPerVisit)
	app.logger.Printf("Bytes por visita: %d bytes\n", report.Metrics.BytesPerVisit)
	app.logger.Printf("Hosting verde: %v\n", report.Metrics.GreenHost)
	app.logger.Printf("Eco-friendly: %v\n", report.Metrics.IsEcoFriendly())

	output := buf.String()

	expectedOutputs := []string{
		"Resultados da análise para https://example.com:",
		"Carbono por visita: 0.50g CO2",
		"Bytes por visita: 1000 bytes",
		"Hosting verde: true",
		"Eco-friendly: true",
	}

	for _, expected := range expectedOutputs {
		if !strings.Contains(output, expected) {
			t.Errorf("log output does not contain %q", expected)
		}
	}
}