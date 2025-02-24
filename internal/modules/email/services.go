package email

import (
	"context"

	"github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/pkg/logger"
)

type svc struct {
	log  logger.Logger
	repo RepositoryInterface
}

func NewService(log logger.Logger, repo RepositoryInterface) ServiceInterface {
	return &svc{log, repo}
}

func (s svc) InsertValidation(ctx context.Context, v Validation) error {
	if err := s.repo.InsertValidation(ctx, v); err != nil {
		s.log.Errorw(ctx, "failed to insert email validation", logger.Err(err))
		return errors.NewBadRequestError("failed to insert email validation", nil)
	}

	return nil
}

func (s svc) UpdateValidation(ctx context.Context, v Validation) error {
	if err := s.repo.UpdateValidation(ctx, v); err != nil {
		s.log.Errorw(ctx, "failed to update email validation", logger.Err(err))
		return errors.NewBadRequestError("failed to update email validation", nil)
	}

	return nil
}

func (s svc) FindValidationByEmail(ctx context.Context, email string) (*Validation, error) {
	validation, err := s.repo.FindValidationByEmail(ctx, email)
	if err != nil {
		s.log.Errorw(ctx, "failed to find email validation", logger.Err(err))
		return nil, errors.NewBadRequestError("failed to find email validation", nil)
	}
	if validation == nil {
		s.log.Errorw(ctx, "email validation not found", logger.Err(err))
		return nil, errors.NewNotFoundError("email validation not found", nil)
	}

	return validation, nil
}

func (s svc) FindValidationByEmailID(ctx context.Context, emailId string) (*Validation, error) {
	validation, err := s.repo.FindValidationByEmailID(ctx, emailId)
	if err != nil {
		s.log.Errorw(ctx, "failed to find email validation", logger.Err(err))
		return nil, errors.NewBadRequestError("failed to find email validation", nil)
	}
	if validation == nil {
		s.log.Errorw(ctx, "email validation not found", logger.Err(err))
		return nil, errors.NewNotFoundError("email validation not found", nil)
	}

	return validation, nil
}

func (s svc) Delete(ctx context.Context, userId, emailId string) error {
	if err := s.repo.Delete(ctx, userId, emailId); err != nil {
		s.log.Errorw(ctx, "failed to delete email", logger.Err(err))
		return errors.NewBadRequestError("failed to delete email", nil)
	}

	return nil
}

func (s svc) FindAllByUser(ctx context.Context, userId string) ([]Entity, error) {
	emails, err := s.repo.FindAllByUser(ctx, userId)
	if err != nil {
		s.log.Errorw(ctx, "failed to find all emails by user", logger.Err(err))
		return nil, errors.NewBadRequestError("failed to find all emails by user", nil)
	}

	return emails, nil
}

func (s svc) FindByEmail(ctx context.Context, email string) (*Entity, error) {
	foundEmail, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		s.log.Errorw(ctx, "failed to find email by email", logger.Err(err))
		return nil, errors.NewBadRequestError("failed to find email by email", nil)
	}
	if foundEmail == nil {
		s.log.Errorw(ctx, "email not found", logger.Err(err))
		return nil, errors.NewNotFoundError("email not found", nil)
	}

	return foundEmail, nil
}

func (s svc) FindByID(ctx context.Context, emailId string) (*Entity, error) {
	foundEmail, err := s.repo.FindByID(ctx, emailId)
	if err != nil {
		s.log.Errorw(ctx, "failed to find email by id", logger.Err(err))
		return nil, errors.NewBadRequestError("failed to find email by id", nil)
	}
	if foundEmail == nil {
		s.log.Errorw(ctx, "email not found", logger.Err(err))
		return nil, errors.NewNotFoundError("email not found", nil)
	}

	return foundEmail, nil
}

func (s svc) Create(ctx context.Context, dto CreateParams) (*Entity, error) {
	email, err := NewEmail(dto.UserID, dto.Email)
	if err != nil {
		s.log.Errorw(ctx, "failed to create email", logger.Err(err))
		return nil, errors.NewBadRequestError("failed to create email", nil)
	}
	toStore := email.Store()

	if err := s.repo.Insert(ctx, toStore); err != nil {
		s.log.Errorw(ctx, "failed to create email", logger.Err(err))
		return nil, errors.NewBadRequestError("failed to create email", nil)
	}

	return &toStore, nil
}

func (s svc) Update(ctx context.Context, entity Entity) (*Entity, error) {
	if err := s.repo.Update(ctx, entity); err != nil {
		s.log.Errorw(ctx, "failed to update email", logger.Err(err))
		return nil, errors.NewBadRequestError("failed to update email", nil)
	}

	return &entity, nil
}
