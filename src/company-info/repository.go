package company_info

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const (
	SELECT_CIK_BY_TICKER = `SELECT cik FROM company WHERE ticker = ?;`
	COUNT_COMPANY_FACT   = `SELECT COUNT (fact_key) FROM fact WHERE cik = ?`
	UPDATE_COMPANY       = `
		UPDATE company
		SET
			sic = ?,
			name = ?,
			phone = ?,
			entry_type = ?,
			owner_org = ?,
			exchanges = ?,
			description = ?,
			fiscal_year_end = ?,
			latest_10k = ?,
			latest_10q = ?
		WHERE cik = ?;`
	SELECT_COMPANY = `
		SELECT
			cik,
			name,
			ticker,
			latest_10k,
			latest_10q
		FROM company
		WHERE cik = ?;`
	INSERT_COMPANY_INFO = `
		INSERT OR IGNORE INTO company (
			cik,
			ticker,
			name
		) VALUES (?, ?, ?);`
	INSERT_COMPANY_FACT = `
		INSERT OR IGNORE INTO fact (
			cik,
			fact_key,
			namespace,
			label,
			description,
			unit
		) VALUES (?, ?, ?, ?, ?, ?);`
	INSERT_COMPANY_REPORT = `
		INSERT OR IGNORE INTO report_data (
			cik,
			fact_key,
			start,
			end,
			val,
			accn,
			fy,
			fp,
			form,
			filed,
			frame,
			hash
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CountCompanyFactsR(cik int) (int, error) {
	var count int
	rowErr := r.db.QueryRow(
		COUNT_COMPANY_FACT,
		cik,
	).Scan(&count)

	if count == 0 {
		return count, nil
	}

	if rowErr != nil {
		return 0, fmt.Errorf(
			"[CountCompanyFactsR] Failed to execute query SELECT_CIK_BY_TICKER: %s \n %w",
			COUNT_COMPANY_FACT,
			rowErr,
		)
	}

	return count, nil
}

func (r *Repository) GetCikByTickerR(ticker string) (string, error) {
	var cik string
	rowErr := r.db.QueryRow(
		SELECT_CIK_BY_TICKER,
		strings.ToUpper(ticker),
	).Scan(&cik)

	if cik == "" {
		return "", fmt.Errorf(
			"[GetCikByTicker] No CIK found for ticker: %s", strings.ToUpper(ticker),
		)
	}

	if rowErr != nil {
		return "", fmt.Errorf(
			"[GetCikByTicker] Failed to execute query SELECT_CIK_BY_TICKER: %s \n %w",
			SELECT_CIK_BY_TICKER,
			rowErr,
		)
	}

	if cik == "" {
		return "", fmt.Errorf("[GetCikByTicker] No CIK found for ticker: %s\n", ticker)
	}

	return cik, nil
}

func (r *Repository) GetCompanyR(cik int) (*DbCompany, error) {
	var response DbCompany
	rowErr := r.db.QueryRow(
		SELECT_COMPANY,
		cik,
	).Scan(
		&response.Cik,
		&response.Name,
		&response.Ticker,
		&response.Latest10k,
		&response.Latest10q,
	)

	if response.Cik == 0 {
		return nil, nil
	}

	if rowErr != nil {
		return nil, fmt.Errorf(
			"[GetCompanyR] Failed to execute query SELECT_COMPANY: %s \n %w",
			SELECT_COMPANY,
			rowErr,
		)
	}

	return &response, nil
}

func (r *Repository) UpdateCompanyR(company DbCompany) error {
	result, resultErr := r.db.Exec(
		UPDATE_COMPANY,
		company.Sic,
		company.Name,
		company.Phone,
		company.EntryType,
		company.OwnerOrg,
		company.Exchanges,
		company.Description,
		company.FiscalYearEnd,
		company.Latest10k,
		company.Latest10q,
		company.Cik,
	)

	if resultErr != nil {
		return fmt.Errorf(
			"[UpdateCompanyR] Failed to execute query UPDATE_COMPANY: %s \n %w",
			UPDATE_COMPANY,
			resultErr,
		)
	}

	affected, affectedErr := result.RowsAffected()
	if affectedErr != nil {
		return fmt.Errorf("[UpdateCompanyR] Cannot get affected rows: %w", affectedErr)
	}

	if affected == 0 {
		return fmt.Errorf("[UpdateCompanyR] No rows affected")
	}

	return nil
}

func (r *Repository) InsertCompanyFactsR(data []DbFact) error {
	if len(data) == 0 {
		log.Println("[InsertCompanyFactsR] No facts to insert")
		return nil
	}

	log.Println("[InsertCompanyFactsR] Inserting company facts into db")

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("[InsertCompanyFactsR] Failed to start transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	stmt, err := tx.Prepare(INSERT_COMPANY_FACT)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("[InsertCompanyFactsR] Failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	var totalRowsAffected int64

	for _, entry := range data {
		res, err := stmt.Exec(
			entry.Cik,
			entry.FactKey,
			entry.Namespace,
			entry.Label,
			entry.Description,
			entry.Unit,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("[InsertCompanyFactsR] Failed to insert data: %w", err)
		}

		if rowsAffected, err := res.RowsAffected(); err == nil {
			totalRowsAffected += rowsAffected
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("[InsertCompanyFactsR] Failed to commit: %w", err)
	}

	log.Printf(
		"[InsertCompanyFactsR] Inserted %d facts (%d new rows)",
		len(data), totalRowsAffected,
	)

	return nil
}

func (r *Repository) InsertCompanyReportsR(data []DbReport) error {
	if len(data) == 0 {
		log.Println("[InsertCompanyReportsR] No facts to insert")
		return nil
	}

	log.Println("[InsertCompanyReportsR] Inserting company facts into db")

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("[InsertCompanyReportsR] Failed to start transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	stmt, err := tx.Prepare(INSERT_COMPANY_REPORT)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("[InsertCompanyReportsR] Failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	var totalRowsAffected int64

	for _, entry := range data {
		res, err := stmt.Exec(
			entry.Cik,
			entry.FactKey,
			entry.Start,
			entry.End,
			entry.Val,
			entry.Accn,
			entry.Fy,
			entry.Fp,
			entry.Form,
			entry.Filed,
			entry.Frame,
			entry.Hash,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("[InsertCompanyReportsR] Failed to insert data: %w", err)
		}

		if rowsAffected, err := res.RowsAffected(); err == nil {
			totalRowsAffected += rowsAffected
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("[InsertCompanyReportsR] Failed to commit: %w", err)
	}

	log.Printf(
		"[InsertCompanyReportsR] Inserted %d reports (%d new rows)",
		len(data), totalRowsAffected,
	)

	return nil
}

func (r *Repository) InsertTickerInfoR(data *SecEntryRes) error {
	log.Println("[InsertTickerInfoR] Inserting company ticker info into db")

	tx, txErr := r.db.Begin()
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

	query, queryErr := tx.Prepare(INSERT_COMPANY_INFO)
	if queryErr != nil {
		return fmt.Errorf(
			"[InsertTickerInfoR] Failed to prepare insert statement: %w",
			queryErr,
		)
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

	log.Printf(
		"[InsertTickerInfoR] Completed inserting company ticker info, rows affected: %d",
		totalRowsAffected,
	)

	return nil
}
