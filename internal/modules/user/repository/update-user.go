package userrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/bernardinorafael/internal/modules/user"
)

func (r repo) Update(ctx context.Context, entity user.Entity) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	entity.Updated = time.Now()
	_, err := r.db.NamedExecContext(
		ctx,
		`
		UPDATE users
		SET
			full_name = :full_name,
			username = :username,
			updated = :updated
		WHERE
			id = :id
	`,
		entity,
	)
	if err != nil {
		return fmt.Errorf("error on update user: %w", err)
	}

	return nil
}
