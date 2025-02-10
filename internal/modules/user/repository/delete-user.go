package userrepo

import (
	"context"
)

func (r repo) Delete(ctx context.Context, userId string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", userId)
	if err != nil {
		return err
	}

	return nil
}
