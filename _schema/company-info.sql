CREATE TABLE IF NOT EXISTS company (
    cik         INTEGER     PRIMARY KEY,
    name        TEXT        NOT NULL,
    ticker      TEXT        NOT NULL,
    next_10k    TIMESTAMP,
    next_10q    TIMESTAMP
);

CREATE TABLE IF NOT EXISTS fact (
    id          SERIAL      PRIMARY KEY,
    fact_key    TEXT        NOT NULL,
    namespace   TEXT        NOT NULL,
    label       TEXT,
    description TEXT,
    unit        TEXT,

    UNIQUE(fact_key, fact_type),
    CHECK (namespace IN ('dei', 'ifrs-full', 'invest', 'srt', 'us-gaap', 'custom'))
);

CREATE TABLE report_data (
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
