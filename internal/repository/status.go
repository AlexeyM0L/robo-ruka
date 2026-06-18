package repository

import (
	"database/sql"
	"errors"

	"robo-ruka/internal/domain"
)

type StatusRepository struct {
	db *sql.DB
}

func NewStatusRepository(db *sql.DB) *StatusRepository {
	return &StatusRepository{db: db}
}

func (r *StatusRepository) Get() (domain.Status, error) {
	var value string
	err := r.db.QueryRow(`SELECT value FROM status WHERE id = 1`).Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.StatusOff, nil
	}
	if err != nil {
		return "", err
	}
	s, ok := domain.ParseStatus(value)
	if !ok {
		return domain.StatusOff, nil
	}
	return s, nil
}

func (r *StatusRepository) Set(s domain.Status) error {
	_, err := r.db.Exec(
		`INSERT INTO status (id, value) VALUES (1, ?)
		 ON CONFLICT(id) DO UPDATE SET value = excluded.value`,
		string(s),
	)
	return err
}
