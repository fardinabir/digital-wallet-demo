-- Migration: Insert provider wallets
-- Description: Pre-insert two provider wallets that act as master accounts
-- for deposit and withdraw operations to maintain transaction consistency

-- Insert deposit provider wallet
-- This wallet acts as the source for all deposit transactions
INSERT INTO wallets (user_id, acnt_type, balance, status, created_at, updated_at)
SELECT 'deposit-provider-master', 'provider', 999999999999, 'active', NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM wallets WHERE user_id = 'deposit-provider-master'
);

-- Insert withdraw provider wallet
-- This wallet acts as the destination for all withdraw transactions
INSERT INTO wallets (user_id, acnt_type, balance, status, created_at, updated_at)
SELECT 'withdraw-provider-master', 'provider', 0, 'active', NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM wallets WHERE user_id = 'withdraw-provider-master'
);