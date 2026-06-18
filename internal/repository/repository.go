package repository

import (
	"database/sql"

	"robo-ruka/internal/domain"
)

type Status interface {
	Get() (domain.Status, error)
	Set(domain.Status) error
}

type Repository struct {
	Status Status
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Status: NewStatusRepository(db),
	}
}
