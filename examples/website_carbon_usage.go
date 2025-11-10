package main

import (
	"fmt"
	"log"
	"time"

	"github.com/pamelamiranda/eco-link/internal/adapters/httpclient"
)

func ExampleWebsiteCarbon() {
	// Criar adapter com timeout de 10 segundos
	adapter := httpclient.NewWebsiteCarbonAdapter(10 * time.Second)

	// Exemplo 1: Obter métricas de um site
	fmt.Println("=== Exemplo 1: Obter métricas de um site ===")
	url := "https://example.com"

	metrics, err := adapter.GetMetrics(url)
	if err != nil {
		log.Printf("Erro ao obter métricas: %v", err)
	} else {
		fmt.Printf("Site: %s\n", url)
		fmt.Printf("CO₂ por visita: %.3f gramas\n", metrics.CarbonPerVisit)
		fmt.Printf("Bytes por visita: %d\n", metrics.BytesPerVisit)
		fmt.Printf("Green Host: %t\n", metrics.GreenHost)
		fmt.Println()
	}

	// Exemplo 2: Calcular CO₂ estimado
	fmt.Println("=== Exemplo 2: Calcular CO₂ estimado ===")
	pageSizes := []int{102400, 1048576, 5242880} // 100KB, 1MB, 5MB
	pageNames := []string{"100KB", "1MB", "5MB"}

	for i, size := range pageSizes {
		estimatedCO2 := adapter.CalculateEstimatedCO2(size)
		fmt.Printf("Página %s (%d bytes): %.3f gramas de CO₂\n",
			pageNames[i], size, estimatedCO2)
	}
	fmt.Println()

	// Exemplo 3: Obter métricas com estimativa como fallback
	fmt.Println("=== Exemplo 3: Métricas com estimativa como fallback ===")
	metricsWithEstimation, err := adapter.GetMetricsWithEstimation(url)
	if err != nil {
		log.Printf("Erro ao obter métricas com estimativa: %v", err)
	} else {
		fmt.Printf("Site: %s\n", url)
		fmt.Printf("CO₂ por visita: %.3f gramas\n", metricsWithEstimation.CarbonPerVisit)
		fmt.Printf("Bytes por visita: %d\n", metricsWithEstimation.BytesPerVisit)
		fmt.Printf("Green Host: %t\n", metricsWithEstimation.GreenHost)
	}

	// Exemplo 4: Demonstração da requisição GET
	fmt.Println("\n=== Exemplo 4: Estrutura da requisição ===")
	fmt.Printf("URL da API: %s/site?url=%s\n",
		"https://api.websitecarbon.com",
		"https%3A%2F%2Fexample.com")
	fmt.Println("Resposta JSON esperada:")
	fmt.Println(`{
  "bytes": 1000000,
  "green": true,
  "cO2": 0.42,
  "rating": "A",
  "cleanerThan": 0.9
}`)
	fmt.Println("\nMapeamento para Metrics:")
	fmt.Println("- cO2 -> CarbonPerVisit")
	fmt.Println("- bytes -> BytesPerVisit")
	fmt.Println("- green -> GreenHost")
}
