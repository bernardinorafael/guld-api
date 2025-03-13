package session

import "context"

type RepositoryInterface interface {
	Insert(ctx context.Context, entity Entity) error
	FindByID(ctx context.Context, sessionId string) (*Entity, error)
	FindByRefreshToken(ctx context.Context, refreshToken string) (*Entity, error)
	Update(ctx context.Context, entity Entity) error
	Delete(ctx context.Context, sessionId string) error
	DeleteAll(ctx context.Context, username string) error
	FindAllByUsername(ctx context.Context, username string) ([]Entity, error)
}
