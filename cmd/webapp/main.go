package main

import (
	"log"
	"os"

	"github.com/pamelamiranda/eco-link/internal/adapters"
	"github.com/pamelamiranda/eco-link/internal/adapters/httpclient"
	"github.com/pamelamiranda/eco-link/internal/app"
)

func main() {
	logger := log.New(os.Stdout, "[eco-link] ", log.LstdFlags)

	cfg := app.NewConfig()

	carbonClient := httpclient.NewWebsiteCarbonAdapter(cfg.Timeout)
	greenClient := httpclient.NewGreenWebAdapter(cfg.Timeout)

	service := adapters.NewService(carbonClient, greenClient)

	application := app.NewApplication(logger, cfg, service)

	// Start BFF server (serves API and optional frontend static files)
	if err := application.StartServer(); err != nil {
		logger.Fatalf("server error: %v", err)
	}
}
