package repository

import (
	"sync"

	"github.com/Egor213/notifyBot/internal/entity"
	"github.com/Egor213/notifyBot/internal/repository/repoerrs"
)

type InMemoryUserRepo struct {
	data map[int64]*entity.User
	mu   sync.RWMutex
}

func NewInMemoryUserRepo() *InMemoryUserRepo {
	return &InMemoryUserRepo{
		data: make(map[int64]*entity.User),
	}
}

func (r *InMemoryUserRepo) Create(user *entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.data[user.TgID]; exists {
		return repoerrs.ErrAlreadyExists
	}
	r.data[user.TgID] = user
	return nil
}

func (r *InMemoryUserRepo) GetByID(id int64) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.data[id]
	if !ok {
		return nil, nil
	}
	return user, nil
}
