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

```
├── cmd/webapp      # Application entrypoint
├── internal/
│   ├── domain     # Business entities
│   ├── ports      # Interfaces
│   ├── adapters   # External implementations
│   └── app        # Application core
└── pkg/           # Reusable packages
```

## Features

- Website sustainability analysis
- Carbon footprint calculation
- Green hosting verification

## Development

```bash
# Run the application
go run cmd/webapp/main.go
```

![CI](https://github.com/pamyydev/eco-link/actions/workflows/ci.yml/badge.svg)