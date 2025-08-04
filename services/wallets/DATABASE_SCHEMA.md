# Digital Wallet Database Schema Documentation

This document provides comprehensive information about the database schema, migration system, and seed data for the Digital Wallet Demo project.

## Overview

The digital wallet system uses PostgreSQL as its primary database with a clean, normalized schema designed for financial transactions. The system employs a three-tier migration approach:

1. **GORM Auto-Migration**: Automatic schema generation from Go models
2. **DDL Migrations**: Explicit schema definitions with constraints and indexes
3. **DML Migrations**: Data seeding and initial records

## Database Schema

### Tables

#### 1. Wallets Table

Stores wallet information for users and system providers.

```sql
CREATE TABLE wallets (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL UNIQUE,
    acnt_type VARCHAR(50) NOT NULL CHECK (acnt_type IN ('user', 'provider')),
    balance BIGINT NOT NULL DEFAULT 0 CHECK (balance >= 0),
    status VARCHAR(50) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'suspended')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
```

**Fields:**
- `id`: Primary key (auto-increment)
- `user_id`: Unique identifier for wallet owner
- `acnt_type`: Account type (`user` or `provider`)
- `balance`: Current balance in cents (prevents floating-point precision issues)
- `status`: Wallet status (`active`, `inactive`, `suspended`)
- `created_at`: Record creation timestamp
- `updated_at`: Last modification timestamp (auto-updated via trigger)

**Constraints:**
- Unique constraint on `user_id`
- Check constraint ensuring balance is non-negative
- Check constraint for valid account types
- Check constraint for valid status values

#### 2. Transactions Table

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

**Wallets Table:**
- `idx_wallets_user_id`: Unique index on user_id (primary lookup)
- `idx_wallets_acnt_type`: Index on account type
- `idx_wallets_status`: Index on status

**Transactions Table:**
- `idx_transactions_type`: Index on transaction type
- `idx_transactions_status`: Index on status
- `idx_transactions_created_at`: Index on creation time

### Triggers

Automatic timestamp management:
- `update_wallets_updated_at`: Updates `updated_at` on wallet modifications
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
   - Reference data insertion
   - Sample data for testing

### Migration Files

#### DDL Migration
- `migrations/ddl/001_create_wallet_schema.sql`: Complete schema definition

#### DML Migration
- `migrations/dml/001_insert_provider_wallets.sql`: Seed data and sample records

### Running Migrations

```bash
# Run all migrations
make migrate

# Or run manually
go run main.go migrate
```

## Seed Data

### Provider Wallets (Required)

System provider wallets for transaction processing:

- **deposit-provider-master**: Source wallet for deposits (balance: 999,999,999,999 cents)
- **withdraw-provider-master**: Destination wallet for withdrawals (balance: 0 cents)

### Sample User Wallets

Demonstration user accounts with various states:

- **user-001**: Active user with 10,000 cents balance
- **user-002**: Active user with 5,000 cents balance
- **user-003**: Active user with 25,000 cents balance
- **user-004**: New user with 0 cents balance
- **user-inactive**: Inactive user with 1,000 cents balance

### Sample Transactions

Demonstration transaction history:

- Completed deposit: user-001, 5,000 cents
- Completed withdrawal: user-002, 1,000 cents
- Completed transfer: user-001 → user-003, 2,000 cents

## Database Configuration

### Development Environment

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

### Test Environment

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

## Data Integrity

### Financial Constraints

1. **Balance Validation**: Wallet balances cannot be negative
2. **Amount Validation**: Transaction amounts must be positive
3. **Status Validation**: Only valid status values are allowed
4. **Type Validation**: Only valid transaction and account types are allowed

### Transaction Consistency

1. **Double-Entry Bookkeeping**: Each transfer creates two transaction records
2. **Atomic Operations**: All balance updates occur within database transactions
3. **Audit Trail**: Complete transaction history is maintained
4. **Status Tracking**: Transaction status progression is tracked

## Performance Considerations

### Indexing Strategy

- Primary lookups optimized with unique indexes
- Composite indexes for complex query patterns
- Time-based indexes for transaction history queries

### Scalability
- Ensured locking granularity for high concurrency
- Efficient indexing for large datasets
- Minimal locking for concurrent operations

### Query Optimization

- Efficient wallet lookups by user_id
- Fast transaction filtering by type and status
- Optimized transaction history retrieval


## Future Enhancements

### Potential Improvements

1. **Partitioning**: Implement table partitioning for large transaction volumes
2. **Archiving**: Archive old transactions for performance
3. **Replication**: Set up read replicas for reporting
4. **Monitoring**: Add database performance monitoring
5. **Encryption**: Implement column-level encryption for sensitive data
