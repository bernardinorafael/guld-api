package userrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/bernardinorafael/internal/modules/email"
	"github.com/bernardinorafael/internal/modules/user"
	"github.com/bernardinorafael/pkg/transaction"
	"github.com/jmoiron/sqlx"
)

func (r repo) FindCompleteByID(ctx context.Context, userId string) (*user.CompleteEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	u := user.CompleteEntity{
		User:   user.Entity{},
		Emails: []email.Entity{},
	}

	err := transaction.ExecTx(ctx, r.db, func(tx *sqlx.Tx) error {
		err := tx.GetContext(
			ctx,
			&u.User,
			"SELECT * FROM users WHERE id = $1",
			userId,
		)
		if err != nil {
			return fmt.Errorf("error on find user by id: %w", err)
		}

		err = tx.SelectContext(
			ctx,
			&u.Emails,
			"SELECT * FROM emails WHERE user_id = $1 ORDER BY created DESC",
			userId,
		)
		if err != nil {
			return fmt.Errorf("error on find emails by user id: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error on find complete user by id: %w", err)
	}

	return &u, nil
}
