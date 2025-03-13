package account

import (
	"context"
	"time"

	. "github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/internal/infra/token"
	"github.com/bernardinorafael/internal/modules/account/session"
	"github.com/bernardinorafael/pkg/logger"
)

func (s svc) RenewAccessToken(ctx context.Context, refreshToken string) (*RenewAccessTokenPayload, error) {
	refreshTokenClaims, err := token.Verify(s.secretKey, refreshToken)
	if err != nil {
		s.log.Errorw(ctx, "error on verify refresh token", logger.Err(err))
		return nil, NewBadRequestError("error on verify refresh token", err)
	}

	acc, err := s.repo.FindByID(ctx, refreshTokenClaims.AccountID)
	if err != nil {
		s.log.Errorw(ctx, "error on find account by id", logger.Err(err))
		return nil, NewBadRequestError("error on find account by id", err)
	}
	user := acc.User

	record, err := s.sessionRepo.FindByRefreshToken(ctx, refreshToken)
	if err != nil {
		s.log.Errorw(ctx, "error on find session by refresh token", logger.Err(err))
		return nil, NewBadRequestError("error on find session by refresh token", err)
	}

	session := session.NewFromDatabase(*record)

	if !session.IsValid() {
		return nil, NewBadRequestError("session is invalid", nil)
	}

	if session.Username() != user.Username {
		return nil, NewBadRequestError("session username does not match account username", nil)
	}

	accessToken, claims, err := token.Generate(s.secretKey, acc.ID, user.ID, user.Username, time.Minute*15)
	if err != nil {
		s.log.Errorw(ctx, "error on generate access token", logger.Err(err))
		return nil, NewBadRequestError("error on generate access token", err)
	}

	payload := RenewAccessTokenPayload{
		AccessToken:        accessToken,
		AccessTokenExpires: claims.ExpiresAt.Unix(),
	}

	return &payload, nil
}
