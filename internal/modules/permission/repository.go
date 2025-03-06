package permission

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type repo struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) RepositoryInterface {
	return &repo{db}
}

func (r *repo) FindByRoleID(ctx context.Context, roleId string) ([]Entity, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var permissions []Entity
	err := r.db.SelectContext(
		ctx,
		&permissions,
		`
		SELECT p.*
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = $1
		`,
		roleId,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find permissions: %w", err)
	}

	return permissions, nil
}

func (r *repo) FindAll(ctx context.Context) ([]Entity, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var permissions []Entity
	err := r.db.SelectContext(ctx, &permissions, "SELECT * FROM permissions ORDER BY key ASC")
	if err != nil {
		return nil, fmt.Errorf("failed to find permissions: %w", err)
	}

	return permissions, nil
}
