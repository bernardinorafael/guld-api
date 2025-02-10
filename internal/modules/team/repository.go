package team

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type repo struct{ db *sqlx.DB }

func NewRepository(db *sqlx.DB) RepositoryInterface {
	return &repo{db}
}

func (r repo) Create(ctx context.Context, team Entity) (*Entity, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var teamId string
	rows, err := r.db.NamedQueryContext(
		ctx,
		`
		INSERT INTO teams (name, slug, owner_id)
		VALUES (:name, :slug, :owner_id)
		RETURNING id
		`,
		team,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		if err := rows.Scan(&teamId); err != nil {
			return nil, err
		}
	}
	team.ID = teamId

	return &team, nil
}

func (r repo) Delete(ctx context.Context, ownerId string, teamId string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(
		ctx,
		"DELETE FROM teams WHERE id = $1 AND owner_id = $2",
		teamId,
		ownerId,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r repo) GetAll(ctx context.Context, ownerId string) ([]Entity, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var teams = make([]Entity, 0)
	err := r.db.SelectContext(
		ctx,
		&teams,
		`
		SELECT * FROM teams
		WHERE owner_id = $1 ORDER BY created DESC
		`,
		ownerId,
	)
	if err != nil {
		return nil, err
	}

	return teams, nil
}
