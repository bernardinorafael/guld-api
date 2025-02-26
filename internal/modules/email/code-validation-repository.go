package email

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func (r repo) FindAllCodesByUser(ctx context.Context, userId string) ([]ValidationEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var entities []ValidationEntity
	err := r.db.SelectContext(
		ctx,
		&entities,
		"SELECT * FROM email_validations WHERE user_id = $1",
		userId,
	)
	if err != nil {
		return nil, fmt.Errorf("error on find all codes by user: %w", err)
	}

	return entities, nil
}

func (r repo) FindCodeValidationByEmailId(ctx context.Context, emailId string) (*ValidationEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var entity ValidationEntity
	err := r.db.GetContext(
		ctx,
		&entity,
		"SELECT * FROM email_validations WHERE email_id = $1 AND is_consumed = false AND is_valid = true",
		emailId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error on find validation by email id: %w", err)
	}

	return &entity, nil
}

func (r repo) UpdateCodeValidation(ctx context.Context, entity ValidationEntity) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var query = `
		UPDATE email_validations
		SET
			attempts = :attempts,
			is_consumed = :is_consumed,
			is_valid = :is_valid
		WHERE id = :id
	`

	if _, err := r.db.NamedExecContext(ctx, query, entity); err != nil {
		return fmt.Errorf("failed to update email validation: %w", err)
	}

	return nil
}

func (r repo) InsertCodeValidation(ctx context.Context, entity ValidationEntity) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	sql := `
		INSERT INTO email_validations (
			id,
			code,
			email_id,
			user_id,
			is_consumed,
			attempts,
			created,
			expires
		) VALUES (
			:id,
			:code,
			:email_id,
			:user_id,
			:is_consumed,
			:attempts,
			:created,
			:expires
		)
	`

	if _, err := r.db.NamedExecContext(ctx, sql, entity); err != nil {
		return fmt.Errorf("failed to insert email validation: %w", err)
	}

	return nil
}
