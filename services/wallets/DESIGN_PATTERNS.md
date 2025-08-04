# Design Patterns & Concurrency Reference

## Design Patterns Implementation

### 1. Repository Pattern
**Location**: `internal/repository/wallet.go`

**Implementation**:
- Interface: `Wallet` interface
- Concrete: `wallet` struct implementation

**Problem Solved**:
Abstracts data access layer, enabling easy testing and database switching

**Reference**: Used in `internal/service/wallet.go` for dependency injection

### 2. Dependency Injection Pattern
**Location**: `internal/server/api.go`

**Implementation**:
`initWalletController()` function

**Problem Solved**:
Loose coupling between layers (Controller → Service → Repository)

**Flow**: Repository → Service → Controller chain with interface-based dependencies

### 3. Factory Pattern
**Locations**: 
- `internal/model/wallet.go`
- `internal/client/transaction_client.go`

**Implementation**:
- `NewWallet()` function for wallet creation
- `NewTxnClient()` function for singleton transaction client
**Problem Solved**:
- Encapsulates wallet creation logic with proper defaults
- Provides clean interface for transaction client instantiation

**Usage**: 
- `NewWallet()` called in `internal/controller/wallet.go` for wallet instantiation
- `NewTxnClient()` called in `internal/service/wallet.go` for transaction operations

### 4. Strategy Pattern
**Location**: `internal/service/wallet.go`

**Implementation**:
`Wallet` interface with concrete implementation

**Problem Solved**:
Allows different business logic implementations without changing controllers

**Benefit**: Easy to extend with new wallet types or business rules

### 5. Singleton Pattern
**Locations**: 
- `internal/client/transaction_client.go`
- `internal/cache/redis.go`

**Implementation**:
- Thread-safe singleton using `sync.Once`
- Factory method: `NewTxnClient()` function
- Instance caching with private variables

**Problem Solved**:
Ensures only one transaction client instance exists, reducing resource usage

**Concurrency Safety**: Uses `sync.Once` to guarantee thread-safe initialization

**Usage**: Called via `client.NewTxnClient().CreateTransactionPair()` in service layer

#### Redis Singleton Pattern with sync.Once

```go
var (
    redisInstance RedisClient
    redisOnce     sync.Once
)

func NewRedisClient() RedisClient {
    redisOnce.Do(func() {
        // Initialize Redis client once
        redisInstance = &redisClient{...}
    })
    return redisInstance
}
```

**Benefits**:
- **Thread-safe initialization**: Prevents race conditions
- **Resource efficiency**: Single connection pool shared across application
- **Consistent configuration**: Global config applied once

## Concurrency Implementation

### 1. Goroutines for Server Management
**Location**: `cmd/server.go`

**Implementation**:
Multiple servers started concurrently in `runServer()` function

**Problem Solved**:
Non-blocking server startup for API and Swagger servers

**Pattern**: `go func()` with error handling for each server

### 2. Goroutines for Transaction History Fetching
**Location**: `internal/service/wallet.go`

**Implementation**:
Asynchronous transaction history retrieval in `GetWalletWithTransactions()` method

**Problem Solved**:
Non-blocking external service calls for transaction data

**Pattern**: Concurrent HTTP requests to transaction microservice while processing wallet data

### 3. Database Row-Level Locking
**Location**: `internal/repository/wallet.go`

**Implementation**:
`UpdateWalletBalance()` function

**Concurrency Mechanism**:
`clause.Locking{Strength: "UPDATE"}`

**Problem Solved**:
Prevents race conditions during concurrent balance updates

**Critical Section**: Wallet balance modification with exclusive lock

### 4. Context-Based Graceful Shutdown
**Location**: `cmd/server.go`

**Implementation**:
- Signal context: `signal.NotifyContext()` in `runServer()`
- Timeout context: `context.WithTimeout()` for graceful shutdown

**Problem Solved**:
Clean server shutdown without losing in-flight requests

**Concurrency Control**: Coordinates shutdown across multiple goroutines

### 5. Database Transaction Management
**Location**: `internal/service/wallet.go`

**Implementation**:
- Transaction begin: `BeginTransaction()` in `Deposit()`, `Withdraw()`, and `Transfer()` methods
- Defer rollback: `defer func()` with panic recovery

**Problem Solved**:
ACID compliance for multi-step wallet operations

**Concurrency Safety**: Ensures atomic operations across multiple database writes

## Key Benefits

1. **Testability**: Repository pattern enables easy mocking
2. **Scalability**: Goroutines handle concurrent requests efficiently
3. **Data Integrity**: Row locking prevents balance corruption
4. **Maintainability**: Clean separation of concerns through dependency injection
5. **Reliability**: Graceful shutdown ensures no data loss during deployment

## Thread Safety Guarantees

- **Wallet Balance Updates**: Protected by database row-level exclusive locks
- **Transaction Creation**: Atomic within database transactions
- **Server Lifecycle**: Coordinated through context cancellation
- **Concurrent Requests**: Handled safely through Echo framework's built-in goroutine pool