package org

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/bernardinorafael/internal/modules/user"
	"github.com/bernardinorafael/pkg/transaction"
	"github.com/jmoiron/sqlx"
)

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) RepositoryInterface {
	return &repo{db: db}
}

func (r repo) FindByID(ctx context.Context, id string) (*EntityWithOwner, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var org Entity
	var owner user.Entity

	err := transaction.ExecTx(ctx, r.db, func(tx *sqlx.Tx) error {
		err := tx.GetContext(ctx, &org, `SELECT * FROM organizations WHERE id = $1`, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("organization with id %s not found", id)
			}
			return err
		}

		err = tx.GetContext(ctx, &owner, `SELECT * FROM users WHERE id = $1`, org.OwnerID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("user with id %s not found", org.OwnerID)
			}
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &EntityWithOwner{
		ID:        org.ID,
		Name:      org.Name,
		Slug:      org.Slug,
		Owner:     owner,
		AvatarURL: org.AvatarURL,
		Created:   org.Created,
		Updated:   org.Updated,
	}, nil
}

func (r repo) FindBySlug(ctx context.Context, slug string) (*EntityWithOwner, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var org Entity
	var owner user.Entity

	err := transaction.ExecTx(ctx, r.db, func(tx *sqlx.Tx) error {
		err := tx.GetContext(ctx, &org, `SELECT * FROM organizations WHERE slug = $1`, slug)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("organization with slug %s not found", slug)
			}
			return err
		}

		err = tx.GetContext(ctx, &owner, `SELECT * FROM users WHERE id = $1`, org.OwnerID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("user with id %s not found", org.OwnerID)
			}
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &EntityWithOwner{
		ID:        org.ID,
		Name:      org.Name,
		Slug:      org.Slug,
		Owner:     owner,
		AvatarURL: org.AvatarURL,
		Created:   org.Created,
		Updated:   org.Updated,
	}, nil
}

func (r repo) Insert(ctx context.Context, org Entity) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	_, err := r.db.NamedExecContext(
		ctx,
		`
		INSERT INTO organizations (
			id,
			name,
			slug,
			owner_id,
			avatar_url,
			created,
			updated
		) VALUES (
			:id,
			:name,
			:slug,
			:owner_id,
			:avatar_url,
			:created,
			:updated
		)
		`,
		org,
	)
	if err != nil {
		return err
	}

	return nil
}
