package user

import (
	"context"
	"database/sql"
	"fmt"
)

const (
	INSERT_USER_EMAIL = `
		INSERT INTO user_email (id, email)
		VALUES (1, ?)
		ON CONFLICT(id) DO UPDATE SET
		  email = excluded.email;`
)

func InsertUserEmail(email string) error {
	db, dbErr := sql.Open("sqlite3", "_data/user-email.db")
	if dbErr != nil {
		return fmt.Errorf("[InsertUserEmail] Failed to open database: %w", dbErr)
	}

	_, execErr := db.Exec(
		INSERT_USER_EMAIL,
		email,
	)

	if execErr != nil {
		return fmt.Errorf("[InsertUserEmail] Failed to insert new user: %w", execErr)
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
		return nil, fmt.Errorf("[GetUserEmail] Failed to get user: %w", queryErr)
	}

	return &email, nil
}
