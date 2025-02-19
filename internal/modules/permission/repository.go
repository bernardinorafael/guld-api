package permission

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type repo struct{ db *sqlx.DB }

func NewRepository(db *sqlx.DB) RepositoryInterface {
	return &repo{db}
}

func (r *repo) GetByID(ctx context.Context, orgId, permId string) (*Entity, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var permission Entity
	err := r.db.GetContext(
		ctx,
		&permission,
		"SELECT * FROM permissions WHERE org_id = $1 AND id = $2",
		orgId,
		permId,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get permission: %w", err)
	}

	return &permission, nil
}

func (r *repo) Update(ctx context.Context, permission Entity) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	permission.Updated = time.Now()
	_, err := r.db.NamedExecContext(
		ctx,
		`
		UPDATE permissions
		SET
			name = :name,
			description = :description,
			key = :key,
			updated = :updated
		WHERE org_id = :org_id AND id = :id
		`,
		permission,
	)
	if err != nil {
		return fmt.Errorf("failed to update permission: %w", err)
	}

	return nil
}
func (r *repo) Delete(ctx context.Context, orgId string, permId string) error {
	_, err := r.db.ExecContext(
		ctx,
		`DELETE from permissions WHERE org_id = $1 AND id = $2`,
		orgId,
		permId,
	)
	if err != nil {
		return fmt.Errorf("failed to delete permission: %w", err)
	}

	return nil
}

func (r repo) Insert(ctx context.Context, permission Entity) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	_, err := r.db.NamedExecContext(
		ctx,
		`
		INSERT INTO permissions (
			id,
			org_id,
			name,
			key,
			description,
			created,
			updated
		)
		VALUES (
			:id,
			:org_id,
			:name,
			:key,
			:description,
			:created,
			:updated
		)
		`,
		permission,
	)
	if err != nil {
		return fmt.Errorf("failed to insert permission: %w", err)
	}

	return nil
}

func (r repo) GetAll(ctx context.Context, orgId string, p PermissionSearchParams) ([]Entity, int, error) {
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
		AND org_id = $4
		ORDER BY %s %s
		LIMIT $1 OFFSET $2
	`, sort, direction)

	permissions := make([]Entity, 0)
	skip := (p.Page - 1) * p.Limit

	err = r.db.SelectContext(ctx, &permissions, sql, p.Limit, skip, p.Query, orgId)
	if err != nil {
		return nil, -1, err
	}

	return permissions, count, nil
}
