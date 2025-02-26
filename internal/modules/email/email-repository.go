package email

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type repo struct{ db *sqlx.DB }

func NewRepository(db *sqlx.DB) RepositoryInterface {
	return &repo{db}
}

func (r *repo) FindByEmail(ctx context.Context, email string) (*Entity, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var entity Entity
	err := r.db.GetContext(ctx, &entity, "SELECT * FROM emails WHERE email = $1", email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find email by email: %w", err)
	}

	return &entity, nil
}

func (r repo) Insert(ctx context.Context, email Entity) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var query = `
		INSERT INTO emails (
			id,
			user_id,
			email,
			is_primary,
			is_verified,
			created,
			updated
		) VALUES (
			:id,
			:user_id,
			:email,
			:is_primary,
			:is_verified,
			:created,
			:updated
		)
	`
	_, err := r.db.NamedExecContext(ctx, query, email)
	if err != nil {
		return fmt.Errorf("failed to insert email: %w", err)
	}

	return nil
}

func (r repo) Delete(ctx context.Context, userId, emailId string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(
		ctx,
		"DELETE FROM emails WHERE id = $1 AND user_id = $2",
		emailId,
		userId,
	)
	if err != nil {
		return fmt.Errorf("failed to delete email: %w", err)
	}

	return nil
}

func (r repo) FindAllByUser(ctx context.Context, userId string) ([]Entity, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var emails []Entity
	err := r.db.SelectContext(
		ctx,
		&emails,
		"SELECT * FROM emails WHERE user_id = $1 ORDER BY created DESC",
		userId,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find all emails by user: %w", err)
	}

	return emails, nil
}

func (r repo) FindByID(ctx context.Context, emailId string) (*Entity, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var email Entity
	err := r.db.GetContext(ctx, &email, "SELECT * FROM emails WHERE id = $1", emailId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &email, nil
}

func (r repo) Update(ctx context.Context, entity Entity) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	entity.Updated = time.Now()
	var query = `
		UPDATE emails
		SET
			is_primary = :is_primary,
			is_verified = :is_verified,
			updated = :updated
		WHERE id = :id
	`
	if _, err := r.db.NamedExecContext(ctx, query, entity); err != nil {
		return fmt.Errorf("failed to update email: %w", err)
	}

	return nil
}
