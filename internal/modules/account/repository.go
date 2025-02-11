package account

import (
	"context"
	"database/sql"
	"errors"
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

func (r repo) FindByUsername(ctx context.Context, username string) (*Entity, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var acc Entity
	err := r.db.GetContext(
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
			return nil, err
		}
		return nil, err
	}

	return &acc, nil
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
			return nil, err
		}
		return nil, err
	}

	return &acc, nil
}

func (r repo) FindByID(ctx context.Context, accountId string) (*EntityWithUser, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var acc Entity
	var user user.Entity
	var organization org.Entity

	err := transaction.ExecTx(ctx, r.db, func(tx *sqlx.Tx) error {
		err := tx.GetContext(
			ctx,
			&acc,
			`SELECT id, user_id, created, updated FROM accounts WHERE id = $1`,
			accountId,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return err
			}
			return err
		}

		err = tx.GetContext(ctx, &user, `SELECT * FROM users WHERE id = $1`, acc.UserID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return err
			}
			return err
		}

		err = tx.GetContext(
			ctx,
			&organization,
			`SELECT * FROM organizations WHERE owner_id = $1`,
			acc.UserID,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return err
			}
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &EntityWithUser{
		ID:       acc.ID,
		Password: acc.Password,
		Org:      organization,
		Created:  acc.Created,
		Updated:  acc.Updated,
		User:     user,
	}, nil
}

func (r repo) Insert(ctx context.Context, acc EntityWithUser) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	err := transaction.ExecTx(ctx, r.db, func(tx *sqlx.Tx) error {
		_, err := tx.NamedExecContext(ctx, user.InsertUserQuery, acc.User)
		if err != nil {
			return err
		}

		emailId := util.GenID("email")
		_, err = tx.ExecContext(ctx, user.InsertEmailQuery, emailId, acc.User.ID, acc.User.EmailAddress)
		if err != nil {
			return err
		}

		phoneId := util.GenID("phone")
		_, err = tx.ExecContext(ctx, user.InsertPhoneQuery, phoneId, acc.User.ID, acc.User.PhoneNumber)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(
			ctx,
			`INSERT INTO accounts (id, user_id, password) VALUES ($1, $2, $3)`,
			acc.ID,
			acc.User.ID,
			acc.Password,
		)
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
