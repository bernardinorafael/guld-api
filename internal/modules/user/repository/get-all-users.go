package userrepo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bernardinorafael/internal/_shared/dto"
	"github.com/bernardinorafael/internal/modules/user"
)

func (r repo) GetAll(ctx context.Context, params dto.SearchParams) ([]user.EntityWithTeam, int, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var count int
	err := r.db.GetContext(ctx, &count, "SELECT COUNT(*) FROM users")
	if err != nil {
		return nil, -1, err
	}

	direction := "DESC"
	sort := params.Sort

	if strings.HasPrefix(sort, "-") {
		direction = "ASC"
		sort = strings.TrimPrefix(sort, "-")
	}

	sql := fmt.Sprintf(`
		SELECT
			u.*,
			t.id as "team.id",
			t.name as "team.name"
		FROM users u
		LEFT JOIN team_members tm ON tm.user_id = u.id
		LEFT JOIN teams t ON t.id = tm.team_id
		WHERE (
			(to_tsvector('simple', u.full_name) || to_tsvector('simple', u.username))
				@@ websearch_to_tsquery('simple', $3)
				OR u.full_name ILIKE '%%' || $3 || '%%'
				OR u.username ILIKE '%%' || $3 || '%%'
		)
		ORDER BY %s %s
		LIMIT $1 OFFSET $2
	`, sort, direction)

	users := make([]user.EntityWithTeam, 0)
	skip := (params.Page - 1) * params.Limit

	err = r.db.SelectContext(ctx, &users, sql, params.Limit, skip, params.Query)
	if err != nil {
		return nil, -1, err
	}

	return users, count, nil
}
