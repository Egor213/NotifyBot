package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Egor213/notifyBot/internal/entity"
	"github.com/Egor213/notifyBot/internal/repository"
	"github.com/Egor213/notifyBot/internal/repository/repoerrs"
	"github.com/Egor213/notifyBot/internal/service/srverrs"
)

type UserService struct {
	userRepo repository.Users
}

func NewUsers(repo repository.Users) *UserService {
	return &UserService{userRepo: repo}
}

func (s *UserService) RegisterUser(ctx context.Context, tgID int64, email string) (*entity.User, error) {
	existing, err := s.userRepo.GetByID(ctx, tgID)
	if err != nil && !errors.Is(err, repoerrs.ErrUserNotFound) {
		return nil, fmt.Errorf("%w: %v", srverrs.ErrUserCheckFailed, err)
	}
	if existing != nil {
		return nil, srverrs.ErrUserAlreadyExists
	}
	user := &entity.User{
		TgID:  tgID,
		Email: email,
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("%w: %v", srverrs.ErrUserCreateFailed, err)
	}

	return user, nil
}

func (s *UserService) GetEmail(ctx context.Context, tgID int64) (string, error) {
	user, err := s.userRepo.GetByID(ctx, tgID)
	if err != nil {
		if errors.Is(err, repoerrs.ErrUserNotFound) {
			return "", srverrs.ErrUserNotFound
		}
		return "", fmt.Errorf("%w: %v", srverrs.ErrUserCheckFailed, err)
	}
	return user.Email, nil
}
