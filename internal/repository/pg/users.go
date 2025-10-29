package pgdb

import (
	"github.com/Egor213/notifyBot/internal/entity"
	"github.com/Egor213/notifyBot/pkg/postgres"
)

type UsersRepo struct {
	pg *postgres.Postgres
}

func NewUsersRepo(pg *postgres.Postgres) *UsersRepo {
	return &UsersRepo{
		pg: pg,
	}
}

func (r *UsersRepo) Create(user *entity.User) error {
	return nil
}

func (r *UsersRepo) GetByID(id int64) (*entity.User, error) {
	return &entity.User{}, nil
}
