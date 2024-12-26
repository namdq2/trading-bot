# Market Data Service

A microservice for collecting, processing, and distributing real-time market data from multiple cryptocurrency exchanges. Built using Domain-Driven Design and Clean Architecture principles.

## Table of Contents
- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Configuration](#configuration)
- [Building and Running](#building-and-running)
- [Deployment](#deployment)
- [API Documentation](#api-documentation)
- [Testing](#testing)

## Architecture

This service follows Clean Architecture and DDD principles:

```
┌───────────────────┐
│     Interface     │ HTTP/gRPC handlers, Event consumers
├───────────────────┤
│   Application     │ Use cases, DTOs, Ports
├───────────────────┤
│     Domain        │ Entities, Value Objects, Domain Services
├───────────────────┤
│ Infrastructure    │ Repositories, External Services
└───────────────────┘
```

### Key Components:
- Domain Layer: Core business logic and rules
- Application Layer: Orchestration and use cases
- Infrastructure Layer: External systems integration
- Interface Layer: API endpoints and event handlers

## Prerequisites

- Go 1.21 or later
- Docker and Docker Compose
- Redis
- PostgreSQL/TimescaleDB
- Kafka

## Development Setup

1. Clone the repository:
```bash
git clone https://github.com/yourusername/marketdata.git
cd marketdata
```

2. Install dependencies:
```bash
go mod download
```

3. Set up local infrastructure:
```bash
docker-compose up -d redis kafka timescaledb
```

4. Create configuration file:
```bash
cp config/config.example.yaml config/config.yaml
```

5. Run the service:
```bash
go run cmd/main.go
```

## Project Structure

```
marketdata/
├── cmd/                    # Application entry points
├── internal/              # Private application code
│   ├── domain/            # Domain layer
│   │   ├── entity/        # Domain entities
│   │   ├── repository/    # Repository interfaces
│   │   ├── service/       # Domain services
│   │   └── valueobject/   # Value objects
│   ├── infrastructure/    # Infrastructure layer
│   │   ├── persistence/   # Database implementations
│   │   ├── exchange/      # Exchange clients
│   │   └── messaging/     # Message queue implementations
│   ├── application/       # Application layer
│   │   ├── dto/          # Data Transfer Objects
│   │   ├── service/      # Application services
│   │   └── port/         # Ports (interfaces)
│   └── interfaces/        # Interface layer
│       ├── api/          # HTTP/gRPC handlers
│       └── event/        # Event handlers
├── pkg/                   # Public libraries
└── config/               # Configuration files
```

## Configuration

Configuration is managed through `config.yaml`:

```yaml
server:
  http_port: 8080
  grpc_port: 9090

database:
  host: localhost
  port: 5432
  user: postgres
  password: password
  dbname: marketdata

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0
  ttl: 1h

kafka:
  brokers:
    - localhost:9092
  topic: orderbook_updates
  group_id: marketdata_service

exchange:
  symbols:
    - BTC-USDT
    - ETH-USDT
  binance:
    api_key: your_api_key
    api_secret: your_api_secret
  okx:
    api_key: your_api_key
    api_secret: your_api_secret
```

## Building and Running

### Local Development
```bash
# Run with hot reload
go install github.com/cosmtrek/air@latest
air

# Run tests
go test ./...

# Run linter
golangci-lint run
```

### Docker Build
```bash
# Build image
docker build -t marketdata:latest .

# Run container
docker run -p 8080:8080 -p 9090:9090 marketdata:latest
```

## Deployment

### Kubernetes Deployment

1. Create Kubernetes manifests:

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: marketdata
spec:
  replicas: 3
  selector:
    matchLabels:
      app: marketdata
  template:
    metadata:
      labels:
        app: marketdata
    spec:
      containers:
      - name: marketdata
        image: marketdata:latest
        ports:
        - containerPort: 8080
        - containerPort: 9090
        envFrom:
        - configMapRef:
            name: marketdata-config
```

2. Deploy to Kubernetes:
```bash
kubectl apply -f k8s/
```

### Monitoring

The service exposes metrics for Prometheus at `/metrics` and includes:
- Request latencies
- Error rates
- Exchange connection status
- Order book update rates

### Logging

Logs are structured using zap and include:
- Request/Response logging
- Error tracking
- Performance metrics
- Exchange interactions

## API Documentation

### HTTP Endpoints

```
GET /api/v1/orderbook
    Query Parameters:
    - exchange: Exchange ID
    - symbol: Trading pair symbol

GET /api/v1/trades
    Query Parameters:
    - exchange: Exchange ID
    - symbol: Trading pair symbol
    - limit: Number of trades (default: 100)
```

### gRPC Services

See `api/proto/marketdata.proto` for service definitions.

## Testing

```bash
# Run unit tests
go test ./...

# Run integration tests
go test -tags=integration ./...

# Generate test coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.