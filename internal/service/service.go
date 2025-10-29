package service

import (
	"github.com/Egor213/notifyBot/internal/entity"
	"github.com/Egor213/notifyBot/internal/repository"
	"github.com/Egor213/notifyBot/internal/service/users"
)

type Users interface {
	RegisterUser(id int64, email string) (*entity.User, error)
}

type Services struct {
	User Users
}

type ServiceDep struct {
	Repos *repository.Repositories
}

func NewServices(dep *ServiceDep) *Services {
	return &Services{
		User: users.New(dep.Repos.Users),
	}
}
