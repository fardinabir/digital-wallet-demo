# Transactions Service

A microservice for managing financial transactions in the digital wallet system.

## Features

- Create new transactions
- Retrieve transaction history with filtering
- Get transaction details by ID
- Support for different transaction types (deposit, withdrawal, transfer)
- RESTful API with Swagger documentation
- PostgreSQL database integration
- Docker support

## API Endpoints

### Transactions

- `POST /api/v1/transactions` - Create a new transaction
- `GET /api/v1/transactions` - Get all transactions with optional filters
- `GET /api/v1/transactions/{id}` - Get transaction by ID

### Health Check

- `GET /health` - Health check endpoint

## Quick Start

### Prerequisites

- Go 1.21+
- PostgreSQL 15+
- Docker (optional)

### Local Development

1. Clone the repository and navigate to the transactions service:
   ```bash
   cd services/transactions
   ```

2. Install dependencies:
   ```bash
   make deps
   ```

3. Start PostgreSQL database:
   ```bash
   make docker-up
   ```

4. Run database migrations:
   ```bash
   make migrate
   ```

5. Start the API server:
   ```bash
   make run
   ```

6. Access the API:
   - API Server: http://localhost:8082
   - Swagger Documentation: http://localhost:8082/swagger/index.html

### Docker Deployment

1. Build and start all services:
   ```bash
   docker-compose up --build
   ```

2. Access the services:
   - API Server: http://localhost:8082
   - Swagger Documentation: http://localhost:8082/swagger/index.html

## Configuration

The service uses YAML configuration files:

- `config.yaml` - Local development configuration
- `config.docker.yaml` - Docker environment configuration

### Configuration Structure

```yaml
api_server:
  host: "localhost"
  port: 8082
  timeout: 30

swagger_server:
  host: "localhost"
  port: 8082
  timeout: 30

postgresql:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "postgres"
  database: "transactions_db"
  ssl_mode: "disable"
  timezone: "UTC"
```

## Transaction Model

```json
{
  "id": 1,
  "subject_wallet_id": "wallet-123",
  "object_wallet_id": "wallet-456",
  "transaction_type": "transfer",
  "operation_type": "debit",
  "amount": 10000,
  "status": "completed",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### Transaction Types

- `deposit` - Money added to a wallet
- `withdrawal` - Money removed from a wallet
- `transfer` - Money moved between wallets

### Operation Types

- `credit` - Increases wallet balance
- `debit` - Decreases wallet balance

### Transaction Status

- `pending` - Transaction is being processed
- `completed` - Transaction completed successfully
- `failed` - Transaction failed
- `cancelled` - Transaction was cancelled

## Development Commands

```bash
# Install dependencies
make deps

# Run the application
make run

# Run tests
make test

# Run tests with coverage
make test-coverage

# Format code
make fmt

# Run linter
make lint

# Generate Swagger documentation
make swagger

# Build the application
make build

# Clean build artifacts
make clean

# Database operations
make migrate      # Run migrations
make reset-db     # Reset database

# Docker operations
make docker-up    # Start containers
make docker-down  # Stop containers
make docker-logs  # View logs
```

## API Examples

### Create Transaction

```bash
curl -X POST http://localhost:8082/api/v1/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "subject_wallet_id": "wallet-123",
    "object_wallet_id": "wallet-456",
    "transaction_type": "transfer",
    "operation_type": "debit",
    "amount": 10000
  }'
```

### Get Transactions

```bash
# Get all transactions
curl http://localhost:8082/api/v1/transactions

# Get transactions with filters
curl "http://localhost:8082/api/v1/transactions?subject_wallet_id=wallet-123&status=completed"
```

### Get Transaction by ID

```bash
curl http://localhost:8082/api/v1/transactions/1
```

## Architecture

The service follows a clean architecture pattern:

```
├── cmd/                 # CLI commands
├── docs/                # Swagger documentation
├── internal/
│   ├── controller/      # HTTP handlers
│   ├── db/             # Database connection
│   ├── errors/         # Error codes
│   ├── model/          # Data models
│   ├── repository/     # Data access layer
│   ├── server/         # Server setup
│   ├── service/        # Business logic
│   └── utils/          # Utilities
├── config.yaml         # Configuration
├── docker-compose.yml  # Docker setup
├── Dockerfile         # Container image
├── go.mod             # Go modules
├── main.go            # Application entry point
└── Makefile           # Build commands
```

## Contributing

1. Follow the existing code style and patterns
2. Add tests for new functionality
3. Update documentation as needed
4. Run linter and tests before submitting

## License

This project is part of the digital wallet demo system.