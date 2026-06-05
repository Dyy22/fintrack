DROP INDEX IF EXISTS idx_accounts_stock_symbol;

UPDATE account_types
SET description = 'Stock brokerage account tracked as manual IDR balance'
WHERE name = 'stock_broker';

ALTER TABLE accounts
DROP COLUMN IF EXISTS stock_price_per_share,
DROP COLUMN IF EXISTS stock_lots,
DROP COLUMN IF EXISTS stock_symbol;
