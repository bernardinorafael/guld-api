package usersvc

import (
	"context"

	"github.com/bernardinorafael/internal/modules/email"
)

func (s svc) FindAllEmails(ctx context.Context, userId string) ([]email.Entity, error) {
	emails, err := s.emailService.FindAllByUser(ctx, userId)
	if err != nil {
		return nil, err
	}

	return emails, nil
}
