package permission

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type repo struct{ db *sqlx.DB }

func NewRepository(db *sqlx.DB) RepositoryInterface {
	return &repo{db}
}

func (r *repo) Delete(ctx context.Context, teamId string, permissionId string) error {
	_, err := r.db.ExecContext(
		ctx,
		`
		DELETE from permissions WHERE team_id = $1 AND id = $2
		`,
		teamId,
		permissionId,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r repo) Insert(ctx context.Context, permission Entity) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	_, err := r.db.NamedExecContext(
		ctx,
		`
		INSERT INTO permissions (team_id, name, key, description)
		VALUES (:team_id, :name, :key, :description)
		`,
		permission,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r repo) GetAll(ctx context.Context, teamId string, p PermissionSearchParams) ([]Entity, int, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var count int
	err := r.db.GetContext(ctx, &count, "SELECT COUNT(*) FROM permissions")
	if err != nil {
		return nil, -1, err
	}

	direction := "DESC"
	sort := p.Sort

	if strings.HasPrefix(sort, "-") {
		direction = "ASC"
		sort = strings.TrimPrefix(sort, "-")
	}

	sql := fmt.Sprintf(`
		SELECT * FROM permissions p
		WHERE (
			(to_tsvector('simple', p.name) ||
			to_tsvector('simple', p.description) ||
			to_tsvector('simple', p.key)) @@ websearch_to_tsquery('simple', $3)
				OR p.name ILIKE '%%' || $3 || '%%'
				OR p.description ILIKE '%%' || $3 || '%%'
				OR p.key ILIKE '%%' || $3 || '%%'
		)
		AND team_id = $4
		ORDER BY %s %s
		LIMIT $1 OFFSET $2
	`, sort, direction)

	permissions := make([]Entity, 0)
	skip := (p.Page - 1) * p.Limit

	err = r.db.SelectContext(ctx, &permissions, sql, p.Limit, skip, p.Query, teamId)
	if err != nil {
		return nil, -1, err
	}

	return permissions, count, nil
}
