package user

import (
	"database/sql"
	"fmt"
)

const (
	INSERT_USER_EMAIL = `
		INSERT INTO user_email (id, email)
		VALUES (1, ?)
		ON CONFLICT(id) DO UPDATE SET
		  email = excluded.email;`

	GET_USER_EMAIL = `
		SELECT email
		FROM user_email
		WHERE id=1;`
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) InsertUserEmail(email string) error {
	_, execErr := r.db.Exec(
		INSERT_USER_EMAIL,
		email,
	)

	if execErr != nil {
		return fmt.Errorf("[InsertUserEmail] Failed to insert new user: %w", execErr)
	}

	return nil
}

func (r *Repository) GetUserEmail() (*string, error) {
	var email string

	queryErr := r.db.QueryRow(GET_USER_EMAIL).Scan(&email)

	if email == "" {
		return nil, fmt.Errorf("No user email set. Run 'init --email <your-email>'")
	}

	if queryErr != nil {
		if queryErr == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("Failed to get user: %w", queryErr)
	}

	return &email, nil
}
