-- Create DB
CREATE DATABASE tinkoff;

-- Change DB
\c tinkoff;

-- Tinkoff Instruments meta information
CREATE TABLE IF NOT EXISTS instrument (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    figi VARCHAR(255) NOT NULL,
    ticker VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    min_price_increment FLOAT NOT NULL,
    currency VARCHAR(255) NOT NULL,
    type VARCHAR(255) NOT NULL,
    global_rank BOOLEAN NOT NULL -- Is included intoForbes global rank
);

-- Candles Interval meta
CREATE TABLE IF NOT EXISTS candle_interval (
    id SMALLINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

-- Candles Data
CREATE TABLE IF NOT EXISTS candle (
    ts timestamptz NOT NULL,
    instrument_id INTEGER REFERENCES instrument (id) NOT NULL,
    interval_id SMALLINT REFERENCES candle_interval (id) NOT NULL,
    open_price REAL NULL,
    close_price REAL NULL,
    high_price REAL NULL,
    low_price REAL NULL,
    volume REAL NULL
);

-- Create Partitioned table by time
SELECT create_hypertable
    ('candle',
     'ts',
     chunk_time_interval => interval '1 week');