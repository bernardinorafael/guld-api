package userrepo

import (
	"context"
	"time"

	"github.com/bernardinorafael/internal/_shared/util"
	"github.com/bernardinorafael/internal/modules/user"
	"github.com/bernardinorafael/pkg/transaction"
	"github.com/jmoiron/sqlx"
)

func (r repo) Create(ctx context.Context, usr user.Entity) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	err := transaction.ExecTx(ctx, r.db, func(tx *sqlx.Tx) error {
		_, err := tx.NamedExecContext(ctx, user.InsertUserQuery, usr)
		if err != nil {
			return err
		}

		emailId := util.GenID("email")
		_, err = tx.ExecContext(ctx, user.InsertEmailQuery, emailId, usr.ID, usr.EmailAddress)
		if err != nil {
			return err
		}

		phoneId := util.GenID("phone")
		_, err = tx.ExecContext(ctx, user.InsertPhoneQuery, phoneId, usr.ID, usr.PhoneNumber)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
