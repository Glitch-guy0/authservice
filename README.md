# Auth Service

A secure, production-ready authentication service built with Go, Gin, and modern best practices. This service provides JWT-based authentication with role-based access control, rate limiting, and comprehensive monitoring.

## ✨ Features

- 🔒 JWT-based authentication with refresh tokens
- 📝 Structured logging with logrus
- ⚙️ Configuration management with viper
- 🛑 Graceful shutdown handling
- 🩺 Health check endpoints
- ✅ Request validation
- 🔄 API versioning
- 📊 Prometheus metrics
- 🚦 Rate limiting
- 🔍 Request/Response logging
- 🔄 CORS support
- 🛡️ Security headers
- 📈 Performance monitoring

## 🚀 Prerequisites

- Go 1.21 or later
- Git
- Make (optional but recommended)
- Docker & Docker Compose (for containerized deployment)

## 🛠️ Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/Glitch-guy0/authService.git
   cd authService
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

## ⚙️ Configuration

1. Copy the example configuration file:
   ```bash
   cp configs/config.example.yaml configs/config.yaml
   ```

2. Update the configuration in `configs/config.yaml` with your settings.

3. (Optional) Set environment variables to override configuration:
   ```bash
   export APP_ENV=development
   export SERVER_PORT=8080
   ```

## 🏃‍♂️ Running the Application

### Development Mode

```bash
make run
```

This will start the server with live reload using `air`.

### Production Mode

```bash
make build && ./authService
```

### Using Docker

```bash
docker-compose up --build
```

## 📚 API Documentation

Once the server is running, you can access:

- **API Documentation**: `http://localhost:8080/docs` (Swagger UI)
- **Health Check**: `http://localhost:8080/health`
- **Metrics**: `http://localhost:8080/metrics`

## 🧪 Running Tests

```bash
# Run unit tests
make test

# Run integration tests
make test-integration

# Run benchmarks
make benchmark

# Check test coverage
make coverage
```

## 🧹 Code Quality

```bash
# Lint the code
make lint

# Format the code
make fmt

# Check for security vulnerabilities
make security
```

## 🤝 Contributing

We welcome contributions! Please read our [Contributing Guidelines](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📝 Environment Variables

For detailed information about environment variables, see [ENV.md](docs/ENV.md).

## 📖 API Reference

For detailed API documentation, see [API.md](docs/API.md).
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
## Project Structure

```
.
├── cmd/              # Main application entry points
│   └── app/main.go   # Service entrypoint
├── configs/          # Configuration files
│   └── config.yaml   # Application configuration
├── modules/          # Feature modules
│   ├── api/          # API layer modules
│   ├── core/         # Core modules (health, config, logger)
│   ├── server/       # Server setup and middleware
│   └── version/      # Version information
├── pkg/              # Reusable packages
│   └── errors/       # Error handling utilities
├── test/             # Test utilities and fixtures
│   ├── benchmark/    # Performance tests
│   ├── helpers/      # Test helpers
│   └── integration/  # Integration tests
└── specs/            # Feature specifications
```
## License

[MIT](LICENSE)
## Updated Module Structure

The project has been restructured to follow a modular architecture:

- **modules/api/**: API layer modules and handlers
- **modules/core/**: Core shared modules (health, config, logger)
- **modules/server/**: Server setup, routing, and middleware
- **modules/version/**: Version management and provider

This structure aligns with our controller pattern where each module exposes its routes through a controller interface.
