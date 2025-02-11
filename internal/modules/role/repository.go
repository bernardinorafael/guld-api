package role

import (
	"context"
	"time"

	"github.com/bernardinorafael/pkg/transaction"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type repo struct{ db *sqlx.DB }

func NewRepository(db *sqlx.DB) RepositoryInterface {
	return &repo{db}
}

func (r *repo) BatchRolePermissions(ctx context.Context, roleId string, permissions []string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	err := transaction.ExecTx(ctx, r.db, func(tx *sqlx.Tx) error {
		_, err := tx.ExecContext(ctx,
			`DELETE FROM role_permissions WHERE role_id = $1 AND permission_id NOT IN ($2)`,
			roleId,
			pq.Array(permissions),
		)
		if err != nil {
			return err
		}

		var perms = make([]RolePermissionBatch, 0)
		for _, permissionId := range permissions {
			perms = append(perms, RolePermissionBatch{
				RoleID:       roleId,
				PermissionID: permissionId,
			})
		}

		_, err = tx.NamedExecContext(
			ctx,
			`
			INSERT INTO role_permissions (
				role_id,
				permission_id
			) VALUES (
				:role_id,
				:permission_id
			)
			ON CONFLICT (role_id, permission_id) DO NOTHING
			`,
			perms,
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

func (r *repo) FindByID(ctx context.Context, teamId string, roleId string) (*Entity, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var role Entity
	role.Permissions = []Permission{}

	sql := `
		SELECT
			r.id,
			r.name,
			r.team_id,
			r.key,
			r.description,
			r.created,
			r.updated,
			p.id,
			p.name,
			p.key
		FROM roles r
		LEFT JOIN role_permissions rp ON r.id = rp.role_id
		LEFT JOIN permissions p ON rp.permission_id = p.id
		WHERE r.id = $1 AND r.team_id = $2;
	`

	rows, err := r.db.QueryContext(ctx, sql, roleId, teamId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p Permission
		var permId *string
		var permName *string
		var permKey *string

		err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.TeamID,
			&role.Key,
			&role.Description,
			&role.Created,
			&role.Updated,
			&permId,
			&permName,
			&permKey,
		)
		if err != nil {
			return nil, err
		}

		if permId != nil {
			p.ID = *permId
			if permName != nil {
				p.Name = *permName
			}
			if permKey != nil {
				p.Key = *permKey
			}
			role.Permissions = append(role.Permissions, p)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &role, nil
}

func (r *repo) Create(ctx context.Context, role Entity) (*Entity, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var roleId string
	err := r.db.GetContext(
		ctx,
		&roleId,
		`
		INSERT INTO roles (
			team_id,
		 	name,
			key,
			description
		) VALUES ($1, $2, $3, $4)
		RETURNING id
		`,
		role.TeamID,
		role.Name,
		role.Key,
		role.Description,
	)
	if err != nil {
		return nil, err
	}
	role.ID = roleId

	return &role, nil
}

func (r *repo) Delete(ctx context.Context, teamId, roleId string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(
		ctx,
		"DELETE FROM roles WHERE id = $1 AND team_id = $2",
		roleId,
		teamId,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *repo) GetAll(ctx context.Context, teamId string) ([]Entity, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var roles = make([]Entity, 0)

	sql := `
		SELECT
			r.id,
			r.name,
			r.team_id,
			r.key,
			r.description,
			r.created,
			r.updated,
			p.id AS perm_id,
			p.name AS perm_name,
			p.key AS perm_key
		FROM roles r
		LEFT JOIN role_permissions rp ON r.id = rp.role_id
		LEFT JOIN permissions p ON rp.permission_id = p.id
		WHERE r.team_id = $1 ORDER BY r.created ASC, r.id DESC;
	`

	rows, err := r.db.QueryContext(ctx, sql, teamId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roleMap := make(map[string]*Entity)

	for rows.Next() {
		var role Entity
		var permId *string
		var permName *string
		var permKey *string

		err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.TeamID,
			&role.Key,
			&role.Description,
			&role.Created,
			&role.Updated,
			&permId,
			&permName,
			&permKey,
		)
		if err != nil {
			return nil, err
		}

		if _, exists := roleMap[role.ID]; !exists {
			role.Permissions = []Permission{}
			roleMap[role.ID] = &role
		}

		if permId != nil {
			var p Permission
			p.ID = *permId
			if permName != nil {
				p.Name = *permName
			}
			if permKey != nil {
				p.Key = *permKey
			}
			roleMap[role.ID].Permissions = append(roleMap[role.ID].Permissions, p)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	for _, role := range roleMap {
		roles = append(roles, *role)
	}

	return roles, nil
}
