# Kong API Gateway Documentation

## Overview
Kong API Gateway has been successfully integrated into the digital wallet microservices architecture, providing a unified entry point with rate limiting and request routing.

## Available Endpoints

### Base URL
- **Kong Gateway**: `http://localhost:8000`
- **Kong Admin API**: `http://localhost:8001`
- **Kong Manager**: `http://localhost:8002`

### Wallet Operations

#### 1. Create User Account
```bash
POST http://localhost:8000/users
Content-Type: application/json

{
  "user_id": "john_doe",
  "acnt_type": "user"
}
```

#### 2. Deposit Funds
```bash
POST http://localhost:8000/wallets/deposit
Content-Type: application/json

{
  "user_id": "john_doe",
  "amount": 5000
}
```
**Note**: Amount is in cents (5000 = $50.00)

#### 3. Withdraw Funds
```bash
POST http://localhost:8000/wallets/withdraw
Content-Type: application/json

{
  "user_id": "john_doe",
  "amount": 1000
}
```
**Note**: Amount is in cents (1000 = $10.00)

#### 4. Transfer Between Users
```bash
POST http://localhost:8000/wallets/transfer
Content-Type: application/json

{
  "from_user_id": "john_doe",
  "to_user_id": "jane_doe",
  "amount": 1500
}
```
**Note**: Amount is in cents (1500 = $15.00)

#### 5. Check Wallet Balance & Transaction History
```bash
GET http://localhost:8000/wallets/{user_id}
```

## Rate Limiting

The Kong API Gateway implements global rate limiting:

### Global Rate Limits
- **100 requests per minute** across all endpoints
- **1000 requests per hour** across all endpoints

Rate limits are enforced using Kong's local policy and are fault-tolerant. This single rate limiting policy applies to all API endpoints for simplicity and consistency.

## Service Architecture
- **Kong Gateway**: Port 8000 (HTTP), 8443 (HTTPS)
- **Wallet Service**: Internal routing to `wallet-app:8081`
- **Transactions Service**: Internal routing to `transactions-app:8082`
- **Database**: Shared PostgreSQL on port 5432

## Testing Examples

### Complete Workflow Test
```bash
# 1. Create a user
curl -X POST http://localhost:8000/users \
  -H "Content-Type: application/json" \
  -d '{"user_id": "test-user", "acnt_type": "user"}'

# 2. Deposit funds
curl -X POST http://localhost:8000/wallets/deposit \
  -H "Content-Type: application/json" \
  -d '{"user_id": "test-user", "amount": 10000}'

# 3. Check balance
curl http://localhost:8000/api/v1/wallets/test-user

# 4. Create second user
curl -X POST http://localhost:8000/users \
  -H "Content-Type: application/json" \
  -d '{"user_id": "test-user-2", "acnt_type": "user"}'

# 5. Transfer funds
curl -X POST http://localhost:8000/wallets/transfer \
  -H "Content-Type: application/json" \
  -d '{"from_user_id": "test-user", "to_user_id": "test-user-2", "amount": 2500}'
```

## Configuration Files
- **Docker Compose**: `docker-compose.yml` - Kong service configuration
- **Kong Declarative Config**: `kong.yaml` - API routes and plugins

## Management
- **Start Services**: `docker-compose up -d`
- **Restart Kong**: `docker-compose restart kong`
- **View Logs**: `docker logs kong`
- **Kong Admin API**: Access via `http://localhost:8001`