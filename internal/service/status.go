package service

import (
	"robo-ruka/internal/domain"
	"robo-ruka/internal/repository"
)

type statusService struct {
	repo *repository.Repository
}

func NewStatusService(repo *repository.Repository) Status {
	return &statusService{repo: repo}
}

func (s *statusService) Current() (domain.Status, error) {
	return s.repo.Status.Get()
}

func (s *statusService) Update(raw string) (domain.Status, error) {
	st, ok := domain.ParseStatus(raw)
	if !ok {
		return "", ErrInvalidStatus
	}
	if err := s.repo.Status.Set(st); err != nil {
		return "", err
	}
	return st, nil
}
