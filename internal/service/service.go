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

type State interface {
	SetState(chatID int64, state entity.StateType, data map[any]any)
	GetState(chatID int64) *UserState
	ClearState(chatID int64)
}

type MailSender interface {
	SendMessage(to string, title string, body string) error
}

type Services struct {
	User           Users
	NotifySettings NotifySettings
	State          State
	MailSender     MailSender
}

type SendMailDep struct {
	SendMail  string
	Port      int
	Protocol  string
	SecretKey string
}

type ServiceDep struct {
	Repos       *repository.Repositories
	SendMailDep SendMailDep
}

func NewServices(dep *ServiceDep) *Services {
	return &Services{
		User:           NewUsers(dep.Repos.Users),
		NotifySettings: NewNotifySettingsService(dep.Repos.NotifySettings),
		State:          NewStateService(),
		MailSender:     NewMailSenderService(dep.SendMailDep),
	}
}
