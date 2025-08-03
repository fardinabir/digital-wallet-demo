-- Migration: Insert provider wallets
-- Description: Pre-insert two provider wallets that act as master accounts
-- for deposit and withdraw operations to maintain transaction consistency

-- Insert deposit provider wallet
-- This wallet acts as the source for all deposit transactions
INSERT INTO wallets (user_id, acnt_type, balance, status, created_at, updated_at) VALUES ('deposit-provider-master', 'provider', 999999999999, 'active', NOW(), NOW()) ON CONFLICT (user_id) DO NOTHING;

-- Insert withdraw provider wallet
-- This wallet acts as the destination for all withdraw transactions
INSERT INTO wallets (user_id, acnt_type, balance, status, created_at, updated_at) VALUES ('withdraw-provider-master', 'provider', 0, 'active', NOW(), NOW()) ON CONFLICT (user_id) DO NOTHING;

-- Insert dummy user wallets for testing
-- These wallets can be used for testing deposit, withdraw, and transfer operations

-- Insert first dummy user wallet
INSERT INTO wallets (user_id, acnt_type, balance, status, created_at, updated_at) VALUES ('user-001', 'user', 10000, 'active', NOW(), NOW()) ON CONFLICT (user_id) DO NOTHING;

-- Insert second dummy user wallet
INSERT INTO wallets (user_id, acnt_type, balance, status, created_at, updated_at) VALUES ('user-002', 'user', 5000, 'active', NOW(), NOW()) ON CONFLICT (user_id) DO NOTHING;
