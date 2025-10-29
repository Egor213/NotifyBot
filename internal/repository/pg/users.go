package pgdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Egor213/notifyBot/internal/entity"
	"github.com/Egor213/notifyBot/internal/repository/repoerrs"
	"github.com/Egor213/notifyBot/pkg/postgres"
	sq "github.com/Masterminds/squirrel"
)

type UsersRepo struct {
	pg *postgres.Postgres
}

func NewUsersRepo(pg *postgres.Postgres) *UsersRepo {
	return &UsersRepo{pg: pg}
}

func (r *UsersRepo) Create(ctx context.Context, user *entity.User) error {
	query, args, err := r.pg.Builder.
		Insert("users_mails").
		Columns("tg_id", "email").
		Values(user.TgID, user.Email).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("%w: %v", repoerrs.ErrCreateUser, err)
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%w: %v", repoerrs.ErrCreateUser, err)
	}

	return nil
}

func (r *UsersRepo) GetByID(ctx context.Context, tgID int64) (*entity.User, error) {
	query, args, err := r.pg.Builder.
		Select("tg_id", "email").
		From("users_mails").
		Where(sq.Eq{"tg_id": tgID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", repoerrs.ErrGetUser, err)
	}

	row := r.pg.Pool.QueryRow(ctx, query, args...)
	user := &entity.User{}
	if err := row.Scan(&user.TgID, &user.Email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repoerrs.ErrUserNotFound
		}
		return nil, fmt.Errorf("%w: %v", repoerrs.ErrGetUser, err)
	}

	return user, nil
}
