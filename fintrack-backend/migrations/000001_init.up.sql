CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE account_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT
);

INSERT INTO account_types (name, description) VALUES
('bank', 'Bank account'),
('ewallet', 'E-Wallet (GoPay, OVO, Dana, etc)'),
('cash', 'Physical cash'),
('gold', 'Physical gold holdings tracked as manual IDR balance'),
('stock_broker', 'Stock brokerage account tracked as manual IDR balance');

CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    account_type_id INT NOT NULL REFERENCES account_types(id),
    name VARCHAR(100) NOT NULL,
    balance NUMERIC(15,2) NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'IDR',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_accounts_user_id ON accounts(user_id);

CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(50) NOT NULL,
    type VARCHAR(10) NOT NULL CHECK (type IN ('income', 'expense')),
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_categories_user_id ON categories(user_id);

CREATE UNIQUE INDEX idx_categories_user_name_type
ON categories(user_id, LOWER(name), type)
WHERE user_id IS NOT NULL;

CREATE UNIQUE INDEX idx_categories_default_name_type
ON categories(LOWER(name), type)
WHERE is_default = TRUE;

INSERT INTO categories (name, type, is_default) VALUES
('Salary', 'income', TRUE),
('Freelance', 'income', TRUE),
('Investment', 'income', TRUE),
('Food', 'expense', TRUE),
('Transport', 'expense', TRUE),
('Shopping', 'expense', TRUE),
('Bills', 'expense', TRUE),
('Entertainment', 'expense', TRUE);

CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    type VARCHAR(10) NOT NULL CHECK (type IN ('income', 'expense', 'transfer')),
    amount NUMERIC(15,2) NOT NULL CHECK (amount > 0),
    description TEXT,
    date TIMESTAMPTZ NOT NULL,
    transfer_account_id UUID REFERENCES accounts(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_transaction_transfer_rules CHECK (
        (
            type = 'transfer'
            AND transfer_account_id IS NOT NULL
            AND transfer_account_id <> account_id
            AND category_id IS NULL
        )
        OR
        (
            type IN ('income', 'expense')
            AND transfer_account_id IS NULL
        )
    )
);

CREATE INDEX idx_transactions_user_id ON transactions(user_id);
CREATE INDEX idx_transactions_account_id ON transactions(account_id);
CREATE INDEX idx_transactions_date ON transactions(date);
