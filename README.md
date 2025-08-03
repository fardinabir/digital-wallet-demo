# Digital Wallet Demo

A robust digital wallet system built with Go, Echo framework, and PostgreSQL. This project demonstrates modern backend development practices including clean architecture, comprehensive testing, and proper database design.

## 🚀 Features

- **Wallet Management**: Create and manage user wallets
- **Deposit Operations**: Add funds to wallets with provider integration
- **Withdrawal Operations**: Secure fund withdrawals with balance validation
- **Transfer Operations**: Peer-to-peer transfers between users
- **Transaction History**: Complete audit trail of all operations
- **RESTful API**: Well-documented REST endpoints
- **Comprehensive Testing**: 51+ test cases with 100% coverage
- **Database Migrations**: Automated schema and seed data management

## 🏗️ Architecture

```
├── cmd/                    # CLI commands and application entry points
├── internal/
│   ├── controller/         # HTTP handlers and routing
│   ├── service/           # Business logic layer
│   ├── repository/        # Data access layer
│   ├── model/            # Domain models and entities
│   ├── db/               # Database connection and migrations
│   └── errors/           # Custom error definitions
├── migrations/
│   ├── ddl/              # Database schema definitions
│   └── dml/              # Seed data and sample records
├── docs/                 # API documentation (Swagger)
└── config/               # Configuration files
```

## 🛠️ Technology Stack

- **Language**: Go 1.21+
- **Web Framework**: Echo v4
- **Database**: PostgreSQL 15+
- **ORM**: GORM v2
- **Testing**: Go testing package + Testify
- **Documentation**: Swagger/OpenAPI 3.0
- **Containerization**: Docker & Docker Compose
- **Build Tool**: Make

## 📋 Prerequisites

- Go 1.21 or higher
- PostgreSQL 15+ (or Docker)
- Make utility
- Git

## 🚀 Quick Start

### 1. Clone the Repository
```bash
git clone <repository-url>
cd digital-wallet-demo
```

### 2. Start Database (Docker)
```bash
# Start PostgreSQL container
docker-compose up -d

# Verify container is running
docker-compose ps
```

### 3. Run Database Migrations
```bash
# Apply all migrations (schema + seed data)
make migrate
```

### 4. Run Tests
```bash
# Execute all backend tests
make test-backend
```

### 5. Start the Application
```bash
# Start the server
make run

# Server will start on http://localhost:8080
```

### 6. Access API Documentation
Open your browser and navigate to:
- **Swagger UI**: http://localhost:8080/swagger/
- **API Docs**: http://localhost:8080/docs/

## 📚 Documentation

### Database
- **[Database Schema](DATABASE_SCHEMA.md)**: Complete schema documentation
- **[Database Setup](DATABASE_SETUP.md)**: Setup and migration guide

### API Endpoints

#### Wallet Operations
```
POST   /api/v1/wallets          # Create wallet
GET    /api/v1/wallets/:userId  # Get wallet by user ID
POST   /api/v1/wallets/deposit  # Deposit funds
POST   /api/v1/wallets/withdraw # Withdraw funds
POST   /api/v1/wallets/transfer # Transfer between wallets
```

#### Example Requests

**Create Wallet:**
```bash
curl -X POST http://localhost:8080/api/v1/wallets \
  -H "Content-Type: application/json" \
  -d '{"user_id": "user-123", "acnt_type": "user"}'
```

**Deposit Funds:**
```bash
curl -X POST http://localhost:8080/api/v1/wallets/deposit \
  -H "Content-Type: application/json" \
  -d '{"user_id": "user-123", "amount": 10000}'
```

**Transfer Funds:**
```bash
curl -X POST http://localhost:8080/api/v1/wallets/transfer \
  -H "Content-Type: application/json" \
  -d '{"from_user_id": "user-123", "to_user_id": "user-456", "amount": 5000}'
```

## 🧪 Testing

### Run All Tests
```bash
make test-backend
```

### Test Coverage
The project includes comprehensive test coverage:
- **Controller Tests**: HTTP endpoint testing
- **Service Tests**: Business logic validation
- **Repository Tests**: Data access testing
- **Integration Tests**: End-to-end scenarios

### Test Database
Tests use a separate database (`wallet_test`) to avoid conflicts:
```bash
# Run migrations for test database
make migrate-test
```

## 🔧 Development

### Available Make Commands
```bash
# Application
make run              # Start the application
make build            # Build the application
make clean            # Clean build artifacts

# Database
make migrate          # Run database migrations
make migrate-test     # Run test database migrations

# Docker
make docker-up        # Start PostgreSQL container
make docker-down      # Stop PostgreSQL container
make docker-clean     # Clean Docker volumes

# Testing
make test             # Run all tests
make test-backend     # Run backend tests only
make test-coverage    # Run tests with coverage

# Documentation
make docs             # Generate API documentation
make swagger          # Update Swagger specs
```

### Configuration

**Development Config** (`config.yaml`):
```yaml
server:
  port: 8080
  host: localhost

postgreSQL:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  dbname: wallet
  sslmode: disable
```

**Test Config** (`config.test.yaml`):
```yaml
postgreSQL:
  dbname: wallet_test
  # ... other settings same as development
```

## 🏦 Business Logic

### Wallet Types
- **User Wallets**: Individual user accounts
- **Provider Wallets**: System accounts for deposits/withdrawals

### Transaction Types
- **Deposit**: Add funds from external source
- **Withdrawal**: Remove funds to external destination
- **Transfer**: Move funds between user wallets

### Business Rules
- All amounts are stored in cents (integer)
- Negative balances are not allowed
- Transfers require sufficient balance
- All transactions are logged for audit
- Provider wallets have unlimited balance

### Default Provider Wallets
- `deposit-provider-master`: Source for deposit operations
- `withdraw-provider-master`: Destination for withdrawal operations

## 🔒 Security Considerations

### Current Implementation
- Input validation on all endpoints
- SQL injection prevention via GORM
- Transaction atomicity for transfers
- Balance validation before operations

### Production Recommendations
- Add authentication/authorization
- Implement rate limiting
- Add request/response encryption
- Enable database SSL
- Implement audit logging
- Add monitoring and alerting

## 📊 Database Schema

### Tables
- **wallets**: User and provider wallet information
- **transactions**: Complete transaction history

### Key Features
- Unique constraints on user wallets
- Automatic timestamp management
- Foreign key relationships
- Optimized indexes for queries

For detailed schema information, see [DATABASE_SCHEMA.md](DATABASE_SCHEMA.md).

## 🚀 Deployment

### Docker Deployment
```bash
# Build application image
docker build -t digital-wallet .

# Run with docker-compose
docker-compose up -d
```

### Environment Variables
```bash
# Database configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=wallet

# Server configuration
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

### Code Style
- Follow Go conventions
- Use `gofmt` for formatting
- Add comments for public functions
- Write comprehensive tests

## 📝 API Documentation

Complete API documentation is available via Swagger UI at `/swagger/` when the application is running.

### Response Format
All API responses follow a consistent format:

**Success Response:**
```json
{
  "status": "success",
  "data": {
    // Response data
  }
}
```

**Error Response:**
```json
{
  "status": "error",
  "message": "Error description",
  "code": "ERROR_CODE"
}
```

## 🐛 Troubleshooting

### Common Issues

1. **Database Connection Failed**
   - Ensure PostgreSQL is running
   - Check connection parameters
   - Verify database exists

2. **Migration Errors**
   - Check database permissions
   - Verify migration file syntax
   - Ensure proper file ordering

3. **Test Failures**
   - Run `make migrate-test` first
   - Check test database configuration
   - Verify test data setup

For detailed troubleshooting, see [DATABASE_SETUP.md](DATABASE_SETUP.md).

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- Echo framework for excellent HTTP routing
- GORM for powerful ORM capabilities
- PostgreSQL for robust database features
- Go community for excellent tooling

---

**Project Status**: ✅ Ready for submission

**Last Updated**: December 2024

**Test Coverage**: 51 tests passing

**Database**: PostgreSQL with complete schema and seed data