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

	query := `
		UPDATE users
		SET
			full_name = :full_name,
			username = :username,
			updated = :updated,
			username_last_updated = :username_last_updated,
			username_lockout_end = :username_lockout_end,
			locked = :locked,
			banned = :banned,
			avatar_url = :avatar_url,
			phone_number = :phone_number,
			email_address = :email_address
		WHERE
			id = :id
	`

	_, err := r.db.NamedExecContext(ctx, query, entity)
	if err != nil {
		return fmt.Errorf("error on update user: %w", err)
	}

	return nil
}
