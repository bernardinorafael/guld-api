package permission

import "context"

type RepositoryInterface interface {
	FindAll(ctx context.Context) ([]Entity, error)
}

type ServiceInterface interface {
	FindAll(ctx context.Context) ([]Entity, error)
}
