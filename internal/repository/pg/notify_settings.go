package pgdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/Egor213/notifyBot/internal/entity"
	"github.com/Egor213/notifyBot/internal/repository/repoerrs"
	"github.com/Egor213/notifyBot/internal/repository/repotypes"
	"github.com/Egor213/notifyBot/pkg/postgres"
	sq "github.com/Masterminds/squirrel"
)

type NotifySettingsRepo struct {
	pg *postgres.Postgres
}

func NewNotifySettingsRepo(pg *postgres.Postgres) *NotifySettingsRepo {
	return &NotifySettingsRepo{pg: pg}
}

func (r *NotifySettingsRepo) GetByUser(ctx context.Context, tgID int64) ([]*entity.NotifySetting, error) {
	query, args, err := r.pg.Builder.
		Select("service", "level").
		From("users_notify_settings").
		Where(sq.Eq{"tg_id": tgID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", repoerrs.ErrGetSettings, err)

	}

	rows, err := r.pg.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", repoerrs.ErrGetSettings, err)
	}
	defer rows.Close()

	var settings []*entity.NotifySetting
	for rows.Next() {
		setting := &entity.NotifySetting{}
		if err := rows.Scan(&setting.Service, &setting.Level); err != nil {
			return nil, fmt.Errorf("%w: %v", repoerrs.ErrGetSettings, err)
		}
		settings = append(settings, setting)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", repoerrs.ErrGetSettings, err)
	}

	return settings, nil
}

func (r *NotifySettingsRepo) Create(ctx context.Context, setting *entity.NotifySetting) error {
	query, args, err := r.pg.Builder.
		Insert("users_notify_settings").
		Columns("tg_id", "service", "level").
		Values(setting.TgID, setting.Service, setting.Level).
		Suffix("ON CONFLICT DO NOTHING").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("%w: %v", repoerrs.ErrCreateSetting, err)
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%w: %v", repoerrs.ErrCreateSetting, err)
	}

	return nil
}

func (r *NotifySettingsRepo) Delete(ctx context.Context, tgID int64, service string, level entity.LogLevel) error {
	query, args, err := r.pg.Builder.
		Delete("users_notify_settings").
		Where(sq.Eq{"tg_id": tgID, "service": service, "level": level}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("%w: %v", repoerrs.ErrDeleteSetting, err)
	}

	ct, err := r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%w: %v", repoerrs.ErrDeleteSetting, err)
	}
	if ct.RowsAffected() == 0 {
		return errors.New("setting not found")
	}

	return nil
}
func (r *NotifySettingsRepo) GetChatIDsByFilters(ctx context.Context, filter repotypes.ChatIDFilter) ([]int64, error) {
	conds := BuildGetChatIDQuery(filter)
	query := r.pg.Builder.
		Select("DISTINCT tg_id").
		From("users_notify_settings n").
		Join("users_mails m USING(tg_id)")

	if len(conds) > 0 {
		query = query.Where(sq.And(conds))
	}

	sql, args, _ := query.ToSql()
	rows, err := r.pg.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("NotifySettingsRepo.GetChatIDsByFilters - r.pg.Pool.Query: %w", err)
	}
	defer rows.Close()

	chatIDS := []int64{}
	for rows.Next() {
		var chatIDTemp int64
		if err := rows.Scan(&chatIDTemp); err != nil {
			return nil, fmt.Errorf("NotifySettingsRepo.GetChatIDsByFilters - rows.Scan: %w", err)
		}
		chatIDS = append(chatIDS, chatIDTemp)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("NotifySettingsRepo.GetChatIDsByFilters - rows.Err: %w", err)

	}

	return chatIDS, nil
}
