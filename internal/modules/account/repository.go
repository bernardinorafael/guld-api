package account

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/bernardinorafael/internal/_shared/util"
	"github.com/bernardinorafael/internal/modules/org"
	"github.com/bernardinorafael/internal/modules/user"
	"github.com/bernardinorafael/pkg/transaction"
	"github.com/jmoiron/sqlx"
)

type repo struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) RepositoryInterface {
	return &repo{db}
}

func (r repo) FindByUsername(ctx context.Context, username string) (*EntityWithUser, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var acc Entity
	var user user.Entity
	var organization org.EntityWithSettings

	err := transaction.ExecTx(ctx, r.db, func(tx *sqlx.Tx) error {
		err := tx.GetContext(
			ctx,
			&acc,
			`
			SELECT a.* FROM accounts a
			INNER JOIN users u ON u.id = a.user_id
			WHERE u.username = $1
		`,
			username,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil
			}
			return fmt.Errorf("error on find account: %w", err)
		}

		err = tx.GetContext(ctx, &user, `SELECT * FROM users WHERE id = $1`, acc.UserID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil
			}
			return fmt.Errorf("error on find user: %w", err)
		}

		err = tx.GetContext(
			ctx,
			&organization,
			`SELECT * FROM organizations WHERE owner_id = $1`,
			acc.UserID,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil
			}
			return fmt.Errorf("error on find organization: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error on find account: %w", err)
	}

	return &EntityWithUser{
		ID:       acc.ID,
		Password: acc.Password,
		Org:      &organization,
		IsActive: acc.IsActive,
		Created:  acc.Created,
		Updated:  acc.Updated,
		User:     user,
	}, nil
}

func (r repo) FindByUserID(ctx context.Context, userId string) (*Entity, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var acc Entity
	err := r.db.GetContext(
		ctx,
		&acc,
		"SELECT * FROM accounts WHERE user_id = $1",
		userId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error on find account: %w", err)
	}

	return &acc, nil
}

func (r repo) Update(ctx context.Context, acc Entity) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	_, err := r.db.NamedExecContext(
		ctx,
		`
		UPDATE accounts
		SET
			password = :password,
			is_active = :is_active,
			updated = :updated
		`,
		acc,
	)
	if err != nil {
		return fmt.Errorf("error on update account: %w", err)
	}

	return nil
}

func (r repo) FindByID(ctx context.Context, accountId string) (*EntityWithUser, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var acc Entity
	var user user.Entity
	var organization org.EntityWithSettings

	err := transaction.ExecTx(ctx, r.db, func(tx *sqlx.Tx) error {
		err := tx.GetContext(
			ctx,
			&acc,
			`SELECT id, user_id, is_active, created, updated FROM accounts WHERE id = $1`,
			accountId,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil
			}
			return fmt.Errorf("error on find account: %w", err)
		}

		err = tx.GetContext(ctx, &user, `SELECT * FROM users WHERE id = $1`, acc.UserID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil
			}
			return fmt.Errorf("error on find user: %w", err)
		}

		err = tx.GetContext(
			ctx,
			&organization,
			`
			SELECT
				o.*,
				os.id as "settings.id",
				os.org_id as "settings.org_id",
				os.is_active as "settings.is_active",
				os.default_membership_password as "settings.default_membership_password",
				os.max_allowed_memberships as "settings.max_allowed_memberships",
				os.max_allowed_roles as "settings.max_allowed_roles",
				os.use_master_password as "settings.use_master_password",
				os.created as "settings.created",
				os.updated as "settings.updated"
			FROM organizations o
			INNER JOIN organization_settings os ON os.org_id = o.id
			WHERE o.owner_id = $1
		`,
			acc.UserID,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil
			}
			return fmt.Errorf("error on find organization: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error on find account: %w", err)
	}

	return &EntityWithUser{
		ID:       acc.ID,
		Password: acc.Password,
		Org:      &organization,
		IsActive: acc.IsActive,
		Created:  acc.Created,
		Updated:  acc.Updated,
		User:     user,
	}, nil
}

func (r repo) Insert(ctx context.Context, acc EntityWithUser) error {
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

		_, err := tx.NamedExecContext(ctx, query, acc.User)
		if err != nil {
			return fmt.Errorf("error on insert user: %w", err)
		}

		emailId := util.GenID("email")
		_, err = tx.ExecContext(
			ctx,
			"INSERT INTO emails (id, user_id, email, is_primary, is_verified) VALUES ($1, $2, $3, true, true)",
			emailId,
			acc.User.ID,
			acc.User.EmailAddress,
		)
		if err != nil {
			return fmt.Errorf("error on insert email: %w", err)
		}

		_, err = tx.ExecContext(
			ctx,
			`INSERT INTO accounts (id, user_id, password) VALUES ($1, $2, $3)`,
			acc.ID,
			acc.User.ID,
			acc.Password,
		)
		if err != nil {
			return fmt.Errorf("error on insert account: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("error on insert account: %w", err)
	}

	return nil
}
