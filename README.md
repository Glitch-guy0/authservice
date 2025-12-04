# Auth Service

A secure authentication service built with Go, Gin, and modern best practices.

## Features

- JWT-based authentication
- Structured logging with logrus
- Configuration management with viper
- Graceful shutdown
- Health check endpoints
- Request validation
- API versioning

## Prerequisites

- Go 1.21 or later
- Git
- Make (optional)

## Getting Started

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/Glitch-guy0/authService.git
   cd authService
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

### Configuration

Copy the example environment file and update the values:

```bash
cp configs/config.example.yaml configs/config.yaml
```

### Running the Application

```bash
# Development mode
make run

# Production mode
make run-prod
```

## Development

### Building

```bash
make build
```

### Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage
```

### Linting

```bash
make lint
```

## Project Structure

```
.
├── cmd/              # Main application entry points
├── configs/          # Configuration files
├── internal/         # Private application code
│   ├── api/          # HTTP handlers and middleware
│   ├── config/       # Configuration loading
│   ├── logger/       # Logging implementation
│   └── modules/      # Feature modules
├── pkg/              # Reusable packages
└── test/             # Test utilities and fixtures
```

## License

[MIT](LICENSE)
