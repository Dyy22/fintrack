CREATE TABLE gold_price_history (
    price_date DATE PRIMARY KEY,
    price_per_gram NUMERIC(15,2) NOT NULL CHECK (price_per_gram > 0),
    source VARCHAR(100) NOT NULL DEFAULT 'manual',
    fetched_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO gold_price_history (price_date, price_per_gram, source, fetched_at)
SELECT fetched_at::date, price_per_gram, source, fetched_at
FROM gold_prices
ON CONFLICT (price_date) DO UPDATE
SET price_per_gram = EXCLUDED.price_per_gram,
    source = EXCLUDED.source,
    fetched_at = EXCLUDED.fetched_at;
