CREATE TABLE IF NOT EXISTS company (
    cik             INTEGER     PRIMARY KEY,
    sic             TEXT,
    name            TEXT,
    ticker          TEXT,
    phone           TEXT,
    entry_type      TEXT,
    owner_org       TEXT,
    exchanges       TEXT,
    description     TEXT,
    fiscal_year_end TEXT,
    latest_10k      TEXT,
    latest_10q      TEXT
);

CREATE TABLE IF NOT EXISTS fact (
    fact_key    TEXT        PRIMARY KEY,
    cik         INTEGER     NOT NULL,
    namespace   TEXT        NOT NULL,
    label       TEXT,
    description TEXT,
    unit        TEXT,

    UNIQUE(cik, fact_key, namespace),
    FOREIGN KEY (cik) REFERENCES company(cik),
    CHECK (namespace IN ('dei', 'ifrs-full', 'invest', 'srt', 'us-gaap', 'custom'))
);

CREATE TABLE IF NOT EXISTS report_data (
    id        INTEGER         PRIMARY KEY,
    cik       INTEGER         NOT NULL,
    fact_key  TEXT            NOT NULL,
    start     TIMESTAMP,
    end       TIMESTAMP,
    val       DECIMAL(20,4)   NOT NULL,
    accn      TEXT            NOT NULL,
    fy        SMALLINT,
    fp        TEXT,
    form      TEXT            NOT NULL,
    filed     TIMESTAMP       NOT NULL,
    frame     TEXT,
    hash      BLOB            NOT NULL,

    UNIQUE(cik, hash),
    FOREIGN KEY (cik) REFERENCES company(cik),
    FOREIGN KEY (fact_key) REFERENCES fact(fact_key)
);
