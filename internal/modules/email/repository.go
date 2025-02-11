package email

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type repo struct{ db *sqlx.DB }

func NewRepository(db *sqlx.DB) RepositoryInterface {
	return &repo{db}
}

func (r *repo) GetPrimary(ctx context.Context, userId string) (*AdditionalEmail, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var addr AdditionalEmail
	err := r.db.GetContext(
		ctx,
		&addr,
		`SELECT * FROM emails WHERE user_id = $1 AND primary = true`,
		userId,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve additional_email: %w", err)
	}

	return &addr, nil
}

func (r *repo) FindByEmail(ctx context.Context, email string) (*AdditionalEmail, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var addr AdditionalEmail
	err := r.db.GetContext(ctx, &addr, "SELECT * FROM emails WHERE email = $1", email)
	if err != nil {
		return nil, err
	}

	return &addr, nil
}

func (r repo) Create(ctx context.Context, email AdditionalEmail) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	query := `
		INSERT INTO emails (
			user_id,
			email,
			primary,
			verified
		) VALUES (
			:user_id,
			:email,
			:primary,
			:verified
		)
	`
	_, err := r.db.NamedExecContext(ctx, query, email)
	if err != nil {
		return err
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
		return err
	}

	return nil
}

func (r repo) GetAllByUser(ctx context.Context, userId string) ([]AdditionalEmail, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var emails []AdditionalEmail
	err := r.db.SelectContext(ctx, &emails, "SELECT * FROM emails WHERE user_id = $1", userId)
	if err != nil {
		return emails, err
	}

	return emails, nil
}

func (r repo) FindByID(ctx context.Context, emailId string) (*AdditionalEmail, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var email AdditionalEmail
	err := r.db.GetContext(ctx, &email, "SELECT * FROM emails WHERE id = $1", emailId)
	if err != nil {
		return nil, err
	}

	return &email, nil
}

func (r repo) Update(ctx context.Context, email EmailUpdateParams) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// The ID field must always be present
	params := map[string]any{"id": email.ID}
	clauses := []string{}

	if email.Email != nil {
		clauses = append(clauses, "email = :email")
		params["email"] = email.Email
	}

	if email.IsPrimary != nil {
		clauses = append(clauses, "primary = :primary")
		params["primary"] = email.IsPrimary
	}

	if email.IsVerified != nil {
		clauses = append(clauses, "verified = :verified")
		params["verified"] = email.IsVerified
	}

	_, err := r.db.NamedExecContext(
		ctx,
		fmt.Sprintf(
			`UPDATE emails SET %s, updated = now() WHERE id = :id`,
			strings.Join(clauses, ", "),
		),
		params,
	)

	if err != nil {
		return err
	}

	return nil
}
