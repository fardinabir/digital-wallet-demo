-- Digital Transaction Database Schema
-- This file contains the database schema for the transaction service
-- It defines the transactions table with proper constraints and indexes

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
CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status);
CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at);

-- Create updated_at trigger function for automatic timestamp updates
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for automatic updated_at timestamp updates
DROP TRIGGER IF EXISTS update_transactions_updated_at ON transactions;
CREATE TRIGGER update_transactions_updated_at
    BEFORE UPDATE ON transactions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Add comments to tables and columns for documentation

COMMENT ON TABLE transactions IS 'Transaction records for all wallet operations';
COMMENT ON COLUMN transactions.id IS 'Primary key for transaction records';
COMMENT ON COLUMN transactions.subject_wallet_id IS 'Txn ID that initiated the transaction';
COMMENT ON COLUMN transactions.object_wallet_id IS 'Target wallet ID for transfers (provider wallet ID for deposits/withdrawals)';
COMMENT ON COLUMN transactions.transaction_type IS 'Type of transaction: deposit, withdraw, or transfer';
COMMENT ON COLUMN transactions.operation_type IS 'Operation type: debit or credit';
COMMENT ON COLUMN transactions.amount IS 'Transaction amount in cents';
COMMENT ON COLUMN transactions.status IS 'Transaction status: pending, completed, failed, or cancelled';