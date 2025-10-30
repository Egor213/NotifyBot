package pgdb

import (
	"github.com/Egor213/notifyBot/internal/repository/repotypes"
	sq "github.com/Masterminds/squirrel"
)

func BuildGetChatIDQuery(filter repotypes.ChatIDFilter) []sq.Sqlizer {
	conds := []sq.Sqlizer{}

	if filter.Level != "" {
		conds = append(conds, sq.Eq{"n.level": filter.Level})
	}

	if filter.Service != "" {
		conds = append(conds, sq.Eq{"n.service": filter.Service})
	}

	return conds
}
