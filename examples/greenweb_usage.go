package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pamelamiranda/eco-link/internal/adapters/httpclient"
)

func ExampleGreenweb() {
	fmt.Println("Exemplo de uso da API Green Web Foundation")
	fmt.Println("================================================")

	// Exemplo 1: Modo produção (API real)
	fmt.Println("\n=== Exemplo 1: Modo Produção (API Real) ===")
	os.Setenv("ENV", "production") // Forçar modo produção
	adapter := httpclient.NewGreenWebAdapter(10 * time.Second)

	urls := []string{
		"https://example.com",
		"https://github.com",
		"https://google.com",
	}

	for _, url := range urls {
		isGreen, err := adapter.CheckGreenHost(url)
		if err != nil {
			log.Printf("Erro ao verificar %s: %v", url, err)
		} else {
			status := "Não Green"
			if isGreen {
				status = "Green"
			}
			fmt.Printf("%s: %s\n", url, status)
		}
	}

	// Exemplo 2: Modo desenvolvimento (Mock)
	fmt.Println("\n=== Exemplo 2: Modo Desenvolvimento (Mock) ===")
	os.Setenv("ENV", "dev") // Forçar modo desenvolvimento
	adapterDev := httpclient.NewGreenWebAdapter(10 * time.Second)

	for _, url := range urls {
		isGreen, err := adapterDev.CheckGreenHost(url)
		if err != nil {
			log.Printf("Erro ao verificar %s: %v", url, err)
		} else {
			status := "Não Green"
			if isGreen {
				status = "Green"
			}
			fmt.Printf("%s: %s (Mock)\n", url, status)
		}
	}

	// Exemplo 3: Mock forçado
	fmt.Println("\n=== Exemplo 3: Mock Forçado ===")
	os.Setenv("GREENWEB_MOCK", "true")
	adapterMock := httpclient.NewGreenWebAdapter(10 * time.Second)

	for _, url := range urls {
		isGreen, err := adapterMock.CheckGreenHost(url)
		if err != nil {
			log.Printf("Erro ao verificar %s: %v", url, err)
		} else {
			status := "Não Green"
			if isGreen {
				status = "Green"
			}
			fmt.Printf("%s: %s (Mock Forçado)\n", url, status)
		}
	}

	// Exemplo 4: Adapter com mock customizado
	fmt.Println("\n=== Exemplo 4: Mock Customizado ===")
	adapterCustom := httpclient.NewGreenWebAdapterWithMock(10*time.Second, "pkg/web/mock.json")

	for _, url := range urls {
		isGreen, err := adapterCustom.CheckGreenHost(url)
		if err != nil {
			log.Printf("Erro ao verificar %s: %v", url, err)
		} else {
			status := "Não Green"
			if isGreen {
				status = "Green"
			}
			fmt.Printf("%s: %s (Mock Customizado)\n", url, status)
		}
	}

	// Exemplo 5: Demonstração de fallback
	fmt.Println("\n=== Exemplo 5: Fallback para Mock ===")
	fmt.Println("Quando a API real falha, automaticamente usa mock como fallback")

	// Simular falha da API (URL inválida)
	adapterFallback := httpclient.NewGreenWebAdapter(10 * time.Second)
	adapterFallback.SetBaseURL("https://api-inexistente.com") // usar setter público

	isGreen, err := adapterFallback.CheckGreenHost("https://example.com")
	if err != nil {
		log.Printf("Erro: %v", err)
	} else {
		status := "Não Green"
		if isGreen {
			status = "Green"
		}
		fmt.Printf("Fallback funcionou: %s\n", status)
	}

	fmt.Println("\nResumo das funcionalidades:")
	fmt.Println("- Chama API real da Green Web Foundation")
	fmt.Println("- Retorna booleano (true/false) para energia renovável")
	fmt.Println("- Integrado ao relatório final")
	fmt.Println("- Usa mock local quando API falha ou ENV=dev")
	fmt.Println("- Permite desenvolvimento offline")
	fmt.Println("- Fallback automático para mock em caso de erro")
}
