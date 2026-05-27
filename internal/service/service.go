package service

import (
	"robo-ruka/internal/domain"
	"robo-ruka/internal/repository"
)

type Status interface {
	Current() (domain.Status, error)
	Update(raw string) (domain.Status, error)
}

type Service struct {
	Status Status
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Status: NewStatusService(repo),
	}
}
