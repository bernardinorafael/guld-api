package usersvc

import (
	"context"

	"github.com/bernardinorafael/internal/modules/email"
)

func (s svc) FindEmail(ctx context.Context, email string) (*email.Entity, error) {
	foundEmail, err := s.emailService.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return foundEmail, nil
}
