package email

import (
	"context"
)

type RepositoryInterface interface {
	Insert(ctx context.Context, entity Entity) error
	Update(ctx context.Context, entity Entity) error
	FindByID(ctx context.Context, emailId string) (*Entity, error)
	FindByEmail(ctx context.Context, entity string) (*Entity, error)
	FindAllByUser(ctx context.Context, userId string) ([]Entity, error)
	Delete(ctx context.Context, userId, emailId string) error

	InsertValidation(ctx context.Context, v Validation) error
	UpdateValidation(ctx context.Context, v Validation) error
	FindValidationByEmail(ctx context.Context, emailId string) (*Validation, error)
}

type ServiceInterface interface {
	Create(ctx context.Context, dto CreateParams) (*Entity, error)
	FindAllByUser(ctx context.Context, userId string) ([]Entity, error)
	FindByEmail(ctx context.Context, email string) (*Entity, error)
	FindByID(ctx context.Context, emailId string) (*Entity, error)
	Update(ctx context.Context, entity Entity) (*Entity, error)
	Delete(ctx context.Context, userId, emailId string) error

	InsertValidation(ctx context.Context, v Validation) error
	UpdateValidation(ctx context.Context, v Validation) error
	FindValidationByEmail(ctx context.Context, emailId string) (*Validation, error)
}
