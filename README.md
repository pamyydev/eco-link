# eco-link

Analyze website sustainability and carbon footprint. Calculate carbon emissions, page size, and verify green hosting status.

## Tests

```bash
go test ./... -v  # Run all tests with verbose output
```

## Features

- Website carbon footprint analysis
- Page size calculation
- Green hosting verification
- Sustainability scoring
- Recent report tracking

## Quick Start

```bash
# Run the application
go run cmd/webapp/main.go
```

## Mock & configuração da Green Web Foundation

- Para forçar o uso do JSON mock local durante desenvolvimento, exporte:
  - ANALYSIS_ENV=dev
  - Opcional: MOCK_GREENCHECK_FILE=./greencheck_mock.json
- Exemplo de mock (JSON): {"carbon":0.6,"green":false,"bytes":512000}
- Em produção, a aplicação tentará chamar a API: https://api.thegreenwebfoundation.org/v3/greencheck/{url} e cairá no mock se a chamada falhar.

## Frontend (dev / build) and BFF notes

- Requisitos para o frontend: Node.js (recomendo >=18) e npm.
- Desenvolvimento (modo rápido): inicie o frontend (Vite) em seu projeto frontend e rode o backend aqui em outra porta.

Exemplo (frontend dev + BFF):

1. No repositório do frontend (ex: `Downloads/carbon-scan-now-main`):
```
npm install
npm run dev
```

2. No repositório backend (`eco-link`):
```
export PORT=8081
go run cmd/webapp/main.go
```

O BFF (porta 8081 por padrão) provê endpoints:
- GET /api/analyze?url={url}
- GET /api/greencheck?url={url}

Obs: o BFF habilita CORS para `/api` para facilitar o desenvolvimento com Vite (sem proxy).

Build e servir static files pelo BFF

Se preferir que o BFF sirva o frontend (produção simples):

1. No projeto frontend, gere a build (`npm run build`).
2. Copie o `dist` para `eco-link/frontend/dist` ou use o script helper do repo:

```
cd /home/pamy/eco-link
make web-build FRONT_SRC=/path/para/seu/frontend
```

O `make web-build` chama `scripts/copy_frontend_dist.sh` que faz `npm install` (se necessário), `npm run build` e copia `dist` para `frontend/dist`.

Depois rode o backend normalmente (`go run cmd/webapp/main.go`) e abra `http://localhost:8081` para ver os arquivos estáticos servidos.

## Project Structure

```
├── cmd/webapp      # Application entrypoint
├── internal/
│   ├── domain     # Business entities
│   ├── ports      # Interfaces
│   ├── adapters   # External implementations
│   └── app        # Application core
└── pkg/           # Reusable packages
```

![CI](https://github.com/pamyydev/eco-link/actions/workflows/ci.yml/badge.svg)