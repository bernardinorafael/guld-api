package userrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/bernardinorafael/internal/modules/user"
)

func (r repo) FindByID(ctx context.Context, userId string) (*user.Entity, error) {
	usr := user.Entity{}

	err := r.db.GetContext(
		ctx,
		&usr,
		"SELECT * FROM users WHERE id = $1",
		userId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}

	return &usr, nil
}
