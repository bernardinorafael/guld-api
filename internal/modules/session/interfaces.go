package session

import "context"

type RepositoryInterface interface {
	Insert(ctx context.Context, entity Entity) error
	FindByID(ctx context.Context, username, sessionId string) (*Entity, error)
	Update(ctx context.Context, entity Entity) error
	Delete(ctx context.Context, sessionId string) error
	FindAllByUsername(ctx context.Context, username string) ([]Entity, error)
}
