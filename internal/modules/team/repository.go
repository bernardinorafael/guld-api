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

func (r *repo) Update(ctx context.Context, team Entity) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	team.Updated = time.Now()
	_, err := r.db.NamedExecContext(
		ctx,
		`
		UPDATE teams
		SET
			name = :name,
			slug = :slug,
			members_count = :members_count,
			updated = :updated
		WHERE
			id = :id
		`,
		team,
	)
	if err != nil {
		return fmt.Errorf("error on update team: %w", err)
	}

	return nil
}

func (r *repo) FindByMember(ctx context.Context, orgId, userId string) (*EntityWithRole, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var team EntityWithRole
	query := `
		SELECT
			t.id,
			t.name,
			t.slug,
			t.owner_id,
			t.org_id,
			t.logo,
			t.members_count,
			t.created,
			t.updated,
			r.id as "role.id",
			r.name as "role.name"
		FROM teams t
			INNER JOIN team_members tm ON tm.team_id = t.id
			INNER JOIN roles r ON r.id = tm.role_id
			WHERE tm.user_id = $1
			AND t.org_id = $2
		LIMIT 1
	`

	err := r.db.GetContext(ctx, &team, query, userId, orgId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error on find team by member: %w", err)
	}

	return &team, nil
}

func (r *repo) InsertMember(ctx context.Context, member TeamMember) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	sql := `
		INSERT INTO team_members (
			id,
			user_id,
			role_id,
			team_id,
			org_id,
			created,
			updated
		) VALUES (
			:id,
			:user_id,
			:role_id,
			:team_id,
			:org_id,
			:created,
			:updated
		)
	`

	_, err := r.db.NamedExecContext(ctx, sql, member)
	if err != nil {
		return fmt.Errorf("error on insert team member: %w", err)
	}

	return nil
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

	var query = `
		INSERT INTO teams (
			id,
			name,
			slug,
			owner_id,
			org_id,
			logo,
			created,
			members_count,
			updated
		)
		VALUES (
			:id,
			:name,
			:slug,
			:owner_id,
			:org_id,
			:logo,
			:created,
			:members_count,
			:updated
		)
	`

	_, err := r.db.NamedExecContext(ctx, query, team)
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

	teams := make([]EntityWithMembers, 0)
	err := r.db.SelectContext(ctx, &teams, "", ownerId, orgId)
	if err != nil {
		return nil, fmt.Errorf("error on find all teams with members: %w", err)
	}

	return teams, nil
}

func (r repo) FindAll(ctx context.Context, orgId, ownerId string) ([]Entity, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	teams := make([]Entity, 0)
	err := r.db.SelectContext(
		ctx,
		&teams,
		"SELECT * FROM teams WHERE owner_id = $1 AND org_id = $2",
		ownerId,
		orgId,
	)
	if err != nil {
		return nil, fmt.Errorf("error on find all teams: %w", err)
	}

	return teams, nil
}
