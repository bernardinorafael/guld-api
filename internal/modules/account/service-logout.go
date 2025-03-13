package account

import (
	"context"

	. "github.com/bernardinorafael/internal/_shared/errors"
)

func (s svc) Logout(ctx context.Context, username string) error {
	err := s.sessionRepo.DeleteAll(ctx, username)
	if err != nil {
		return NewBadRequestError("error on delete all sessions", err)
	}

	return nil
}
