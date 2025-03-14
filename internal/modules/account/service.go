package account

import (
	"context"
	"fmt"
	"time"

	"github.com/bernardinorafael/internal/_shared/dto"
	. "github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/internal/infra/http/middleware"
	"github.com/bernardinorafael/internal/infra/token"
	"github.com/bernardinorafael/internal/mailer"
	"github.com/bernardinorafael/internal/modules/account/session"
	"github.com/bernardinorafael/internal/modules/user"
	"github.com/bernardinorafael/pkg/crypto"
	"github.com/bernardinorafael/pkg/logger"
	"github.com/medama-io/go-useragent"
)

var (
	errInvalidCredential = NewConflictError("invalid credentials", InvalidCredentials, nil, nil)
)

type svc struct {
	ctx         context.Context
	log         logger.Logger
	repo        RepositoryInterface
	userRepo    user.RepositoryInterface
	sessionRepo session.RepositoryInterface
	mailer      mailer.Mailer
	secretKey   string
}

func NewService(
	ctx context.Context,
	log logger.Logger,
	repo RepositoryInterface,
	userRepo user.RepositoryInterface,
	sessionRepo session.RepositoryInterface,
	mailer mailer.Mailer,
	secretKey string,
) ServiceInterface {
	return &svc{
		ctx:         ctx,
		log:         log,
		repo:        repo,
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		mailer:      mailer,
		secretKey:   secretKey,
	}
}

func (s svc) GetSession(ctx context.Context) (*dto.SessionResponse, error) {
	claims, ok := ctx.Value(middleware.AuthKey{}).(*token.AccountClaims)
	if !ok {
		return nil, NewBadRequestError("user ID not found in context", nil)
	}

	records, err := s.sessionRepo.FindAllByUsername(ctx, claims.Username)
	if err != nil {
		return nil, NewBadRequestError("error retrieving sessions by username", err)
	}

	if len(records) == 0 {
		s.log.Info(ctx, "no sessions found for user")
		return nil, nil
	}

	var session *session.Entity
	for _, record := range records {
		if !record.Revoked {
			session = &record
			break
		}
	}

	if session == nil {
		s.log.Info(ctx, "no active session found")
		return nil, nil
	}

	res := &dto.SessionResponse{
		ID:      session.ID,
		Agent:   session.Agent,
		IP:      session.IP,
		Revoked: session.Revoked,
		Expired: time.Now().After(session.Expires),
		Created: session.Created,
	}

	return res, nil
}

func (s svc) RevokeSession(ctx context.Context, username, sessionId string) error {
	record, err := s.sessionRepo.FindByID(ctx, sessionId)
	if err != nil {
		return NewBadRequestError("error on get session by id", err)
	}
	if record == nil {
		return NewNotFoundError("session not found", nil)
	}

	session := session.NewFromDatabase(*record)
	session.Revoke()
	sessionData := session.Store()

	err = s.sessionRepo.Update(ctx, sessionData)
	if err != nil {
		return NewBadRequestError("error on update session", err)
	}

	return nil
}

func (s svc) GetAllSessions(ctx context.Context, username string) ([]*dto.SessionResponse, error) {
	records, err := s.sessionRepo.FindAllByUsername(ctx, username)
	if err != nil {
		return nil, NewBadRequestError("error on get all sessions by username", err)
	}

	sessions := make([]*dto.SessionResponse, 0, len(records))

	for _, s := range records {
		sessions = append(sessions, &dto.SessionResponse{
			ID:      s.ID,
			Agent:   s.Agent,
			IP:      s.IP,
			Revoked: s.Revoked,
			Expired: time.Now().After(s.Expires),
			Created: s.Created,
		})
	}

	return sessions, nil
}

func (s svc) ChangePassword(ctx context.Context, userId string, oldPassword string, newPassword string) error {
	user, err := s.userRepo.FindByID(ctx, userId)
	if err != nil {
		return NewBadRequestError("error on get user by id", err)
	}
	if user == nil {
		return NewNotFoundError("user not found", nil)
	}

	account, err := s.repo.FindByUserID(ctx, userId)
	if err != nil {
		return NewBadRequestError("error on get account by id", err)
	}
	if account == nil {
		return NewNotFoundError("account not found", nil)
	}

	newAcc, err := NewFromDatabase(*account)
	if err != nil {
		return NewBadRequestError("error on create account entity", nil)
	}

	if !crypto.PasswordMatches(oldPassword, newAcc.password) {
		return NewConflictError("passwords does not matches", InvalidCredentials, nil, nil)
	}

	hashedPassword, err := crypto.HashPassword(newPassword)
	if err != nil {
		return NewBadRequestError("failed to encrypt password", nil)
	}

	err = newAcc.ChangePassword(hashedPassword, user.IgnorePasswordPolicy)
	if err != nil {
		return NewBadRequestError("error on change password", err)
	}

	accountData := newAcc.Store()

	err = s.repo.Update(ctx, accountData)
	if err != nil {
		return NewBadRequestError("error on updating account password", err)
	}

	go func() {
		params := mailer.SendParams{
			From:    mailer.NotificationSender,
			To:      "rafaelferreirab2@gmail.com",
			Subject: "Sua senha foi alterada",
			File:    "change_password.html",
		}
		if err := s.mailer.Send(params); err != nil {
			s.log.Errorw(ctx, "error on send email", logger.Err(err))
		}
	}()

	return nil
}

func (s svc) Login(ctx context.Context, username, password, userAgent, ip string) (*dto.AccountResponse, error) {
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

	// sessions, err := s.sessionRepo.FindAllByUsername(ctx, username)
	// if err != nil {
	// 	return nil, NewBadRequestError("error on retrieve all sessions by username", err)
	// }

	// TODO: When the maximum number of sessions is reached
	// it should log out of one session and continue the login
	// if len(sessions) == 3 {
	// 	return nil, NewConflictError("max sessions reached", MaxSessionsReached, nil, nil)
	// }

	// Access token with 15 minutes expiration
	accessToken, accessClaims, err := token.Generate(s.secretKey, account.ID, user.ID, user.Username, time.Minute*15)
	if err != nil {
		return nil, NewBadRequestError("error on generate access token", err)
	}
	// Refresh token with 30 days expiration
	refreshToken, refreshClaims, err := token.Generate(s.secretKey, account.ID, user.ID, user.Username, time.Hour*24*30)
	if err != nil {
		return nil, NewBadRequestError("error on generate refresh token", err)
	}

	agent := useragent.NewParser().Parse(userAgent)
	userAgent = fmt.Sprintf("%s em %s", agent.GetBrowser(), agent.GetOS())

	newSession := session.New(user.Username, refreshToken, userAgent, ip)
	sessionData := newSession.Store()

	err = s.sessionRepo.Insert(ctx, sessionData)
	if err != nil {
		return nil, NewBadRequestError("error on insert session", err)
	}

	payload := dto.AccountResponse{
		SessionID:           newSession.ID(),
		AccessToken:         accessToken,
		RefreshToken:        refreshToken,
		AccessTokenExpires:  accessClaims.RegisteredClaims.ExpiresAt.Unix(),
		RefreshTokenExpires: refreshClaims.RegisteredClaims.ExpiresAt.Unix(),
	}

	return &payload, nil
}

func (s svc) Logout(ctx context.Context) error {
	claims, ok := ctx.Value(middleware.AuthKey{}).(*token.AccountClaims)
	if !ok {
		return NewBadRequestError("user ID not found in context", nil)
	}

	sessions, err := s.sessionRepo.FindAllByUsername(ctx, claims.Username)
	if err != nil {
		return NewBadRequestError("error on find all sessions by username", err)
	}

	var sess *session.Entity
	for _, v := range sessions {
		if !v.Revoked {
			sess = &v
			break
		}
	}
	if sess == nil {
		s.log.Info(ctx, "no revoked session found")
		return nil
	}

	err = s.RevokeSession(ctx, claims.Username, sess.ID)
	if err != nil {
		return NewBadRequestError("error on revoke session", err)
	}

	return nil
}

func (s svc) RenewAccessToken(ctx context.Context, refreshToken string) (*dto.RenewAccessToken, error) {
	refreshTokenClaims, err := token.Verify(s.secretKey, refreshToken)
	if err != nil {
		return nil, NewBadRequestError("error on verify refresh token", err)
	}

	acc, err := s.repo.FindByID(ctx, refreshTokenClaims.AccountID)
	if err != nil {
		return nil, NewBadRequestError("error on find account by id", err)
	}
	user := acc.User

	record, err := s.sessionRepo.FindByRefreshToken(ctx, refreshToken)
	if err != nil {
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

	payload := dto.RenewAccessToken{
		AccessToken:        accessToken,
		AccessTokenExpires: claims.ExpiresAt.Unix(),
	}

	return &payload, nil
}

func (s svc) GetSignedInAccount(ctx context.Context) (*EntityWithUser, error) {
	claims, ok := ctx.Value(middleware.AuthKey{}).(*token.AccountClaims)
	if !ok {
		return nil, NewBadRequestError("user ID not found in context", nil)
	}

	acc, err := s.repo.FindByID(ctx, claims.AccountID)
	if err != nil {
		return nil, NewBadRequestError("error on get account by id", err)
	}
	if acc == nil {
		return nil, NewNotFoundError("account not found", nil)
	}

	return acc, nil
}
