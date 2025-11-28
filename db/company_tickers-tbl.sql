CREATE TABLE IF NOT EXISTS company_tickers (
    cik     INTEGER PRIMARY KEY,
    ticker  TEXT    NOT NULL UNIQUE,
    title   TEXT    NOT NULL
);
