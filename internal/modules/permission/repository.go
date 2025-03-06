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

func (r *repo) FindAll(ctx context.Context) ([]Entity, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var permissions []Entity
	err := r.db.SelectContext(ctx, &permissions, "select * from permissions")
	if err != nil {
		return nil, fmt.Errorf("failed to find permissions: %w", err)
	}

	return permissions, nil
}
