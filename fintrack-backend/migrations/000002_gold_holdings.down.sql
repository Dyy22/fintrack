UPDATE account_types
SET description = 'Physical gold holdings tracked as manual IDR balance'
WHERE name = 'gold';

DROP TABLE IF EXISTS gold_prices;

ALTER TABLE accounts
DROP COLUMN IF EXISTS gold_price_per_gram,
DROP COLUMN IF EXISTS gold_grams;
