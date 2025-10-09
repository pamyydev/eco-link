package app

import (
	"fmt"
	"github.com/pamelamiranda/eco-link/internal/adapters"
)

func Run() {
	fmt.Println("eco-link app iniciado!")
	
	// Initialize analysis service
	service := adapters.NewService(
		"https://api.websitecarbon.com", // Replace with actual API endpoint
		"https://api.green-host.com",    // Replace with actual green host API
	)

	// Example: Analyze a website
	report, err := service.Analyze("https://example.com")
	if err != nil {
		fmt.Printf("Error analyzing website: %v\n", err)
		return
	}

	// Print analysis results
	fmt.Printf("\nAnalysis Results for %s:\n", report.Site.URL)
	fmt.Printf("Carbon per visit: %.2fg CO2\n", report.Metrics.CarbonPerVisit)
	fmt.Printf("Bytes per visit: %d bytes\n", report.Metrics.BytesPerVisit)
	fmt.Printf("Green hosting: %v\n", report.Metrics.GreenHost)
	fmt.Printf("Eco-friendly: %v\n", report.Metrics.IsEcoFriendly())
}
