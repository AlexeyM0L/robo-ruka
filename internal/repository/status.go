package repository

import (
	"errors"
	"io/fs"
	"os"
	"sync"

	"robo-ruka/internal/domain"
)

type StatusRepository struct {
	path string
	mu   sync.RWMutex
}

func NewStatusRepository(path string) *StatusRepository {
	return &StatusRepository{path: path}
}

func (r *StatusRepository) Get() (domain.Status, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	data, err := os.ReadFile(r.path)
	if errors.Is(err, fs.ErrNotExist) {
		return domain.StatusOff, nil
	}
	if err != nil {
		return "", err
	}
	s, ok := domain.ParseStatus(string(data))
	if !ok {
		return domain.StatusOff, nil
	}
	return s, nil
}

func (r *StatusRepository) Set(s domain.Status) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return os.WriteFile(r.path, []byte(s), 0o644)
}
