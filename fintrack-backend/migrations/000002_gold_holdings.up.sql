ALTER TABLE accounts
ADD COLUMN gold_grams NUMERIC(18,4),
ADD COLUMN gold_price_per_gram NUMERIC(15,2);

CREATE TABLE gold_prices (
    id SMALLINT PRIMARY KEY DEFAULT 1 CHECK (id = 1),
    price_per_gram NUMERIC(15,2) NOT NULL CHECK (price_per_gram > 0),
    source VARCHAR(100) NOT NULL DEFAULT 'manual',
    fetched_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

UPDATE account_types
SET description = 'Physical gold holdings tracked in grams and valued using the latest configured Antam gold price'
WHERE name = 'gold';
