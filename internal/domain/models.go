package domain

import (
	"errors"
	"net/url"
	"time"
)

// Site represents a website to be analyzed
type Site struct {
	URL string
}

// NewSite creates a new Site with URL validation
func NewSite(urlStr string) (*Site, error) {
	if _, err := url.ParseRequestURI(urlStr); err != nil {
		return nil, errors.New("invalid URL")
	}
	return &Site{URL: urlStr}, nil
}

// Metrics represents sustainability metrics for a website
type Metrics struct {
	CarbonPerVisit float64 // in grams of CO2
	BytesPerVisit  int     // in bytes
	GreenHost      bool    // whether the host uses renewable energy
}

// IsEcoFriendly returns true if the site is considered ecological
func (m Metrics) IsEcoFriendly() bool {
	return m.GreenHost && m.CarbonPerVisit < 1.0 // less than 1g CO2 per visit
}

// Report represents a complete website analysis report
type Report struct {
	Site    Site
	Metrics Metrics
	Date    time.Time
}

// NewReport creates a new report with current timestamp
func NewReport(site Site, metrics Metrics) Report {
	return Report{
		Site:    site,
		Metrics: metrics,
		Date:    time.Now(),
	}
}

// IsRecent returns true if the report is from the last 7 days
func (r Report) IsRecent() bool {
	return time.Since(r.Date).Hours() < 24*7
}
