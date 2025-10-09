# eco-link

Analyze website sustainability and carbon footprint.

## Structure

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