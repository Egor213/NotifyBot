package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Egor213/notifyBot/internal/entity"
	"github.com/Egor213/notifyBot/internal/repository"
	"github.com/Egor213/notifyBot/internal/service/srverrs"
)

type NotifySettingsService struct {
	repo repository.NotifySettings
}

func NewNotifySettingsService(repo repository.NotifySettings) *NotifySettingsService {
	return &NotifySettingsService{repo: repo}
}

func (s *NotifySettingsService) GetSettings(ctx context.Context, tgID int64) ([]*entity.NotifySetting, error) {
	settings, err := s.repo.GetByUser(ctx, tgID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", srverrs.ErrGetNotifySettings, err)
	}
	return settings, nil
}

func (s *NotifySettingsService) SetSettings(ctx context.Context, tgID int64, services []string, levels []entity.LogLevel) error {
	if len(services) == 0 || len(levels) == 0 {
		return srverrs.ErrSetNotifySettings
	}

	for _, svc := range services {
		for _, lvl := range levels {
			setting := &entity.NotifySetting{
				TgID:    tgID,
				Service: svc,
				Level:   lvl,
			}
			if err := s.repo.Create(ctx, setting); err != nil {
				return fmt.Errorf("%w: %v", srverrs.ErrSetNotifySettings, err)
			}
		}
	}
	return nil
}

func (s *NotifySettingsService) RemoveSettings(ctx context.Context, tgID int64, service string, level entity.LogLevel) error {
	err := s.repo.Delete(ctx, tgID, service, level)
	if err != nil {
		if errors.Is(err, errors.New("setting not found")) {
			return srverrs.ErrNotifySettingNotFound
		}
		return fmt.Errorf("%w: %v", srverrs.ErrRemoveNotifySettings, err)
	}
	return nil
}
