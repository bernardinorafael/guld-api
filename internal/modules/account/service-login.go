package account

import (
	"context"
	"time"

	. "github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/internal/infra/token"
	"github.com/bernardinorafael/internal/modules/account/session"
	"github.com/bernardinorafael/pkg/crypto"
)

var (
	errInvalidCredential = NewConflictError("invalid credentials", InvalidCredentials, nil, nil)
)

func (s svc) Login(ctx context.Context, username string, password string) (*AccountPayload, error) {
	account, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return nil, errInvalidCredential
	}
	user := account.User
	// Check if password is correct
	if !crypto.PasswordMatches(password, account.Password) {
		return nil, errInvalidCredential
	}
	// Check if account is active
	if !account.IsActive {
		return nil, NewBadRequestError("account is not active", nil)
	}

	sessions, err := s.sessionRepo.FindAllByUsername(ctx, username)
	if err != nil {
		return nil, NewBadRequestError("error on retrieve all sessions by username", err)
	}

	// TODO: When the maximum number of sessions is reached
	// it should log out of one session and continue the login
	if len(sessions) == 3 {
		return nil, NewConflictError("max sessions reached", MaxSessionsReached, nil, nil)
	}

	// Access token with 15 minutes expiration
	accessToken, accessClaims, err := token.Generate(s.secretKey, account.ID, user.ID, user.Username, time.Second*15)
	if err != nil {
		return nil, NewBadRequestError("error on generate access token", err)
	}
	// Refresh token with 30 days expiration
	refreshToken, refreshClaims, err := token.Generate(s.secretKey, account.ID, user.ID, user.Username, time.Hour*24*30)
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
		AccessTokenExpires:  accessClaims.RegisteredClaims.ExpiresAt.Unix(),
		RefreshTokenExpires: refreshClaims.RegisteredClaims.ExpiresAt.Unix(),
	}

	return &payload, nil
}
