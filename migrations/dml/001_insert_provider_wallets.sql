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

-- =============================================================================
-- SAMPLE TRANSACTION DATA (FOR DEMONSTRATION)
-- =============================================================================

-- Insert sample completed transactions to demonstrate transaction history
-- These show examples of different transaction types with proper double-entry bookkeeping

-- Sample deposit transaction (debit from provider, credit to user)
-- Debit entry: Provider wallet loses money
INSERT INTO transactions (subject_wallet_id, object_wallet_id, transaction_type, operation_type, amount, status, created_at, updated_at)
VALUES ('deposit-provider-master', 'user-001', 'deposit', 'debit', 5000, 'completed', NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days')
ON CONFLICT DO NOTHING;

-- Credit entry: User wallet gains money
INSERT INTO transactions (subject_wallet_id, object_wallet_id, transaction_type, operation_type, amount, status, created_at, updated_at)
VALUES ('user-001', 'deposit-provider-master', 'deposit', 'credit', 5000, 'completed', NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days')
ON CONFLICT DO NOTHING;

-- Sample withdraw transaction (debit from user, credit to provider)
-- Debit entry: User wallet loses money
INSERT INTO transactions (subject_wallet_id, object_wallet_id, transaction_type, operation_type, amount, status, created_at, updated_at)
VALUES ('user-002', 'withdraw-provider-master', 'withdraw', 'debit', 1000, 'completed', NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day')
ON CONFLICT DO NOTHING;

-- Credit entry: Provider wallet gains money
INSERT INTO transactions (subject_wallet_id, object_wallet_id, transaction_type, operation_type, amount, status, created_at, updated_at)
VALUES ('withdraw-provider-master', 'user-002', 'withdraw', 'credit', 1000, 'completed', NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day')
ON CONFLICT DO NOTHING;

-- Sample transfer transaction (debit from sender, credit to receiver)
-- Debit entry: Sender wallet loses money
INSERT INTO transactions (subject_wallet_id, object_wallet_id, transaction_type, operation_type, amount, status, created_at, updated_at)
VALUES ('user-001', 'user-003', 'transfer', 'debit', 2000, 'completed', NOW() - INTERVAL '12 hours', NOW() - INTERVAL '12 hours')
ON CONFLICT DO NOTHING;

-- Credit entry: Receiver wallet gains money
INSERT INTO transactions (subject_wallet_id, object_wallet_id, transaction_type, operation_type, amount, status, created_at, updated_at)
VALUES ('user-003', 'user-001', 'transfer', 'credit', 2000, 'completed', NOW() - INTERVAL '12 hours', NOW() - INTERVAL '12 hours')
ON CONFLICT DO NOTHING;

-- Sample pending deposit transaction (both entries pending)
-- Debit entry: Provider wallet loses money (pending)
INSERT INTO transactions (subject_wallet_id, object_wallet_id, transaction_type, operation_type, amount, status, created_at, updated_at)
VALUES ('deposit-provider-master', 'user-002', 'deposit', 'debit', 3000, 'pending', NOW() - INTERVAL '1 hour', NOW() - INTERVAL '1 hour')
ON CONFLICT DO NOTHING;

-- Credit entry: User wallet gains money (pending)
INSERT INTO transactions (subject_wallet_id, object_wallet_id, transaction_type, operation_type, amount, status, created_at, updated_at)
VALUES ('user-002', 'deposit-provider-master', 'deposit', 'credit', 3000, 'pending', NOW() - INTERVAL '1 hour', NOW() - INTERVAL '1 hour')
ON CONFLICT DO NOTHING;

-- Additional completed transfer transaction (user-001 → user-003, 2,000 cents)
-- Debit entry: Sender wallet loses money
INSERT INTO transactions (subject_wallet_id, object_wallet_id, transaction_type, operation_type, amount, status, created_at, updated_at)
VALUES ('user-001', 'user-003', 'transfer', 'debit', 2000, 'completed', NOW() - INTERVAL '6 hours', NOW() - INTERVAL '6 hours')
ON CONFLICT DO NOTHING;

-- Credit entry: Receiver wallet gains money
INSERT INTO transactions (subject_wallet_id, object_wallet_id, transaction_type, operation_type, amount, status, created_at, updated_at)
VALUES ('user-003', 'user-001', 'transfer', 'credit', 2000, 'completed', NOW() - INTERVAL '6 hours', NOW() - INTERVAL '6 hours')
ON CONFLICT DO NOTHING;

-- =============================================================================
-- DATA VERIFICATION
-- =============================================================================

-- The following comments show expected data after seeding:
-- Provider Wallets: 2 (deposit-provider-master, withdraw-provider-master)
-- User Wallets: 5 (user-001 through user-004, plus user-inactive)
-- Sample Transactions: 12 (6 transaction pairs with proper double-entry bookkeeping)
-- Transaction Types: deposit (2 pairs), withdraw (1 pair), transfer (2 pairs), pending deposit (1 pair)
-- Total Balance in System: 40,000 cents (excluding provider balances)
