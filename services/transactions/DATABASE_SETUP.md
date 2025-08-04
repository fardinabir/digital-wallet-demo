# Database Setup Guide

This guide provides step-by-step instructions for setting up the PostgreSQL database for the Transaction Service.

## Prerequisites

- PostgreSQL 15+ installed locally OR Docker with Docker Compose
- Go 1.24+ installed
- Make utility installed

## Quick Setup (Recommended)

### Using Docker Compose

1. **Start PostgreSQL container:**
   ```bash
   docker-compose up -d
   ```

2. **Run database migrations:**
   ```bash
   make migrate
   ```

3. **Verify setup:**
   ```bash
   make test
   ```

4. **Start the application:**
   ```bash
   make serve
   ```

## Database Schema Overview

The system creates one main table:

### Transactions Table
- Records all financial transactions
- Maintains complete audit trail
- Supports deposits, withdrawals, and transfers
- Uses double-entry bookkeeping pattern

## Migration System

The project uses a three-tier migration approach:

1. **GORM Auto-Migration**: Creates basic table structure from Go models
2. **DDL Migrations**: Adds constraints, indexes, and triggers
3. **DML Migrations**: Inserts sample transaction data

### Migration Files

```
migrations/
├── ddl/
│   └── 001_create_transaction_schema.sql    # Schema definition
└── dml/
    └── 001_insert_sample_transactions.sql   # Sample data
```

## Sample Data

After migration, the database contains sample transactions demonstrating:
- Deposit transactions (provider to user)
- Withdrawal transactions (user to provider)
- Transfer transactions (user to user)
- Various transaction statuses (completed, pending)

## Configuration Files

### Development Configuration
```yaml
# config.yaml
postgreSQL:
  host: localhost
  port: 5433
  user: postgres
  password: postgres
  dbname: transaction
  sslmode: disable
```

### Test Configuration
```yaml
# config.test.yaml
postgreSQL:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  dbname: transaction_test
  sslmode: disable
```

## Available Make Commands

```bash
# Database operations
make migrate          # Run all database migrations
make migrate-test     # Run migrations for test database
make reset-db         # Reset development database
make reset-test-db    # Reset test database

# Docker operations
make start            # Start PostgreSQL container
make stop             # Stop PostgreSQL container
make clear            # Stop and remove volumes

# Testing
make test             # Run all tests
make test-ci          # Run tests with coverage

# Application
make serve            # Start the application
make lint             # Run linter
make fmt              # Format code
```

## Troubleshooting

### Common Issues

1. **Database connection failed**: Ensure PostgreSQL is running and configuration is correct
2. **Migration errors**: Check if database exists and user has proper permissions
3. **Test failures**: Ensure test database is properly reset before running tests

### Database Reset

If you encounter issues, reset the database:

```bash
# For development
make reset-db

# For testing
make reset-test-db
```
