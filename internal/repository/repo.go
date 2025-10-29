package repository

import (
	"context"

	"github.com/Egor213/notifyBot/internal/entity"
	pgdb "github.com/Egor213/notifyBot/internal/repository/pg"
	"github.com/Egor213/notifyBot/pkg/postgres"
)

type Users interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id int64) (*entity.User, error)
}

type NotifySettings interface {
	GetByUser(ctx context.Context, tgID int64) ([]*entity.NotifySetting, error)
	Create(ctx context.Context, setting *entity.NotifySetting) error
	Delete(ctx context.Context, tgID int64, service string, level entity.LogLevel) error
}

type Repositories struct {
	Users
	NotifySettings
}

func NewRepositoriesPG(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		Users:          pgdb.NewUsersRepo(pg),
		NotifySettings: pgdb.NewNotifySettingsRepo(pg),
	}
}
