package userrepo

import (
	"context"
	"time"
)

func (r repo) Delete(ctx context.Context, userId string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", userId)
	if err != nil {
		return err
	}

	return nil
}
