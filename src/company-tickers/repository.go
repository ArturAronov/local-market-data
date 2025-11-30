package company_ticker

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const (
	INSERT_TICKER_INFO = `
		INSERT OR IGNORE INTO company_tickers (
			cik,
			ticker,
			title
		) VALUES (?, ?, ?);`
)

func InsertTickerInfoR(data *secResponse) error {
	log.Println("[InsertTickerInfoR] Inserting company ticker info into db")
	db, dbErr := sql.Open("sqlite3", "_data/company-tickers.db")
	if dbErr != nil {
		return fmt.Errorf("[InsertTickerInfoR] Failed to open database: %w", dbErr)
	}

	defer db.Close()

	tx, txErr := db.Begin()
	if txErr != nil {
		return fmt.Errorf("[InsertTickerInfoR] Failed to start transaction: %w", txErr)
	}

	var deferErr error
	defer func() {
		panicErr := recover()
		if panicErr != nil {
			tx.Rollback()
			panic(panicErr)
		} else if deferErr != nil {
			tx.Rollback()
		} else {
			deferErr = tx.Commit()
		}
	}()

	query, queryErr := tx.Prepare(INSERT_TICKER_INFO)
	if queryErr != nil {
		return fmt.Errorf("[InsertTickerInfoR] Failed to prepare insert statement: %w", queryErr)
	}

	defer query.Close()

	var totalRowsAffected int64

	for _, entry := range *data {
		res, execErr := query.Exec(entry.CIKStr, entry.Ticker, entry.Title)
		if execErr != nil {
			return fmt.Errorf("[InsertTickerInfoR] Failed to insert data: %w", execErr)
		}

		rowsAffected, rowsAffectedErr := res.RowsAffected()
		if rowsAffectedErr != nil {
			log.Printf("Warning: Could not retrieve RowsAffected: %v", rowsAffectedErr)
		} else {
			totalRowsAffected += rowsAffected
		}
	}

	log.Printf("[InsertTickerInfoR] Completed inserting company ticker info, rows affected: %d", totalRowsAffected)

	return nil
}
