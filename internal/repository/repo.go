package repository

import (
	"github.com/Egor213/notifyBot/internal/entity"
	pgdb "github.com/Egor213/notifyBot/internal/repository/pg"
	"github.com/Egor213/notifyBot/pkg/postgres"
)

type Users interface {
	Create(user *entity.User) error
	GetByID(id int64) (*entity.User, error)
}

type Repositories struct {
	Users
}

func NewRepositoriesPG(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		Users: pgdb.NewUsersRepo(pg),
	}
}
