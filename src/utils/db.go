package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

type DbTables string

const (
	USER_EMAIL   DbTables = "user-email"
	COMPANY_INFO DbTables = "company-info"
)

func CreateTable(tableName DbTables, db *sql.DB) error {
	tablePath := fmt.Sprintf("_schema/%s.sql", tableName)

	sqlBytes, sqlBytesErr := os.ReadFile(tablePath)
	if sqlBytesErr != nil {
		return fmt.Errorf("[CreateTable] Failed to read file %s: %v", tableName, sqlBytesErr)
	}

	_, queryErr := db.Exec(string(sqlBytes))
	if queryErr != nil {
		return fmt.Errorf("[CreateTable] Failed to create %v: %v", tablePath, queryErr)
	}

	return nil
}

func InitDB() (map[DbTables]*sql.DB, error) {
	log.Println("[InitDB] Initializing databases")
	var dbPaths = []DbTables{
		COMPANY_INFO,
		USER_EMAIL,
	}

	dbMap := make(map[DbTables]*sql.DB)
	for _, path := range dbPaths {
		sqlPath := fmt.Sprintf("_data/%s.db", path)
		log.Printf("[InitDB] Creating new db %s: \n", sqlPath)

		db, dbErr := sql.Open("sqlite3", sqlPath)
		if dbErr != nil {
			return nil, fmt.Errorf("[InitDB] Failed to open %s: %w", path, dbErr)
		}

		createTblErr := CreateTable(path, db)
		if createTblErr != nil {
			return nil, fmt.Errorf("[InitDB] Failed to create new db %s: %w", path, createTblErr)
		}

		dbMap[path] = db
	}

	log.Println("[InitDB] Finished initializing databases")

	return dbMap, nil
}
