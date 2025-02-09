package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/magneless/todo-app/internal/models"
	"github.com/magneless/todo-app/internal/storage"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(name, username, password_hash string) (int64, error) {
	const op = "ropository.CreateUser"

	stmt, err := r.db.Prepare(`
	INSERT INTO users(name, username, password_hash) VALUES($1, $2, $3) RETURNING id
	`)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	defer stmt.Close()

	var id int64

	err = stmt.QueryRow(name, username, password_hash).Scan(&id)
	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == storage.UniqueViolationErrorCode {
		return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
	}
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (r *Repository) GetUser(username, password_hash string) (*models.User, error) {
	const op = "repository.GetUser"

	stmt, err := r.db.Prepare(`
	SELECT id, name, username, password_hash FROM users WHERE username = $1 and password_hash = $2
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer stmt.Close()

	user := &models.User{}
	err = stmt.QueryRow(username, password_hash).Scan(
		&user.ID, &user.Name, &user.Username, &user.PasswrodHash,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("%s: wrong username or password", op)
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (r *Repository) AddRefreshToken(refresh_token string, expires_at time.Time, username string) error {
	const op = "repository.AddRefreshToken"

	stmt, err := r.db.Prepare(`
	UPDATE users
	SET refresh_token = $1, expires_at = $2
	WHERE username = $3
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(refresh_token, expires_at, username)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
