-- Digital Wallet Seed Data
-- Description: Initial data for the digital wallet system including provider wallets and sample user accounts
-- This file contains essential seed data for system operation and testing

-- =============================================================================
-- PROVIDER WALLETS (REQUIRED FOR SYSTEM OPERATION)
-- =============================================================================

-- Insert deposit provider wallet
-- This wallet acts as the source for all deposit transactions
-- It has a large balance to handle deposit operations
INSERT INTO wallets (user_id, acnt_type, balance, status, created_at, updated_at)
VALUES ('deposit-provider-master', 'provider', 999999999999, 'active', NOW(), NOW())
ON CONFLICT (user_id) DO NOTHING;

-- Insert withdraw provider wallet
-- This wallet acts as the destination for all withdraw transactions
-- It starts with zero balance as it receives withdrawn funds
INSERT INTO wallets (user_id, acnt_type, balance, status, created_at, updated_at)
VALUES ('withdraw-provider-master', 'provider', 0, 'active', NOW(), NOW())
ON CONFLICT (user_id) DO NOTHING;

-- =============================================================================
-- SAMPLE USER WALLETS (FOR TESTING AND DEMONSTRATION)
-- =============================================================================

-- Insert sample user wallets with different balances and statuses
-- These wallets demonstrate various account states and can be used for testing

-- Active user with moderate balance
INSERT INTO wallets (user_id, acnt_type, balance, status, created_at, updated_at)
VALUES ('user-001', 'user', 10000, 'active', NOW(), NOW())
ON CONFLICT (user_id) DO NOTHING;

-- Active user with lower balance
INSERT INTO wallets (user_id, acnt_type, balance, status, created_at, updated_at)
VALUES ('user-002', 'user', 5000, 'active', NOW(), NOW())
ON CONFLICT (user_id) DO NOTHING;

-- Active user with higher balance
INSERT INTO wallets (user_id, acnt_type, balance, status, created_at, updated_at)
VALUES ('user-003', 'user', 25000, 'active', NOW(), NOW())
ON CONFLICT (user_id) DO NOTHING;

-- New user with zero balance
INSERT INTO wallets (user_id, acnt_type, balance, status, created_at, updated_at)
VALUES ('user-004', 'user', 0, 'active', NOW(), NOW())
ON CONFLICT (user_id) DO NOTHING;

-- Inactive user wallet (for testing status scenarios)
INSERT INTO wallets (user_id, acnt_type, balance, status, created_at, updated_at)
VALUES ('user-inactive', 'user', 1000, 'inactive', NOW(), NOW())
ON CONFLICT (user_id) DO NOTHING;
