package account

import (
	"context"
	"time"

	. "github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/internal/infra/token"
	"github.com/bernardinorafael/internal/modules/session"
	"github.com/bernardinorafael/pkg/crypto"
)

var (
	errInvalidCredential = NewConflictError("invalid credentials", InvalidCredentials, nil, nil)
)

func (s svc) Login(ctx context.Context, username string, password string) (*AccountPayload, error) {
	acc, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return nil, errInvalidCredential
	}
	if !crypto.PasswordMatches(password, acc.Password) {
		return nil, errInvalidCredential
	}
	user := acc.User

	if !acc.IsActive {
		return nil, NewBadRequestError("account is not active", nil)
	}

	// Access token with 15 minutes expiration
	accessToken, accessClaims, err := token.Generate(s.secretKey, acc.ID, user.ID, user.Username, time.Minute*15)
	if err != nil {
		return nil, NewBadRequestError("error on generate access token", err)
	}

	// Refresh token with 30 days expiration
	refreshToken, refreshClaims, err := token.Generate(s.secretKey, acc.ID, user.ID, user.Username, time.Hour*24*30)
	if err != nil {
		return nil, NewBadRequestError("error on generate refresh token", err)
	}

	// TODO: get agent and ip from context
	newSession := session.New(user.Username, refreshToken, "agent", "ip")
	sessionData := newSession.Store()

	err = s.sessionRepo.Insert(ctx, sessionData)
	if err != nil {
		return nil, NewBadRequestError("error on insert session", err)
	}

	payload := AccountPayload{
		SessionID:           newSession.ID(),
		AccessToken:         accessToken,
		RefreshToken:        refreshToken,
		AccessTokenExpires:  accessClaims.RegisteredClaims.ExpiresAt.Time,
		RefreshTokenExpires: refreshClaims.RegisteredClaims.ExpiresAt.Time,
		User:                user,
	}

	return &payload, nil
}
