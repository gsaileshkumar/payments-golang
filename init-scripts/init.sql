-- Create the accounts table
CREATE TABLE accounts (
    id NUMERIC PRIMARY KEY,
    balance DECIMAL(21, 6) NOT NULL
);

-- Create the transactions table to log all transactions for audit
CREATE TABLE transactions (
    transaction_id SERIAL PRIMARY KEY,  -- Auto-incrementing ID for each transaction
    source_account_id NUMERIC NOT NULL,
    destination_account_id NUMERIC NOT NULL,
    amount DECIMAL(21, 6) NOT NULL, 
    transaction_date TIMESTAMP WITH TIME ZONE DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),  -- Timestamp of the transaction in UTC
    FOREIGN KEY (source_account_id) REFERENCES accounts(id),
    FOREIGN KEY (destination_account_id) REFERENCES accounts(id)
);
