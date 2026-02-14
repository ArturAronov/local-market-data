package company_info

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	COUNT_COMPANY_FACT = `SELECT COUNT (fact_key) FROM fact WHERE cik = ?`
	UPDATE_COMPANY     = `
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
	SELECT_CIK_BY_TICKER = `SELECT cik FROM company WHERE ticker = ?;`
	SELECT_COMPANY       = `
		SELECT
			cik,
			sic,
			name,
			ticker,
			phone,
			entry_type,
			owner_org,
			exchanges,
			description,
			fiscal_year_end,
			latest_10k,
			latest_10q
		FROM company
		WHERE cik = ?;`
	SELECT_FACTS_REPORTS = `SELECT
		    f.fact_key,
		    f.cik,
		    f.namespace,
		    f.label,
		    f.description,
		    f.unit,
		    r.id,
		    r.start,
		    r.end,
		    r.val,
		    r.accn,
		    r.fy,
		    r.fp,
		    r.form,
		    r.filed,
		    r.frame,
		    r.hash
		FROM fact f
		LEFT JOIN report_data r
			ON f.fact_key = r.fact_key
			AND f.cik = r.cik
		WHERE f.cik = 1018724
		ORDER BY f.fact_key, r.filed DESC;`
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

func (r *Repository) GetCompanyFactsReportsR(cik int) ([]DbFactsReports, error) {
	rows, rowsErr := r.db.Query(SELECT_FACTS_REPORTS, cik)
	if rowsErr != nil {
		return nil, fmt.Errorf(
			"[GetCompanyFactsReportsR]: Failed to get company facts & reports: \n%w",
			rowsErr,
		)
	}
	defer rows.Close()

	factsMap := make(map[string]*DbFactsReports)
	for rows.Next() {
		var (
			factKey   string
			factCik   int
			namespace string
			label     string
			desc      string
			unit      string
			reportId  sql.NullInt64
			reportCik sql.NullInt64
			start     sql.NullTime
			end       sql.NullTime
			val       sql.NullFloat64
			accn      sql.NullString
			fy        sql.NullInt64
			fp        sql.NullString
			form      sql.NullString
			filed     sql.NullTime
			frame     sql.NullString
			hash      []byte
		)

		scanErr := rows.Scan(
			&factKey,
			&factCik,
			&namespace,
			&label,
			&desc,
			&unit,
			&reportId,
			&start,
			&end,
			&val,
			&accn,
			&fy,
			&fp,
			&form,
			&filed,
			&frame,
			&hash,
		)
		if scanErr != nil {
			return nil, fmt.Errorf("[GetCompanyFactsReportsR]: Failed to scan row: %w", scanErr)
		}

		if _, exists := factsMap[factKey]; !exists {
			fr := &DbFactsReports{
				Reports: []DbReport{},
			}
			fr.Cik = factCik
			fr.FactKey = factKey
			fr.Namespace = namespace
			fr.Label = label
			fr.Description = desc
			fr.Unit = unit
			factsMap[factKey] = fr
		}

		if reportId.Valid {
			var accnPtr, fpPtr, formPtr, framePtr *string
			if accn.Valid {
				accnPtr = &accn.String
			}
			if fp.Valid {
				fpPtr = &fp.String
			}
			if form.Valid {
				formPtr = &form.String
			}
			if frame.Valid {
				framePtr = &frame.String
			}

			var startPtr, endPtr *time.Time
			if start.Valid {
				startPtr = &start.Time
			}
			if end.Valid {
				endPtr = &end.Time
			}

			var fyPtr *int
			if fy.Valid {
				fyVal := int(fy.Int64)
				fyPtr = &fyVal
			}

			factsMap[factKey].Reports = append(factsMap[factKey].Reports, DbReport{
				Id:      int(reportId.Int64),
				Cik:     int(reportCik.Int64),
				FactKey: factKey,
				Start:   startPtr,
				End:     endPtr,
				Val:     val.Float64,
				Accn:    accnPtr,
				Fy:      fyPtr,
				Fp:      fpPtr,
				Form:    *formPtr,
				Filed:   filed.Time,
				Frame:   framePtr,
				Hash:    hash,
			})
		}
	}

	if rowsErr = rows.Err(); rowsErr != nil {
		return nil, fmt.Errorf("[GetCompanyFactsReportsR]: Error iterating rows: %w", rowsErr)
	}

	result := make([]DbFactsReports, 0, len(factsMap))
	for _, fr := range factsMap {
		result = append(result, *fr)
	}

	return result, nil
}

func (r *Repository) CountCompanyReportDataR(cik int) (int, error) {
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
			"[CountCompanyReportDataR] Failed to execute query SELECT_CIK_BY_TICKER: %s \n %w",
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
		&response.Sic,
		&response.Name,
		&response.Ticker,
		&response.Phone,
		&response.EntryType,
		&response.OwnerOrg,
		&response.Exchanges,
		&response.Description,
		&response.FiscalYearEnd,
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

func (r *Repository) InsertCompanyFinancialDataR(data []DbFact) error {
	if len(data) == 0 {
		log.Println("[InsertCompanyFinancialDataR] No facts to insert")
		return nil
	}

	log.Println("[InsertCompanyFinancialDataR] Inserting company facts into db")

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("[InsertCompanyFinancialDataR] Failed to start transaction: %w", err)
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
		return fmt.Errorf("[InsertCompanyFinancialDataR] Failed to prepare statement: %w", err)
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
			return fmt.Errorf("[InsertCompanyFinancialDataR] Failed to insert data: %w", err)
		}

		if rowsAffected, err := res.RowsAffected(); err == nil {
			totalRowsAffected += rowsAffected
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("[InsertCompanyFinancialDataR] Failed to commit: %w", err)
	}

	log.Printf(
		"[InsertCompanyFinancialDataR] Inserted %d facts (%d new rows)",
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
