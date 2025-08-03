# Database Setup Guide

This guide provides step-by-step instructions for setting up the PostgreSQL database for the Digital Wallet Demo project.

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
   make test-backend
   ```

4. **Start the application:**
   ```bash
   make run
   ```


## Database Schema Overview

The system creates two main tables:

### Wallets Table
- Stores user and provider wallet information
- Tracks balance, status, and account type
- Enforces business constraints at database level

### Transactions Table
- Records all financial transactions
- Maintains complete audit trail
- Supports deposits, withdrawals, and transfers

## Migration System

The project uses a three-tier migration approach:

1. **GORM Auto-Migration**: Creates basic table structure from Go models
2. **DDL Migrations**: Adds constraints, indexes, and triggers
3. **DML Migrations**: Inserts seed data and sample records

### Migration Files

```
migrations/
├── ddl/
│   └── 001_create_wallet_schema.sql    # Schema definition
└── dml/
    └── 001_insert_provider_wallets.sql # Seed data
```

## Seed Data

After migration, the database contains:

### Provider Wallets (Required for System Operation)
- `deposit-provider-master`: Source for deposit transactions
- `withdraw-provider-master`: Destination for withdrawal transactions

## Configuration Files

### Development Configuration
```yaml
# config.yaml
postgreSQL:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  dbname: wallet
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
  dbname: wallet_test
  sslmode: disable
```

## Available Make Commands

```bash
# Database operations
make migrate          # Run all database migrations
make migrate-test     # Run migrations for test database

# Docker operations
make docker-up        # Start PostgreSQL container
make docker-down      # Stop PostgreSQL container
make docker-clean     # Stop and remove volumes

# Testing
make test-backend     # Run all backend tests
make test             # Run all tests

# Application
make run              # Start the application
make build            # Build the application
```
