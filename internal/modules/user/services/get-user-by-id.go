package usersvc

import (
	"context"

	. "github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/internal/modules/user"
	"github.com/bernardinorafael/pkg/logger"
)

func (s svc) FindByID(ctx context.Context, id string) (*user.CompleteEntity, error) {
	meta := make(map[string]any)

	found, err := s.userRepo.FindCompleteByID(ctx, id)
	if err != nil {
		msg := "failed to retrieve user"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return nil, NewNotFoundError(msg, err)
	}

	// if email != nil {
	// 	meta["email_verify"] = map[string]any{
	// 		"email":   "marilia_expo@gmail.com",
	// 		"expires": time.Now().Add(time.Hour * 24),
	// 	}
	// }

	return &user.CompleteEntity{
		User:   found.User,
		Emails: found.Emails,
		Meta:   meta,
	}, nil
}
