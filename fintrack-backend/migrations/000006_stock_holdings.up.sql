ALTER TABLE accounts
ADD COLUMN stock_symbol VARCHAR(12),
ADD COLUMN stock_lots NUMERIC(18,4),
ADD COLUMN stock_price_per_share NUMERIC(15,2);

UPDATE account_types
SET description = 'Indonesian stock holdings tracked in lots and valued using the latest market price'
WHERE name = 'stock_broker';

CREATE INDEX idx_accounts_stock_symbol ON accounts(stock_symbol)
WHERE stock_symbol IS NOT NULL;
