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
    id          SERIAL      PRIMARY KEY,
    fact_key    TEXT        NOT NULL,
    namespace   TEXT        NOT NULL,
    label       TEXT,
    description TEXT,
    unit        TEXT,

    UNIQUE(fact_key, namespace),
    CHECK (namespace IN ('dei', 'ifrs-full', 'invest', 'srt', 'us-gaap', 'custom'))
);

CREATE TABLE IF NOT EXISTS report_data (
  id        INTEGER         PRIMARY KEY,
  cik       INTEGER         NOT NULL,
  fact_id   INTEGER         NOT NULL,
  start     TIMESTAMP,
  end       TIMESTAMP,
  val       DECIMAL(20,4)   NOT NULL,
  accn      TEXT,
  fy        SMALLINT,
  fp        TEXT,
  form      TEXT            NOT NULL,
  filed     TIMESTAMP       NOT NULL,
  frame     TEXT,

  FOREIGN KEY (cik) REFERENCES company(cik),
  FOREIGN KEY (fact_id) REFERENCES fact(id)
);
