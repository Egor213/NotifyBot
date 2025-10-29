package service

import (
	"context"

	"github.com/Egor213/notifyBot/internal/entity"
	"github.com/Egor213/notifyBot/internal/repository"
)

type Users interface {
	RegisterUser(ctx context.Context, id int64, email string) (*entity.User, error)
	GetEmail(ctx context.Context, tgID int64) (string, error)
}

type Services struct {
	User Users
}

type ServiceDep struct {
	Repos *repository.Repositories
}

func NewServices(dep *ServiceDep) *Services {
	return &Services{
		User: NewUsers(dep.Repos.Users),
	}
}
