package company_ticker

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	DELETE_COMPANY     = `DELETE FROM company_tickers WHERE cik = ?;`
	SELECT_ALL         = `SELECT cik, ticker, title FROM company_tickers;`
	INSERT_TICKER_INFO = `
		INSERT OR IGNORE INTO company_tickers (
			cik,
			ticker,
			title
		) VALUES (?, ?, ?);`
)

func DeleteTickerR(cik int64) error {
	db, dbErr := sql.Open("sqlite3", "_data/company-tickers.db")
	if dbErr != nil {
		return fmt.Errorf("[DeleteTickerR] Failed to open database: %w", dbErr)
	}

	defer db.Close()
	_, execErr := db.Exec(DELETE_COMPANY, cik)
	if execErr != nil {
		return fmt.Errorf(
			"[DeleteTickerR] Failed to delete company from company_tickers with cik %d. %w",
			cik,
			execErr,
		)
	}

	return nil
}

func GetTickerInfoR() ([]SecEntry, error) {
	db, dbErr := sql.Open("sqlite3", "_data/company-tickers.db")
	if dbErr != nil {
		return nil, fmt.Errorf("[GetTickerInfoR] Failed to open database: %w", dbErr)
	}

	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	rows, rowsErr := db.QueryContext(ctx, SELECT_ALL)
	if rowsErr != nil {
		return nil, fmt.Errorf(
			"[GetTickerInfoR] Failed to execute query SELECT_ALL: %s \n %w",
			SELECT_ALL,
			rowsErr,
		)
	}

	defer rows.Close()

	var response []SecEntry

	for rows.Next() {
		var _response SecEntry
		rowErr := rows.Scan(
			&_response.CIKStr,
			&_response.Ticker,
			&_response.Title,
		)
		if rowErr != nil {
			return nil, fmt.Errorf("[GetTickerInfoR] Failed to scan row: %w\n", rowErr)
		}

		response = append(response, _response)
	}

	queryErr := rows.Err()
	if queryErr != nil {
		return nil, fmt.Errorf("[GetTickerInfoR] Error in rows query: %w\n", queryErr)
	}

	return response, nil
}

func InsertTickerInfoR(data *SecResponse) error {
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
