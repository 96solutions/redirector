# Redirective

A high-performance HTTP redirection service built with Go, utilizing Domain-Driven Design with clean architecture.

## Run

### Using Docker Compose

```bash
# Start the service with all dependencies
docker-compose up app

# For development with hot reload
docker-compose up dev
```

### Building from Source

```bash
# Build the application
make build

# Or clean and build
make all

# Install the application to bin directory
make install
```

## Tests

Run all tests:
```bash
go test ./...
```

Run a specific test:
```bash
go test ./path/to/package -run TestName
```

Generate mocks for testing:
```bash
make mocks
```

## Development

### Requirements

- Go 1.22+
- Docker and Docker Compose
- PostgreSQL
- Redis
- ClickHouse
- OpenSearch (for logging)

### Environment Configuration

The application uses environment variables for configuration. Set them directly or create a `.env` file in the project root.

Key configuration options:
- Database settings: `DB_HOST`, `DB_PORT`, `DB_USERNAME`, `DB_PASSWORD`
- HTTP server: `HTTP_SERVER_PORT` (default: 8080)
- Redis cache: `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASS`
- Logging: `LOG_LEVEL`, `LOG_IS_JSON`

Run linting:
```bash
golangci-lint --exclude-use-default=false --out-format tab run ./...
```

## Metrics

The service exposes Prometheus metrics and includes:

- Grafana dashboards for visualization
- Prometheus for metrics collection
- Pre-configured dashboards in `docker/grafana/provisioning/dashboards/`

Access the dashboards:
- Grafana: http://localhost:3000 (default credentials: admin/admin)
- Prometheus: http://localhost:9090

## Deployment

### Docker

The service includes production-ready Dockerfiles:
- `Dockerfile`: Optimized multi-stage build for production
- `Dockerfile.builder`: For CI/CD pipeline builds

The application is designed to run in Kubernetes environments with:
- Health check endpoints
- Prometheus metrics
- Optimized resource usage

## TODO
- Tests
- ClickHandler impl
  - Save click to queue (kafka)
  - Save click to storage (clickhouse)
  - gRPC call
- Add health check
- Add homepage
- Extend redirector functionality
  - deeplinking
  - landing pages
  - transparent redirection
  - parallel tracking
  - use publisher ID + campaign ID (source id optional) instead of slug
- CI/CD
  - K8 + Helm
  - Github Actions
  - AWS / GCP
- Performance Improvements
  - EasyJSON
  + Redis cache
  - Clickhouse storage
  - load landing pages only if needed.