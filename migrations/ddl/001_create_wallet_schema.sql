-- Digital Wallet Database Schema
-- This file contains the complete database schema for the digital wallet system
-- It defines tables for wallets and transactions with proper constraints and indexes

-- Create wallets table
CREATE TABLE IF NOT EXISTS wallets (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    acnt_type VARCHAR(50) NOT NULL CHECK (acnt_type IN ('user', 'provider')),
    balance BIGINT NOT NULL DEFAULT 0 CHECK (balance >= 0),
    status VARCHAR(50) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'suspended')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create unique index on user_id for wallets
CREATE UNIQUE INDEX IF NOT EXISTS idx_wallets_user_id ON wallets(user_id);

-- Create index on account type for efficient filtering
CREATE INDEX IF NOT EXISTS idx_wallets_acnt_type ON wallets(acnt_type);

-- Create index on status for efficient filtering
CREATE INDEX IF NOT EXISTS idx_wallets_status ON wallets(status);

-- Create transactions table
CREATE TABLE IF NOT EXISTS transactions (
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

-- Create indexes for transactions table
CREATE INDEX IF NOT EXISTS idx_transactions_subject_wallet_id ON transactions(subject_wallet_id);
CREATE INDEX IF NOT EXISTS idx_transactions_object_wallet_id ON transactions(object_wallet_id);
CREATE INDEX IF NOT EXISTS idx_transactions_type ON transactions(transaction_type);
CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status);
CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at);

-- Add foreign key constraints (optional, but good practice)
-- Note: These are commented out as the current implementation uses string-based wallet IDs
-- If you want to enforce referential integrity, uncomment these lines
-- ALTER TABLE transactions ADD CONSTRAINT fk_transactions_subject_wallet 
--     FOREIGN KEY (subject_wallet_id) REFERENCES wallets(user_id) ON DELETE CASCADE;
-- ALTER TABLE transactions ADD CONSTRAINT fk_transactions_object_wallet 
--     FOREIGN KEY (object_wallet_id) REFERENCES wallets(user_id) ON DELETE CASCADE;

-- Create updated_at trigger function for automatic timestamp updates
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for automatic updated_at timestamp updates
DROP TRIGGER IF EXISTS update_wallets_updated_at ON wallets;
CREATE TRIGGER update_wallets_updated_at
    BEFORE UPDATE ON wallets
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_transactions_updated_at ON transactions;
CREATE TRIGGER update_transactions_updated_at
    BEFORE UPDATE ON transactions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Add comments to tables and columns for documentation
COMMENT ON TABLE wallets IS 'Digital wallet accounts for users and providers';
COMMENT ON COLUMN wallets.id IS 'Primary key for wallet records';
COMMENT ON COLUMN wallets.user_id IS 'Unique identifier for wallet owner';
COMMENT ON COLUMN wallets.acnt_type IS 'Account type: user or provider';
COMMENT ON COLUMN wallets.balance IS 'Current wallet balance in cents';
COMMENT ON COLUMN wallets.status IS 'Wallet status: active, inactive, or suspended';

COMMENT ON TABLE transactions IS 'Transaction records for all wallet operations';
COMMENT ON COLUMN transactions.id IS 'Primary key for transaction records';
COMMENT ON COLUMN transactions.subject_wallet_id IS 'Wallet ID that initiated the transaction';
COMMENT ON COLUMN transactions.object_wallet_id IS 'Target wallet ID for transfers (provider wallet ID for deposits/withdrawals)';
COMMENT ON COLUMN transactions.transaction_type IS 'Type of transaction: deposit, withdraw, or transfer';
COMMENT ON COLUMN transactions.operation_type IS 'Operation type: debit or credit';
COMMENT ON COLUMN transactions.amount IS 'Transaction amount in cents';
COMMENT ON COLUMN transactions.status IS 'Transaction status: pending, completed, failed, or cancelled';