package domain

import (
	"testing"
	"time"
)

func TestNewSite(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"valid url", "https://example.com", false},
		{"invalid url", "not-a-url", true},
		{"empty url", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewSite(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSite() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMetrics_IsEcoFriendly(t *testing.T) {
	tests := []struct {
		name     string
		metrics  Metrics
		expected bool
	}{
		{
			"eco friendly site",
			Metrics{CarbonPerVisit: 0.5, GreenHost: true},
			true,
		},
		{
			"high carbon site",
			Metrics{CarbonPerVisit: 1.5, GreenHost: true},
			false,
		},
		{
			"non green host",
			Metrics{CarbonPerVisit: 0.5, GreenHost: false},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.metrics.IsEcoFriendly(); got != tt.expected {
				t.Errorf("IsEcoFriendly() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestReport_IsRecent(t *testing.T) {
	now := time.Now()
	oldDate := now.Add(-8 * 24 * time.Hour)

	tests := []struct {
		name     string
		report   Report
		expected bool
	}{
		{
			"recent report",
			Report{Date: now},
			true,
		},
		{
			"old report",
			Report{Date: oldDate},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.report.IsRecent(); got != tt.expected {
				t.Errorf("IsRecent() = %v, want %v", got, tt.expected)
			}
		})
	}
}
