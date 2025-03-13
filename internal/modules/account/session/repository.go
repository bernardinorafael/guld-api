package session

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) RepositoryInterface {
	return &repo{db: db}
}

func (r repo) FindByRefreshToken(ctx context.Context, refreshToken string) (*Entity, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var session Entity
	err := r.db.GetContext(
		ctx,
		&session,
		"SELECT * FROM sessions WHERE refresh_token = $1",
		refreshToken,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error on find session by refresh token: %w", err)
	}

	return &session, nil
}

func (r repo) FindByID(ctx context.Context, sessionId string) (*Entity, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var session Entity
	err := r.db.GetContext(
		ctx,
		&session,
		"SELECT * FROM sessions WHERE id = $1",
		sessionId,
	)
	if err != nil {
		return nil, fmt.Errorf("error on find session by id: %w", err)
	}

	return &session, nil
}

func (r repo) Insert(ctx context.Context, entity Entity) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	_, err := r.db.NamedExecContext(
		ctx,
		`
			INSERT INTO sessions (
				id,
				username,
				refresh_token,
				agent,
				ip,
				revoked,
				expires,
				created,
				updated
			)	VALUES (
				:id,
				:username,
				:refresh_token,
				:agent,
				:ip,
				:revoked,
				:expires,
				:created,
				:updated
			)
		`,
		entity,
	)
	if err != nil {
		return fmt.Errorf("error on insert session: %w", err)
	}

	return nil
}

func (r repo) Delete(ctx context.Context, sessionId string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, "DELETE FROM sessions where id = $1", sessionId)
	if err != nil {
		return fmt.Errorf("error on delete session: %w", err)
	}

	return nil
}

func (r repo) DeleteAll(ctx context.Context, username string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, "DELETE FROM sessions WHERE username = $1", username)
	if err != nil {
		return fmt.Errorf("error on delete all sessions: %w", err)
	}

	return nil
}

func (r repo) Update(ctx context.Context, entity Entity) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	_, err := r.db.NamedExecContext(
		ctx,
		`
			UPDATE sessions
			SET
				ip = :ip,
				agent = :agent,
				refresh_token = :refresh_token,
				revoked = :revoked,
				expires = :expires,
				updated = :updated
			WHERE id = :id
		`,
		entity,
	)
	if err != nil {
		return fmt.Errorf("error on update session: %w", err)
	}

	return nil
}

func (r repo) FindAllByUsername(ctx context.Context, username string) ([]Entity, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var sessions = []Entity{}
	err := r.db.SelectContext(ctx, &sessions, "SELECT * FROM sessions WHERE username = $1", username)
	if err != nil {
		return nil, fmt.Errorf("error on find sessions by username: %w", err)
	}

	return sessions, nil
}
