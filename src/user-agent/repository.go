package user_agent

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
)

func CreateUserAgentTable(db *sql.DB) {
	const TABLE_NAME = "db/user-agent-tbl.sql"

	sqlBytes, sqlBytesErr := os.ReadFile(TABLE_NAME)
	if sqlBytesErr != nil {
		log.Fatalf("Failed to read %s: %v", TABLE_NAME, sqlBytesErr)
	}
	_, queryErr := db.Exec(string(sqlBytes))
	if queryErr != nil {
		log.Fatalf("Failed to create %v: %v", TABLE_NAME, queryErr)
	}
}

func InsertUserAgent(ctx context.Context, db *sql.DB, name string, email string) error {
	const query = `
	INSERT INTO user_agent (id, name, email)
	VALUES (1, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
	  name = excluded.name,
	  email = excluded.email;
	`

	_, execErr := db.ExecContext(
		ctx,
		query,
		name, email,
	)

	if execErr != nil {
		return fmt.Errorf("Failed to insert new user: %w", execErr)
	}

	return nil
}

func GetUserAgent(ctx context.Context, db *sql.DB) (*UserAgent, error) {
	const query = `SELECT name, email FROM user_agent WHERE id=1;`
	var userAgent UserAgent

	queryErr := db.QueryRowContext(ctx, query).Scan(&userAgent.name, &userAgent.email)

	if queryErr != nil {
		if queryErr == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("Failed to get user: %w", queryErr)
	}

	return &userAgent, nil
}
