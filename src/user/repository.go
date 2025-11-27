package user

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
)

func CreateUserTable(db *sql.DB) {
	const TABLE_NAME = "db/user_email-tbl.sql"

	sqlBytes, sqlBytesErr := os.ReadFile(TABLE_NAME)
	if sqlBytesErr != nil {
		log.Fatalf("Failed to read %s: %v", TABLE_NAME, sqlBytesErr)
	}
	_, queryErr := db.Exec(string(sqlBytes))
	if queryErr != nil {
		log.Fatalf("Failed to create %v: %v", TABLE_NAME, queryErr)
	}
}

func InsertUserEmail(ctx context.Context, db *sql.DB, name string, email string) error {
	const query = `
	INSERT INTO user (id, email)
	VALUES (1, ?)
	ON CONFLICT(id) DO UPDATE SET
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

func GetUserEmail(ctx context.Context, db *sql.DB) (*string, error) {
	const query = `SELECT email FROM user_email WHERE id=1;`
	var email string

	queryErr := db.QueryRowContext(ctx, query).Scan(&email)

	if queryErr != nil {
		if queryErr == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("Failed to get user: %w", queryErr)
	}

	return &email, nil
}
