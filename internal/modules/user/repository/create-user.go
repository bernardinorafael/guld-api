package userrepo

import (
	"context"
	"fmt"
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
		var query = `
			INSERT INTO users (
				id,
				full_name,
				username,
				phone_number,
				email_address,
				avatar_url,
				banned,
				locked,
				username_last_updated,
				username_lockout_end,
				created,
				updated
			) VALUES (
				:id,
				:full_name,
				:username,
				:phone_number,
				:email_address,
				:avatar_url,
				:banned,
				:locked,
				:username_last_updated,
				:username_lockout_end,
				:created,
				:updated
			)
		`

		_, err := tx.NamedExecContext(ctx, query, usr)
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		emailId := util.GenID("email")
		_, err = tx.ExecContext(
			ctx,
			"INSERT INTO emails (id, user_id, email, is_primary, is_verified) VALUES ($1, $2, $3, true, true)",
			emailId,
			usr.ID,
			usr.EmailAddress,
		)
		if err != nil {
			return fmt.Errorf("failed to create email: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}
