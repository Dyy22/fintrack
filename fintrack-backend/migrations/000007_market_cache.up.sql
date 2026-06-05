CREATE TABLE stock_quotes (
    symbol VARCHAR(12) PRIMARY KEY,
    name VARCHAR(255),
    price NUMERIC(15,2) NOT NULL CHECK (price > 0),
    currency VARCHAR(8) NOT NULL DEFAULT 'IDR',
    source VARCHAR(100) NOT NULL DEFAULT 'unknown',
    fetched_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE market_charts (
    symbol VARCHAR(12) NOT NULL,
    chart_range VARCHAR(16) NOT NULL,
    interval VARCHAR(16) NOT NULL,
    name VARCHAR(255),
    currency VARCHAR(8) NOT NULL DEFAULT 'IDR',
    source VARCHAR(100) NOT NULL DEFAULT 'unknown',
    fetched_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    points JSONB NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (symbol, chart_range, interval)
);

CREATE INDEX idx_stock_quotes_fetched_at ON stock_quotes(fetched_at DESC);
CREATE INDEX idx_market_charts_fetched_at ON market_charts(fetched_at DESC);
