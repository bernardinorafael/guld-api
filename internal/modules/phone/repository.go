package phone

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

func (r repo) Create(ctx context.Context, phone AdditionalPhone) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	query := `
		INSERT INTO phones (
			user_id,
			phone,
			primary,
			verified
		) VALUES (
			:user_id,
			:phone,
			:primary,
			:verified
		)
	`
	_, err := r.db.NamedExecContext(ctx, query, phone)
	if err != nil {
		return err
	}

	return nil
}

func (r repo) Delete(ctx context.Context, userId, phoneId string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(
		ctx,
		"DELETE FROM phones WHERE id = $1 AND user_id = $2",
		phoneId,
		userId,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r repo) GetAllByUser(ctx context.Context, userId string) ([]AdditionalPhone, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var phones []AdditionalPhone
	err := r.db.SelectContext(ctx, &phones, "SELECT * FROM phones WHERE user_id = $1", userId)
	if err != nil {
		return phones, err
	}

	return phones, nil
}

func (r repo) FindByID(ctx context.Context, phoneId string) (*AdditionalPhone, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var p AdditionalPhone
	err := r.db.GetContext(ctx, &p, "SELECT * FROM phones WHERE id = $1", phoneId)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r repo) Update(ctx context.Context, phone PhoneUpdateParams) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// The ID field must always be present
	params := map[string]any{"id": phone.ID}
	clauses := []string{}

	if phone.Phone != nil {
		clauses = append(clauses, "phone = :phone")
		params["phone"] = phone.Phone
	}
	if phone.IsPrimary != nil {
		clauses = append(clauses, "primary = :primary")
		params["primary"] = phone.IsPrimary
	}
	if phone.IsVerified != nil {
		clauses = append(clauses, "verified = :verified")
		params["verified"] = phone.IsVerified
	}

	_, err := r.db.NamedExecContext(
		ctx,
		fmt.Sprintf(
			`UPDATE phones SET %s, updated = now() WHERE id = :id`,
			strings.Join(clauses, ", "),
		),
		params,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) FindByPhone(ctx context.Context, phone string) (*AdditionalPhone, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var p AdditionalPhone
	err := r.db.GetContext(ctx, &p, "SELECT * FROM phones WHERE phone = $1", phone)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *repo) GetPrimary(ctx context.Context, userId string) (*AdditionalPhone, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var p AdditionalPhone
	err := r.db.GetContext(
		ctx,
		&p,
		`SELECT * FROM phones WHERE user_id = $1 AND primary = true`,
		userId,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve additional_phone: %w", err)
	}

	return &p, nil
}
