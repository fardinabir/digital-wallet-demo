# Transaction Service Database Schema Documentation

This document provides information about the database schema, migration system, and seed data for the Transaction Service.

## Overview

The transaction service uses PostgreSQL as its database with a schema designed for financial transaction records. The system employs a three-tier migration approach:

1. **GORM Auto-Migration**: Automatic schema generation from Go models
2. **DDL Migrations**: Explicit schema definitions with constraints and indexes
3. **DML Migrations**: Data seeding and sample transactions

## Database Schema

### Transactions Table

Stores all transaction records with complete audit trail.

```sql
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    subject_wallet_id VARCHAR(255) NOT NULL,
    object_wallet_id VARCHAR(255),
    transaction_type VARCHAR(50) NOT NULL CHECK (transaction_type IN ('deposit', 'withdraw', 'transfer')),
    operation_type VARCHAR(50) NOT NULL CHECK (operation_type IN ('debit', 'credit')),
    amount BIGINT NOT NULL CHECK (amount > 0),
    status VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'completed', 'failed', 'cancelled')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
```

**Fields:**
- `id`: Primary key (auto-increment)
- `subject_wallet_id`: Wallet initiating the transaction
- `object_wallet_id`: Target wallet (provider wallet ID for deposits/withdrawals)
- `transaction_type`: Type of transaction (`deposit`, `withdraw`, `transfer`)
- `operation_type`: Operation type (`debit` or `credit`)
- `amount`: Transaction amount in cents
- `status`: Transaction status (`pending`, `completed`, `failed`, `cancelled`)
- `created_at`: Transaction creation timestamp
- `updated_at`: Last modification timestamp

**Constraints:**
- Check constraint ensuring amount is positive
- Check constraint for valid transaction types
- Check constraint for valid operation types
- Check constraint for valid status values

### Indexes

Optimized indexes for common query patterns:

- `idx_transactions_subject_wallet_id`: Index on subject wallet ID
- `idx_transactions_object_wallet_id`: Index on object wallet ID
- `idx_transactions_status`: Index on status
- `idx_transactions_created_at`: Index on creation time

### Triggers

Automatic timestamp management:
- `update_transactions_updated_at`: Updates `updated_at` on transaction modifications

## Migration System

### Migration Process

The system uses a three-phase migration approach:

1. **GORM Auto-Migration** (`internal/db/migration.go`)
   - Automatically creates tables from Go struct definitions
   - Handles basic schema changes
   - Ensures compatibility between code and database

2. **DDL Migrations** (`migrations/ddl/`)
   - Explicit schema definitions
   - Constraints, indexes, and triggers
   - Database-specific optimizations

3. **DML Migrations** (`migrations/dml/`)
   - Initial data seeding
   - Sample transaction data

### Migration Files

#### DDL Migration
- `migrations/ddl/001_create_transaction_schema.sql`: Complete schema definition

#### DML Migration
- `migrations/dml/001_insert_sample_transactions.sql`: Sample transaction data

### Running Migrations

```bash
# Run all migrations
make migrate

# Or run manually
go run main.go migrate
```

## Sample Transactions

Demonstration transaction history with double-entry bookkeeping:

- Completed deposit: Provider to user-001, 5,000 cents
- Completed withdrawal: user-002 to Provider, 1,000 cents
- Completed transfer: user-001 → user-003, 2,000 cents
- Pending deposit: Provider to user-002, 3,000 cents

## Database Configuration

### Development Environment

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

## Data Integrity

### Financial Constraints

1. **Amount Validation**: Transaction amounts must be positive
2. **Status Validation**: Only valid status values are allowed
3. **Type Validation**: Only valid transaction and operation types are allowed

### Transaction Consistency

1. **Double-Entry Bookkeeping**: Each transaction creates two records (debit and credit)
2. **Audit Trail**: Complete transaction history is maintained
3. **Status Tracking**: Transaction status progression is tracked
