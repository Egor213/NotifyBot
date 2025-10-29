package users

import (
	"errors"

	"github.com/Egor213/notifyBot/internal/entity"
	"github.com/Egor213/notifyBot/internal/repository"
)

type UserService struct {
	userRepo repository.Users
}

func New(repo repository.Users) *UserService {
	return &UserService{userRepo: repo}
}

func (s *UserService) RegisterUser(id int64, email string) (*entity.User, error) {
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}

	user := &entity.User{
		ID:    id,
		Email: email,
	}

	err := s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
