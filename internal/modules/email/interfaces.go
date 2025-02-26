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

	InsertCodeValidation(ctx context.Context, entity ValidationEntity) error
	UpdateCodeValidation(ctx context.Context, entity ValidationEntity) error
	FindCodeValidationByEmailId(ctx context.Context, emailId string) (*ValidationEntity, error)
	FindAllCodesByUser(ctx context.Context, userId string) ([]ValidationEntity, error)
}

type ServiceInterface interface {
	FindAllByUser(ctx context.Context, userId string) ([]Entity, error)
	FindByEmail(ctx context.Context, email string) (*Entity, error)
	FindByID(ctx context.Context, emailId string) (*Entity, error)
	Update(ctx context.Context, entity Entity) (*Entity, error)
	Delete(ctx context.Context, userId, emailId string) error
	AddEmail(ctx context.Context, dto CreateEmailDTO) error

	GenerateValidationCode(ctx context.Context, dto GenerateEmailValidationDTO) error
	ValidateEmail(ctx context.Context, dto ValidateEmailDTO) error
	FindActiveCodesByUser(ctx context.Context, userId string) ([]ValidationEntity, error)
	DeleteEmail(ctx context.Context, userId, emailId string) error
}
