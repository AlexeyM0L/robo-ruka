package repository

import "robo-ruka/internal/domain"

type Status interface {
	Get() (domain.Status, error)
	Set(domain.Status) error
}

type Repository struct {
	Status Status
}

func NewRepository(statePath string) *Repository {
	return &Repository{
		Status: NewStatusRepository(statePath),
	}
}
