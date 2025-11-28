package utils

import (
	"database/sql"
	"fmt"
	"os"
)

func CreateTable(tableName string, db *sql.DB) error {
	var tablePath string
	fileExtension := ".sql"

	if len(tableName) <= len(fileExtension) {
		return fmt.Errorf("[CreateTable] tableName length is too short")
	}

	if tableName[len(tableName)-len(fileExtension):] == fileExtension {
		tablePath = fmt.Sprintf("db/%s", tableName)
	} else {
		tablePath = fmt.Sprintf("db/%s.sql", tableName)
	}

	sqlBytes, sqlBytesErr := os.ReadFile(tablePath)
	if sqlBytesErr != nil {
		return fmt.Errorf("[CreateTable] Failed to read %s: %v", tablePath, sqlBytesErr)
	}
	_, queryErr := db.Exec(string(sqlBytes))
	if queryErr != nil {
		return fmt.Errorf("[CreateTable] Failed to create %v: %v", tablePath, queryErr)
	}

	return nil
}
