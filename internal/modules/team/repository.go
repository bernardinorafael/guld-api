package team

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

func (r *repo) FindByID(ctx context.Context, orgId, teamId string) (*Entity, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var team Entity
	err := r.db.GetContext(
		ctx,
		&team,
		"SELECT * FROM teams WHERE id = $1 AND org_id = $2",
		teamId,
		orgId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error on find team by id: %w", err)
	}

	return &team, nil
}

func (r *repo) FindBySlug(ctx context.Context, orgId, slug string) (*Entity, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var team Entity
	err := r.db.GetContext(
		ctx,
		&team,
		"SELECT * FROM teams WHERE slug = $1 AND org_id = $2",
		slug,
		orgId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error on find team by slug: %w", err)
	}

	return &team, nil
}

func (r repo) Insert(ctx context.Context, team Entity) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	_, err := r.db.NamedExecContext(ctx, InsertTeamQuery, team)
	if err != nil {
		return fmt.Errorf("error on insert team: %w", err)
	}

	return nil
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
		return fmt.Errorf("error on delete team: %w", err)
	}

	return nil
}

func (r repo) FindAllWithMembers(ctx context.Context, orgId, ownerId string) ([]EntityWithMembers, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	sql := `
		SELECT t.*, json_agg(u.*) as members
		FROM teams t
		LEFT JOIN users u ON u.team_id = t.id
		WHERE t.owner_id = $1 AND t.org_id = $2
		GROUP BY t.id
	`

	teams := make([]EntityWithMembers, 0)
	err := r.db.SelectContext(ctx, &teams, sql, ownerId, orgId)
	if err != nil {
		return nil, fmt.Errorf("error on find all teams with members: %w", err)
	}

	return teams, nil
}

func (r repo) FindAll(ctx context.Context, orgId, ownerId string) ([]Entity, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	teams := make([]Entity, 0)
	err := r.db.SelectContext(ctx, &teams, FindAllTeamsQuery, ownerId, orgId)
	if err != nil {
		return nil, fmt.Errorf("error on find all teams: %w", err)
	}

	return teams, nil
}

func (r repo) FindAllWithOrg(ctx context.Context, orgId, ownerId string) ([]EntityWithOrg, error) {
	panic("unimplemented")
}
