# Creator API

> Go service for creator account management, avatar upload, and analytics

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Quick Start](#quick-start)
- [Project Structure](#project-structure)
- [Development Commands](#development-commands)
- [API Endpoints](#api-endpoints)
- [Configuration](#configuration)
- [Docker Deployment](#docker-deployment)
- [HTTP Framework: Chi](#http-framework-chi)
- [Code Quality & Linting Standards](#code-quality--linting-standards)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
- [Roadmap](#roadmap)
- [Resources](#resources)
- [License](#license)
- [Support](#support)

## Overview

The Creator API provides backend services for the Masquerade Creator Studio platform, enabling avatar creators to manage their accounts, upload avatars, and track analytics. Built with Go for high performance and scalability.

**Core Capabilities**:
- Creator signup and trial management
- Avatar upload and processing workflows
- Creator dashboard analytics
- Earnings and payout tracking

**Current Status**: Milestone 2 - Foundation with mocked endpoints (no database integration yet)

## Features

- ✅ RESTful API with versioned endpoints (`/api/v1`)
- ✅ Structured logging with zerolog (JSON/console modes)
- ✅ Comprehensive middleware (CORS, request ID, recovery, timing)
- ✅ OpenAPI 3.0 specification with Swagger UI
- ✅ >80% test coverage with table-driven tests
- ✅ Docker support with multi-stage builds
- ✅ Pre-commit hooks for code quality
- ✅ Encrypted environment management with envault
- ⏳ Database integration (coming in Milestone 3)
- ⏳ JWT authentication (coming in Milestone 4)
- ⏳ S3 avatar storage (coming in Milestone 5)

## Quick Start

**Prerequisites**:
- Go 1.23 or later ([install](https://go.dev/doc/install))
- Make (optional, for build automation)
- envault (for environment management - [install](https://github.com/envault/envault))

**Setup**:
```bash
# 1. Navigate to the project
cd services/creator-api

# 2. Download dependencies
go mod download

# 3. Verify setup
go mod verify

# 4. Install development tools (optional but recommended)
make install-tools
# Installs: golangci-lint, goimports, air (hot reload)

# 5. Initialize envault (first-time only)
envault init
envault add-key ~/.ssh/id_rsa.pub

# 6. Run the server
envault dev && make run
# Or with auto-reload:
envault dev && make dev
```

**Verify it's working**:
```bash
# Health check
curl http://localhost:34075/health
# Expected response:
# {"status":"healthy","version":"0.1.0","timestamp":"2025-11-20T..."}

# Test ping endpoint
curl http://localhost:34075/api/v1/ping
# Expected response:
# {"message":"pong","timestamp":"2025-11-20T..."}
```

**Explore the API**:
- Interactive docs: http://localhost:34075/docs
- OpenAPI spec: http://localhost:34075/docs/openapi.yaml

That's it! You're ready to start developing.

## Project Structure

```
services/creator-api/
├── cmd/server/        # Application entry point
├── internal/          # Private application code
│   ├── api/          # HTTP handlers, middleware, routes
│   ├── config/       # Configuration management
│   └── mock/         # Mock data generators
├── pkg/models/       # Shared types/structs
├── docs/             # API documentation (OpenAPI, Swagger UI)
├── .golangci.yml     # Linter configuration
├── .env.example      # Environment variable template
├── Dockerfile        # Multi-stage Docker build
├── Makefile          # Development commands
└── README.md         # This file
```

## Development Commands

### Prerequisites
Install development tools (if not already installed):
```bash
make install-tools
```

### Common Tasks

| Command | Description |
|---------|-------------|
| `make help` | Show all available commands (default) |
| `make build` | Build production binary to bin/creator-api |
| `make run` | Run the server |
| `make dev` | Run with auto-reload (hot reload) |
| `make test` | Run tests with coverage |
| `make lint` | Run code linters |
| `make fmt` | Format code (gofmt + goimports) |
| `make fmt-check` | Check formatting without modifying |
| `make deadcode` | Detect unused code |
| `make complexity` | Report code complexity (threshold: >10) |
| `make security` | Run security scans (gosec) |
| `make quality` | Run ALL quality checks |
| `make clean` | Remove build artifacts |
| `make all` | Run all quality checks + build |

### Development Workflow
```bash
# First time setup (if tools not installed)
make install-tools

# Daily development
envault dev && make dev    # Start server with auto-reload
make test                  # Run tests (in another terminal)
make quality              # Check code quality (all checks)
make all                  # Format + quality + build (before committing)
```

## API Endpoints

### Health & Monitoring

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Service health check |
| GET | `/api/v1/ping` | Simple test endpoint |

### Creator Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/creators/signup` | Register new creator account |
| GET | `/api/v1/creators/:id/dashboard` | Get creator dashboard metrics |
| GET | `/api/v1/creators/:id/earnings` | Get earnings breakdown |

### Avatar Operations

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/avatars/upload` | Initiate avatar upload |
| GET | `/api/v1/avatars/:id` | Retrieve avatar metadata |

### Documentation

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/docs` | Swagger UI (interactive API docs) |
| GET | `/docs/openapi.yaml` | OpenAPI 3.0 specification |

**Full Interactive Documentation**: http://localhost:34075/docs (when server is running)

**Example Requests**:
```bash
# Creator signup
curl -X POST http://localhost:34075/api/v1/creators/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "creator@example.com",
    "username": "avatar_master",
    "display_name": "Avatar Master"
  }'

# Get dashboard metrics
curl http://localhost:34075/api/v1/creators/123/dashboard

# Upload avatar
curl -X POST http://localhost:34075/api/v1/avatars/upload \
  -H "Content-Type: application/json" \
  -d '{
    "creator_id": "123",
    "avatar_name": "Cool Avatar",
    "file_size": 1048576
  }'
```

## Configuration

This service uses **envault** for encrypted environment management. All configuration is loaded from environment variables (no .env files at runtime).

**Environment Variables**:

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `34075` | HTTP server port |
| `HOST` | `localhost` | Server host |
| `ENV` | `development` | Environment: development, staging, production |
| `LOG_LEVEL` | `debug` | Log level: debug, info, warn, error |
| `LOG_FORMAT` | `console` | Log format: console (dev), json (prod) |
| `API_VERSION` | `v1` | API version |
| `API_PREFIX` | `/api` | API prefix for all routes |
| `CORS_ALLOWED_ORIGINS` | `http://localhost:3000,http://localhost:3001` | Allowed CORS origins |
| `CORS_ALLOWED_METHODS` | `GET,POST,PUT,DELETE,OPTIONS,PATCH` | Allowed HTTP methods |
| `CORS_ALLOWED_HEADERS` | `Content-Type,Authorization,X-Request-ID` | Allowed headers |

**Configuration Reference**:
- See `.env.example` for all available configuration options
- `.env.example` is for documentation only (NOT loaded at runtime)
- Actual values are managed via envault (encrypted in repository)

**Envault Setup** (first-time):
```bash
# Initialize envault
envault init

# Add your SSH public key
envault add-key ~/.ssh/id_rsa.pub

# Load development environment
envault dev
```

## Docker Deployment

### Build Image

```bash
# Using Make
make docker-build

# Or manually
docker build -t masquerade-creator-api:latest .
```

The Dockerfile uses multi-stage builds for optimal image size:
- Build stage: Compiles Go binary
- Runtime stage: Minimal Alpine Linux with binary only

### Run Container

```bash
# Using Make
make docker-run

# Or manually
docker run --rm -p 34075:34075 \
  --env-file .env \
  masquerade-creator-api:latest
```

### Docker Compose

```bash
# Start all services (from workspace root)
docker-compose up -d

# View logs
docker-compose logs -f creator-api

# Stop services
docker-compose down
```

## HTTP Framework: Chi

The Creator API uses **Chi v5** (`github.com/go-chi/chi/v5`) as its HTTP routing framework.

### Why Chi?

Chi was selected for the following reasons:

1. **Idiomatic Go**: Built on `net/http` standard library patterns, minimal abstraction
2. **Lightweight**: <1000 lines of code, easy to understand internals
3. **Context-based**: Middleware uses request context properly (no custom context types)
4. **Excellent Routing**: Clean path parameter support (`/:id`, `/*`)
5. **Composable**: Middleware chaining is explicit and follows standard patterns
6. **Active Maintenance**: Well-documented, stable, and actively maintained
7. **Performance**: Fast routing with minimal overhead

**Alternatives Rejected**:
- **Gin**: Faster in benchmarks but less idiomatic (custom context)
- **Echo**: Good framework but Chi is more minimalist
- **stdlib only**: Too low-level, would require reinventing routing

### Router Structure

The router is initialized in `internal/api/routes.go` with global middleware and versioned API routes:

```go
func NewRouter() *chi.Mux {
    r := chi.NewRouter()

    // Global middleware stack
    r.Use(middleware.RequestID)   // Inject request ID into context
    r.Use(middleware.RealIP)      // Set RemoteAddr to real IP
    r.Use(middleware.Logger)      // Log all requests
    r.Use(middleware.Recoverer)   // Recover from panics

    // API v1 routes (versioned for future compatibility)
    r.Route("/api/v1", func(r chi.Router) {
        r.Get("/ping", handlers.HandlePing)
        r.Post("/creators/signup", handlers.HandleCreatorSignup)
        // ... other routes
    })

    return r
}
```

### Middleware

Chi uses the standard Go middleware signature: `func(http.Handler) http.Handler`

**Built-in Middleware** (applied globally):
- `RequestID`: Injects unique request ID into context (accessible via `middleware.GetReqID(r.Context())`)
- `RealIP`: Sets `RemoteAddr` to the real client IP (handles X-Forwarded-For)
- `Logger`: Logs all requests with method, path, status, duration
- `Recoverer`: Recovers from panics and returns 500 status

### Routing Patterns

**Path Parameters**:
```go
r.Get("/avatars/{id}", handlers.HandleGetAvatar)
// Extract param: chi.URLParam(r, "id")
```

**Wildcards**:
```go
r.Get("/files/*", handlers.HandleFiles)
// Access via: chi.URLParam(r, "*")
```

**Route Groups**:
```go
r.Route("/creators/{id}", func(r chi.Router) {
    r.Get("/dashboard", handlers.HandleCreatorDashboard)
    r.Get("/earnings", handlers.HandleCreatorEarnings)
})
```

### Testing Routes

Routes are tested using `net/http/httptest`:

```go
router := NewRouter()
req := httptest.NewRequest("GET", "/api/v1/ping", nil)
w := httptest.NewRecorder()
router.ServeHTTP(w, req)

assert.Equal(t, http.StatusOK, w.Code)
```

See `internal/api/routes_test.go` for examples.

## Code Quality & Linting Standards

The Creator API enforces strict code quality standards using **golangci-lint**, a comprehensive linting tool that runs 16+ linters in parallel to catch bugs, enforce style, and ensure security.

### Running Linters

```bash
# Run all linters (recommended before committing)
make lint

# Run with verbose output (see which linters are running)
golangci-lint run -v ./...

# Auto-fix issues where possible
golangci-lint run --fix ./...

# Run comprehensive quality checks (formatting, linting, tests, security)
make quality
```

### Enabled Linters

**Formatting & Style**:
- `gofmt` - Enforces standard Go formatting
- `goimports` - Organizes imports properly
- `staticcheck` - Advanced static analysis (includes gosimple)
- `unconvert` - Removes unnecessary type conversions
- `whitespace` - Detects leading/trailing whitespace

**Bug Detection**:
- `govet` - Reports suspicious constructs (vet all passes enabled)
- `errcheck` - Ensures all errors are checked
- `ineffassign` - Detects ineffectual assignments
- `unused` - Finds unused constants, variables, functions
- `errorlint` - Catches common error wrapping mistakes

**Security**:
- `gosec` - Scans for security vulnerabilities (medium severity+)

**Performance**:
- `prealloc` - Identifies slice declarations that should pre-allocate

**Complexity**:
- `gocyclo` - Cyclomatic complexity (threshold: >15)
- `gocognit` - Cognitive complexity (threshold: >20)

**Code Quality**:
- `nakedret` - Flags naked returns in long functions (>5 lines)
- `misspell` - Catches spelling errors in comments/strings

### Configuration

Linting is configured in `.golangci.yml` with project-specific rules and exclusions:

- **Test files**: More lenient rules (complexity, duplicate code OK)
- **Mock/generated code**: Excluded from style checks
- **HTTP handlers**: Duplicate patterns allowed (common structure)
- **Timeout**: 5 minutes maximum
- **Go version**: 1.23+

### Common Issues & Fixes

**Unchecked errors**:
```go
// Bad
json.NewEncoder(w).Encode(response)

// Good
if err := json.NewEncoder(w).Encode(response); err != nil {
    log.Error().Err(err).Msg("failed to encode response")
}
```

**Error string formatting**:
```go
// Bad
errors.New("Failed to load config")

// Good
errors.New("failed to load config")
```

**Context parameter ordering**:
```go
// Bad
func DoSomething(id string, ctx context.Context)

// Good
func DoSomething(ctx context.Context, id string)
```

### Pre-commit Checklist

Before submitting code, ensure:
1. `make fmt` - Code is formatted
2. `make lint` - All linters pass (0 issues)
3. `make test` - All tests pass
4. `make quality` - Comprehensive checks pass

Or run `make all` to execute all checks and build.

## Troubleshooting

### Port already in use

**Problem**: Server fails to start with "address already in use" error

**Solution**:
```bash
# Find process using port 34075
lsof -i :34075

# Kill the process
kill -9 <PID>

# Or change the port in envault configuration
envault edit dev
# Update PORT=34076
```

### Module errors

**Problem**: "module not found" or "checksum mismatch" errors

**Solution**:
```bash
# Tidy dependencies
go mod tidy

# Verify dependencies
go mod verify

# Clear module cache if needed
go clean -modcache

# Re-download dependencies
go mod download
```

### Tests failing

**Problem**: Tests pass locally but fail in CI or vice versa

**Solution**:
```bash
# Clean test cache
go clean -testcache

# Run tests with verbose output
go test -v ./...

# Run with race detector
go test -v -race ./...

# Run specific test
go test -v -run TestHandleCreatorSignup ./internal/api/handlers/
```

### Linter issues

**Problem**: golangci-lint reports errors

**Solution**:
```bash
# Auto-fix formatting
make fmt

# Auto-fix linter issues (where possible)
golangci-lint run --fix

# Check specific linter
golangci-lint run --enable-only=errcheck ./...

# See verbose output
golangci-lint run -v ./...
```

### envault errors

**Problem**: "failed to decrypt" or "key not found" errors

**Solution**:
```bash
# Re-initialize envault
envault init

# Add your SSH key
envault add-key ~/.ssh/id_rsa.pub

# Verify key is added
envault list-keys

# Load environment manually
envault dev
```

### Hot reload not working

**Problem**: `make dev` doesn't reload on file changes

**Solution**:
```bash
# Reinstall air
go install github.com/cosmtrek/air@latest

# Check air is in PATH
which air

# Run air manually with verbose output
air -c .air.toml
```

### Docker build fails

**Problem**: Docker build errors or slow builds

**Solution**:
```bash
# Clear Docker cache
docker system prune -a

# Build with no cache
docker build --no-cache -t masquerade-creator-api:latest .

# Check Dockerfile syntax
docker build --check -t masquerade-creator-api:latest .
```

### High memory usage

**Problem**: Server consumes excessive memory

**Solution**:
```bash
# Profile memory usage
go test -memprofile=mem.out ./...
go tool pprof mem.out

# Check for goroutine leaks
curl http://localhost:34075/debug/pprof/goroutine

# Enable memory profiling
ENV=development make run
# Visit http://localhost:34075/debug/pprof/
```

## Contributing

We welcome contributions! Follow these guidelines to ensure smooth integration.

### Workflow

1. **Create a feature stream/branch**:
   ```bash
   # Using Perforce
   p4 switch -c feature/my-feature
   # Or use P4V to create stream
   ```

2. **Make your changes**:
   - Write clear, focused code
   - Follow existing patterns and conventions
   - Add tests for new functionality

3. **Run quality checks**:
   ```bash
   make all  # Format + lint + test + build
   ```

4. **Write tests**:
   - Maintain >80% coverage
   - Use table-driven tests for multiple scenarios
   - Include edge cases and error handling

5. **Submit with conventional commit messages**:
   - `feat:` New feature
   - `fix:` Bug fix
   - `docs:` Documentation
   - `test:` Tests
   - `refactor:` Code refactoring
   - `perf:` Performance improvements
   - `chore:` Maintenance tasks

6. **Merge to parent stream** or submit for review

### Code Review Checklist

Before submitting for review, ensure:

- [ ] Tests added for new functionality
- [ ] Test coverage >80% (run `make test` to verify)
- [ ] Linter passes with zero warnings (`make lint`)
- [ ] Documentation updated (README, comments, OpenAPI spec)
- [ ] No breaking API changes (or properly documented)
- [ ] Error handling follows established patterns
- [ ] Logging includes context (request_id, relevant data)
- [ ] No secrets committed (API keys, passwords, etc.)
- [ ] Code is formatted (`make fmt`)
- [ ] Pre-commit hooks pass

### Coding Standards

**Functions**:
- Keep functions under 50 lines where possible
- Single responsibility principle
- Clear, descriptive names

**Error Handling**:
- Always check errors
- Use error wrapping for context: `fmt.Errorf("failed to %s: %w", action, err)`
- Log errors with appropriate level (Error, Warn, Info)

**Testing**:
- Table-driven tests for multiple scenarios
- Use `testify/assert` for assertions
- Mock external dependencies
- Test edge cases and error paths

**Logging**:
- Use structured logging (zerolog)
- Include request ID in all logs
- Log at appropriate levels (Debug, Info, Warn, Error)

## Roadmap

### Milestone 2 (Current) ✅
- ✅ Mock API endpoints with realistic data
- ✅ Quality infrastructure (linting, testing, pre-commit hooks)
- ✅ Docker support with multi-stage builds
- ✅ OpenAPI documentation with Swagger UI
- ✅ Comprehensive README and documentation

### Milestone 3 (Next)
- PostgreSQL integration with migrations
- Real data storage and retrieval
- Database connection pooling
- Transaction management

### Milestone 4
- JWT authentication and authorization
- API key management
- Role-based access control (RBAC)
- Rate limiting

### Milestone 5
- S3 avatar storage integration
- File upload handling with multipart support
- CDN integration for avatar delivery
- Image processing pipeline

## Resources

- **API Documentation**: http://localhost:34075/docs (Swagger UI)
- **OpenAPI Spec**: [docs/openapi.yaml](docs/openapi.yaml)
- **Product Specs**: [../../product-specs/creator.md](../../product-specs/creator.md)
- **Roadmap**: [../../roadmap/milestone-2/](../../roadmap/milestone-2/)
- **Go Best Practices**: https://go.dev/doc/effective_go
- **Chi Router Docs**: https://github.com/go-chi/chi
- **zerolog Docs**: https://github.com/rs/zerolog

## Port Allocation

- **34075**: Creator API (this service)

See `CLAUDE.md` in the workspace root for complete Masquerade platform port allocation (34070-34079 range).

## Technology Stack

- **Language**: Go 1.23+
- **HTTP Framework**: Chi router (go-chi/chi/v5)
- **Logging**: zerolog (structured logging)
- **Configuration**: envault (encrypted environment management)
- **Testing**: Standard library + testify/assert
- **Linting**: golangci-lint (16+ linters)
- **Documentation**: OpenAPI 3.0 + Swagger UI

## Development Status

**Current Milestone**: Milestone 2 - Foundation with Mocked Endpoints
**Phase**: Phase 3 - Polish
**Status**: Complete (pending final review)

See `roadmap/milestone-2/` for detailed task breakdown and progress.

## License

MIT License - see [LICENSE](LICENSE) file for details

## Support

- **Issues**: File on GitHub or internal issue tracker
- **Questions**: #creator-api Slack channel
- **Email**: engineering@masquerade.ai

---

**Built with ❤️ by the Masquerade Engineering Team**
