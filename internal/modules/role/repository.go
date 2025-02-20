package role

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

func (r *repo) Delete(ctx context.Context, orgId, roleId string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, "DELETE FROM roles WHERE id = $1 AND org_id = $2", roleId, orgId)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	return nil
}

func (r *repo) FindAll(ctx context.Context, orgId string) ([]EntityWithPermission, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var roles []EntityWithPermission
	err := r.db.SelectContext(ctx, &roles, "select * from roles where org_id = $1", orgId)
	if err != nil {
		return nil, fmt.Errorf("failed to find all roles: %w", err)
	}

	for i := range roles {
		roles[i].Permissions = make([]Permission, 0)
	}

	return roles, nil
}

func (r *repo) FindByID(ctx context.Context, orgId, roleId string) (*EntityWithPermission, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var role EntityWithPermission
	role.Permissions = make([]Permission, 0)
	q := `
		SELECT
			r.id,
			r.name,
			r.org_id,
			r.key,
			r.description,
			r.created,
			r.updated,
		FROM roles r
		LEFT JOIN role_permissions rp ON r.id = rp.role_id
		LEFT JOIN permissions p ON rp.permission_id = p.id
		WHERE r.id = $1 AND r.org_id = $2
	`

	err := r.db.GetContext(ctx, &role, q, roleId, orgId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("role not found")
		}
		return nil, fmt.Errorf("failed to find role by id: %w", err)
	}

	return &role, nil
}

func (r *repo) Update(ctx context.Context, entity Entity) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	entity.Updated = time.Now()
	_, err := r.db.NamedExecContext(
		ctx,
		`
		UPDATE roles
		SET
			name = :name,
			key = :key,
			description = :description,
			updated = :updated
		WHERE id = :id
		AND org_id = :org_id
	`,
		entity,
	)
	if err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	return nil
}

func (r *repo) Insert(ctx context.Context, entity Entity) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	_, err := r.db.NamedExecContext(
		ctx,
		`
			INSERT INTO roles (
				id,
				name,
				key,
				org_id,
				description,
				created,
				updated
			) VALUES (
				:id,
				:name,
				:key,
				:org_id,
				:description,
				:created,
				:updated
			)
		`,
		entity,
	)
	if err != nil {
		return fmt.Errorf("failed to insert role: %w", err)
	}

	return nil
}

// func (r *repo) BatchRolePermissions(ctx context.Context, roleId string, permissions []string) error {
// 	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
// 	defer cancel()

// 	err := transaction.ExecTx(ctx, r.db, func(tx *sqlx.Tx) error {
// 		_, err := tx.ExecContext(ctx,
// 			`DELETE FROM role_permissions WHERE role_id = $1 AND permission_id NOT IN ($2)`,
// 			roleId,
// 			pq.Array(permissions),
// 		)
// 		if err != nil {
// 			return err
// 		}

// 		var perms = make([]RolePermissionBatch, 0)
// 		for _, permissionId := range permissions {
// 			perms = append(perms, RolePermissionBatch{
// 				RoleID:       roleId,
// 				PermissionID: permissionId,
// 			})
// 		}

// 		_, err = tx.NamedExecContext(
// 			ctx,
// 			`
// 			INSERT INTO role_permissions (
// 				role_id,
// 				permission_id
// 			) VALUES (
// 				:role_id,
// 				:permission_id
// 			)
// 			ON CONFLICT (role_id, permission_id) DO NOTHING
// 			`,
// 			perms,
// 		)
// 		if err != nil {
// 			return err
// 		}

// 		return nil
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
