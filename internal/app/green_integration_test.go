package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"
)

// Test direto que chama a API oficial da Green Web Foundation (sem fallback).
// Executa apenas quando RUN_INTEGRATION_TESTS=1.
func TestGreenWeb_API_DirectIntegration(t *testing.T) {
    if os.Getenv("RUN_INTEGRATION_TESTS") != "1" {
        t.Skip("skipping integration test; set RUN_INTEGRATION_TESTS=1 to enable")
    }

    target := "https://example.com"
    // Endpoint correto inclui /api/ conforme documentação
    base := "https://api.thegreenwebfoundation.org/api/v3/greencheck/"
    // A API geralmente espera o domínio/hostname sem o esquema. Extrair hostname
    parsed, err := url.Parse(target)
    if err != nil {
        t.Fatalf("failed to parse target URL: %v", err)
    }
    host := parsed.Hostname()
    reqURL := base + url.PathEscape(host)

    client := http.Client{Timeout: 15 * time.Second}
    resp, err := client.Get(reqURL)
    if err != nil {
        t.Fatalf("failed to call Green Web API: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        t.Fatalf("unexpected status code from Green Web API: %d", resp.StatusCode)
    }

    // Tentar decodificar resposta; o campo esperado normalmente é 'green'.
    var body map[string]any
    if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
        t.Fatalf("failed to decode JSON from Green Web API: %v", err)
    }

    // Logar a resposta completa para inspeção
    t.Logf("Green Web API response: %+v", body)

    // Verificar presença do campo 'green' como sanity check
    if _, ok := body["green"]; !ok {
        // Alguns endpoints podem usar outras chaves; aceitaremos se houver qualquer chave booleana
        foundBool := false
        for k, v := range body {
            switch v.(type) {
            case bool:
                t.Logf("found boolean key '%s' in response", k)
                foundBool = true
            }
        }
        if !foundBool {
            t.Fatalf("response did not contain expected boolean 'green' field and no boolean fields found")
        }
    }

    fmt.Println("Green Web API direct integration ok")
}
