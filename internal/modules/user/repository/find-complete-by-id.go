package userrepo

import (
	"context"
	"time"

	"github.com/bernardinorafael/internal/modules/email"
	"github.com/bernardinorafael/internal/modules/phone"
	"github.com/bernardinorafael/internal/modules/user"
	"github.com/bernardinorafael/pkg/transaction"
	"github.com/jmoiron/sqlx"
)

func (r repo) FindCompleteByID(ctx context.Context, userId string) (*user.CompleteEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	usr := user.CompleteEntity{
		User:   user.Entity{},
		Emails: []email.AdditionalEmail{},
		Phones: []phone.AdditionalPhone{},
	}

	err := transaction.ExecTx(ctx, r.db, func(tx *sqlx.Tx) error {
		err := tx.GetContext(
			ctx,
			&usr.User,
			"SELECT * FROM users WHERE id = $1 ORDER BY created DESC",
			userId,
		)
		if err != nil {
			return err
		}

		err = tx.SelectContext(
			ctx,
			&usr.Emails,
			"SELECT * FROM emails WHERE user_id = $1 ORDER BY created DESC",
			userId,
		)
		if err != nil {
			return err
		}

		err = tx.SelectContext(
			ctx,
			&usr.Phones,
			"SELECT * FROM phones WHERE user_id = $1 ORDER BY created DESC",
			userId,
		)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &usr, nil
}
