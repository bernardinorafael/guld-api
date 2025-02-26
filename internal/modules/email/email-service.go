package email

import (
	"context"
	"errors"

	. "github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/internal/_shared/util"
	"github.com/bernardinorafael/internal/mailer"
	"github.com/bernardinorafael/pkg/logger"
	"github.com/lib/pq"
)

type svc struct {
	log    logger.Logger
	repo   RepositoryInterface
	mailer mailer.Mailer
}

func NewService(log logger.Logger, repo RepositoryInterface, mailer mailer.Mailer) ServiceInterface {
	return &svc{log, repo, mailer}
}

func (s svc) AddEmail(ctx context.Context, dto CreateEmailDTO) error {
	emails, err := s.FindAllByUser(ctx, dto.UserID)
	if err != nil {
		return err
	}

	if len(emails) >= maxEmailAndPhoneByUser {
		s.log.Info(ctx, "max emails reached")
		return NewForbiddenError("max emails reached", MaxLimitResourceReached, nil)
	}

	for _, v := range emails {
		if v.Email == dto.Email {
			return NewConflictError("email already taken", ResourceAlreadyTaken, nil,
				[]Field{{Field: "email", Msg: "email already taken"}},
			)
		}
	}

	email, err := NewEmail(dto.UserID, dto.Email)
	if err != nil {
		return NewBadRequestError("error on create email", err)
	}
	toStore := email.Store()

	if err := s.repo.Insert(ctx, toStore); err != nil {
		s.log.Errorw(ctx, "failed to insert email", logger.Err(err))
		return NewBadRequestError("failed to insert email", err)
	}

	if dto.SendCode {
		if err := s.GenerateValidationCode(ctx, GenerateEmailValidationDTO{
			EmailID: email.ID(),
			UserID:  dto.UserID,
		}); err != nil {
			s.log.Errorw(ctx, "failed to generate validation code", logger.Err(err))
		}
	}

	return nil
}

func (s svc) Delete(ctx context.Context, userId, emailId string) error {
	if err := s.repo.Delete(ctx, userId, emailId); err != nil {
		s.log.Errorw(ctx, "failed to delete email", logger.Err(err))
		return NewBadRequestError("failed to delete email", nil)
	}

	return nil
}

func (s svc) FindAllByUser(ctx context.Context, userId string) ([]Entity, error) {
	emails, err := s.repo.FindAllByUser(ctx, userId)
	if err != nil {
		s.log.Errorw(ctx, "failed to find all emails by user", logger.Err(err))
		return nil, NewBadRequestError("failed to find all emails by user", err)
	}

	return emails, nil
}

func (s svc) FindByEmail(ctx context.Context, email string) (*Entity, error) {
	foundEmail, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		s.log.Errorw(ctx, "failed to find email by email", logger.Err(err))
		return nil, NewBadRequestError("failed to find email by email", err)
	}
	if foundEmail == nil {
		s.log.Errorw(ctx, "email not found", logger.Err(err))
		return nil, NewNotFoundError("email not found", err)
	}

	return foundEmail, nil
}

func (s svc) FindByID(ctx context.Context, emailId string) (*Entity, error) {
	foundEmail, err := s.repo.FindByID(ctx, emailId)
	if err != nil {
		s.log.Errorw(ctx, "failed to find email by id", logger.Err(err))
		return nil, NewBadRequestError("failed to find email by id", err)
	}
	if foundEmail == nil {
		s.log.Errorw(ctx, "email not found", logger.Err(err))
		return nil, NewNotFoundError("email not found", err)
	}

	return foundEmail, nil
}

func (s svc) Create(ctx context.Context, dto CreateParams) (*Entity, error) {
	email, err := NewEmail(dto.UserID, dto.Email)
	if err != nil {
		s.log.Errorw(ctx, "failed to create email", logger.Err(err))
		return nil, NewBadRequestError("failed to create email", err)
	}
	toStore := email.Store()

	if err := s.repo.Insert(ctx, toStore); err != nil {
		var pqErr *pq.Error
		// 23505 is the code for unique constraint violation
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			msg := "email already exists"
			var appErr = NewConflictError(msg, ResourceAlreadyTaken, err, nil)
			field := util.ExtractFieldFromDetail(pqErr.Detail)
			s.log.Errorw(ctx, msg, logger.Err(err))
			appErr.AddField(field, field+" already exists")
			return nil, appErr
		}
		s.log.Errorw(ctx, "failed to create email", logger.Err(err))
		return nil, NewBadRequestError("failed to create email", err)
	}

	return &toStore, nil
}

func (s svc) Update(ctx context.Context, entity Entity) (*Entity, error) {
	if err := s.repo.Update(ctx, entity); err != nil {
		s.log.Errorw(ctx, "failed to update email", logger.Err(err))
		return nil, NewBadRequestError("failed to update email", err)
	}

	return &entity, nil
}
