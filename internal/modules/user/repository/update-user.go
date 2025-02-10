package userrepo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bernardinorafael/internal/modules/user"
)

func (r repo) Update(ctx context.Context, input user.PartialEntity) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	params := map[string]any{"id": input.ID}
	clauses := []string{}

	if input.FullName != nil {
		clauses = append(clauses, "full_name = :full_name")
		params["full_name"] = input.FullName
	}

	if input.Username != nil {
		clauses = append(clauses, "username = :username")
		params["username"] = input.Username
	}

	if input.EmailAddress != nil {
		clauses = append(clauses, "email_address = :email_address")
		params["email_address"] = input.EmailAddress
	}

	if input.PhoneNumber != nil {
		clauses = append(clauses, "phone_number = :phone_number")
		params["phone_number"] = input.PhoneNumber
	}

	if input.AvatarURL != nil {
		clauses = append(clauses, "avatar_url = :avatar_url")
		params["avatar_url"] = input.AvatarURL
	}

	if input.Banned != nil {
		clauses = append(clauses, "banned = :banned")
		params["banned"] = input.Banned
	}

	if input.Locked != nil {
		clauses = append(clauses, "locked = :locked")
		params["locked"] = input.Locked
	}

	sql := fmt.Sprintf(
		`UPDATE users SET %s, updated = now() WHERE id = :id`,
		strings.Join(clauses, ", "),
	)

	_, err := r.db.NamedExecContext(ctx, sql, params)
	if err != nil {
		return err
	}

	return nil
}
