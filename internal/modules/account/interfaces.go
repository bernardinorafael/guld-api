package account

import (
	"context"

	"github.com/bernardinorafael/internal/_shared/dto"
)

type RepositoryInterface interface {
	Insert(ctx context.Context, acc EntityWithUser) error
	FindByID(ctx context.Context, accountId string) (*EntityWithUser, error)
	FindByUserID(ctx context.Context, userId string) (*Entity, error)
	FindByUsername(ctx context.Context, username string) (*EntityWithUser, error)
	Update(ctx context.Context, acc Entity) error
}

type ServiceInterface interface {
	Login(ctx context.Context, username, password, userAgent, ip string) (*dto.AccountResponse, error)
	Logout(ctx context.Context, username string) error
	RenewAccessToken(ctx context.Context, refreshToken string) (*dto.RenewAccessToken, error)
	GetSignedInAccount(ctx context.Context) (*EntityWithUser, error)
	ChangePassword(ctx context.Context, userId string, old string, new string) error
	GetAllSessions(ctx context.Context, username string) ([]*dto.SessionResponse, error)
	RevokeSession(ctx context.Context, username, sessionId string) error
}
