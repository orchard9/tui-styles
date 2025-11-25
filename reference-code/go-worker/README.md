# /var/folders/70/y0td7vjs3hz520207tgz_y0m0000gn/T/bootstrap-email-worker-3343678467/email-worker

/var/folders/70/y0td7vjs3hz520207tgz_y0m0000gn/T/bootstrap-email-worker-3343678467/email-worker service built with go-core-http-toolkit.

## Features

- **Service Type**: http
- **Databases**: postgres, redis
- **Authentication**: 
- **Additional**: metrics, health, job-processing

## Quick Start

### Prerequisites

- Go 1.23+
- Docker & Docker Compose
- Make

### Setup

1. Clone the repository:
```bash
git clone <repository-url>
cd /var/folders/70/y0td7vjs3hz520207tgz_y0m0000gn/T/bootstrap-email-worker-3343678467/email-worker
```

2. Copy environment configuration:
```bash
cp .env.example .env
# Edit .env with your configuration
```

3. Install dependencies:
```bash
make setup
```

4. Start development environment:
```bash
make dev-up
```

5. Run database migrations:
```bash
make migrate-up
```

6. Start development server:
```bash
make dev
```

The server will start on http://localhost:20323

## Development

### Available Commands

```bash
make help          # Show all available commands
make dev           # Run with hot reload
make test          # Run tests
make lint          # Run linters
make build         # Build production binary
make complexity    # Check code complexity
make file-length   # Check file/function length
make deadcode      # Detect dead code
make circular      # Check for circular dependencies
```

### Continuous Monitoring

This project includes watch-now for continuous code quality monitoring:

```bash
# Install watch-now (first time only)
go install github.com/orchard9/watch-now@latest

# Start continuous monitoring
watch-now

# Or run checks once
watch-now --once
```

watch-now monitors your code in real-time and runs:
- Formatting checks
- Linting
- Code complexity analysis
- Dead code detection
- Circular dependency checks
- Tests and test coverage
- Build verification

Configuration is in `.watch-now.yaml`. Customize as needed.

### Project Structure

```
.
├── cmd/
│   ├── server/          # Main server application
│   └── migrate/         # Database migration tool
├── internal/
│   ├── config/          # Configuration management
│   ├── handlers/        # HTTP handlers
│   ├── repository/      # Database access layer
│   └── service/         # Business logic
├── api/
├── migrations/          # Database migrations
├── docker-compose.yml   # Development environment
└── Makefile            # Development commands
```

## API Documentation

API documentation coming soon.

## Testing

Run unit tests:
```bash
make test
```

Run integration tests:
```bash
go test -tags=integration ./...
```

## Deployment

### Build Docker Image

```bash
make docker-build
```

### Environment Variables

See `.env.example` for all available configuration options.

## License

[Your License Here]
