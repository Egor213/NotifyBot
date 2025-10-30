package service

import (
	"context"

	"github.com/Egor213/notifyBot/internal/entity"
	"github.com/Egor213/notifyBot/internal/repository"
	"github.com/Egor213/notifyBot/internal/repository/repotypes"
)

type Users interface {
	RegisterUser(ctx context.Context, id int64, email string) (*entity.User, error)
	GetEmail(ctx context.Context, tgID int64) (string, error)
}

type NotifySettings interface {
	GetSettings(ctx context.Context, tgID int64) ([]*entity.NotifySetting, error)
	SetSettings(ctx context.Context, tgID int64, services []string, levels []entity.LogLevel) error
	RemoveSettings(ctx context.Context, tgID int64, service string, level entity.LogLevel) error
	GetChatIDsByFilters(ctx context.Context, filter repotypes.ChatIDFilter) ([]int64, error)
}

type Services struct {
	User           Users
	NotifySettings NotifySettings
}

type ServiceDep struct {
	Repos *repository.Repositories
}

func NewServices(dep *ServiceDep) *Services {
	return &Services{
		User:           NewUsers(dep.Repos.Users),
		NotifySettings: NewNotifySettingsService(dep.Repos.NotifySettings),
	}
}
