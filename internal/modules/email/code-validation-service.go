package email

import (
	"context"
	"fmt"

	. "github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/internal/mailer"
	"github.com/bernardinorafael/pkg/logger"
)

func (s svc) FindActiveCodesByUser(ctx context.Context, userId string) ([]ValidationEntity, error) {
	foundCodes, err := s.repo.FindAllCodesByUser(ctx, userId)
	if err != nil {
		return nil, NewBadRequestError("failed to find active codes by user", nil)
	}

	codes := make([]ValidationEntity, 0)

	for _, code := range foundCodes {
		if code.IsValid {
			codes = append(codes, code)
		}
	}

	return codes, nil
}

func (s svc) ValidateEmail(ctx context.Context, dto ValidateEmailDTO) error {
	email, err := s.repo.FindByID(ctx, dto.EmailID)
	if err != nil {
		return NewBadRequestError("failed to find email", nil)
	}

	existsCode, err := s.repo.FindCodeValidationByEmailId(ctx, email.ID)
	if err != nil {
		return NewBadRequestError("failed to find email validation", nil)
	}
	if existsCode == nil {
		return NewConflictError("email validation not found", ResourceNotFound, nil, nil)
	}

	code := NewCodeValidationFromEntity(*existsCode)

	if !code.IsValid() || code.IsConsumed() {
		return NewForbiddenError("code is invalid and/or already taken", Expired, nil)
	}

	if code.IsMaxAttempts() {
		// If the code has reached the max attempts:
		// - Invalidate the code verification
		// - Update the code verification in the database
		// - Return a forbidden error
		code.Invalidate()
		if err := s.repo.UpdateCodeValidation(ctx, code.Store()); err != nil {
			return NewBadRequestError("failed to update email validation", err)
		}
		return NewForbiddenError("max attempts reached", MaxLimitResourceReached, nil)
	}

	if code.IsExpired() {
		// If the code has expired:
		// - Invalidate the code verification
		// - Update the code verification in the database
		// - Return a forbidden error
		code.Invalidate()
		if err := s.repo.UpdateCodeValidation(ctx, code.Store()); err != nil {
			return NewBadRequestError("failed to update email validation", err)
		}
		return NewForbiddenError("code expired", Expired, nil)
	}

	code.IncrementAttempts()

	if !code.ValidateCode(dto.Code) {
		if err := s.repo.UpdateCodeValidation(ctx, code.Store()); err != nil {
			return NewBadRequestError("failed to update email validation", err)
		}
		return NewValidationFieldError("invalid code", nil, []Field{
			{Field: "code", Msg: "invalid code"},
		})
	}

	emailEntity, err := NewFromEntity(email)
	if err != nil {
		return NewBadRequestError("failed to create email entity", err)
	}
	emailEntity.Verify()

	if err := s.repo.Update(ctx, emailEntity.Store()); err != nil {
		return NewBadRequestError("failed to update email", err)
	}

	code.Consume()
	if err := s.repo.UpdateCodeValidation(ctx, code.Store()); err != nil {
		return NewBadRequestError("failed to update email validation", err)
	}

	s.log.Info(ctx, "code validated successfully")
	return nil
}

func (s svc) GenerateValidationCode(ctx context.Context, dto GenerateEmailValidationDTO) error {
	email, err := s.repo.FindByID(ctx, dto.EmailID)
	if err != nil {
		return NewBadRequestError("failed to find email", nil)
	}

	// Check if the email validation already exists
	exists, err := s.repo.FindCodeValidationByEmailId(ctx, email.ID)
	if err != nil {
		return NewBadRequestError("failed to find email validation", nil)
	}
	// If the email validation already exists, check if it is expired or valid
	if exists != nil {
		exists := NewCodeValidationFromEntity(*exists)

		if !exists.IsExpired() && exists.IsValid() {
			return NewForbiddenError("email validation already exists", ResourceAlreadyTaken, nil)
		}
	}

	newCode := NewCodeValidation(email.ID, dto.UserID)

	go func() {
		params := mailer.SendParams{
			From:    mailer.NotificationSender,
			To:      "rafaelferreirab2@gmail.com",
			Subject: fmt.Sprintf("%s - Este é seu código de validação", newCode.Code()),
			File:    "email_validation.html",
			Data: map[string]any{
				"Email": email.Email,
				"Code":  newCode.Code(),
			},
		}

		err := s.mailer.Send(params)
		if err != nil {
			s.log.Errorw(ctx, "failed to send email validation", logger.Err(err))
		}

		s.log.Infow(ctx, "email validation sent", logger.String("email", email.Email))
	}()

	err = s.repo.InsertCodeValidation(ctx, newCode.Store())
	if err != nil {
		s.log.Errorw(ctx, "failed to insert email validation", logger.Err(err))
		return NewBadRequestError("failed to insert email validation", nil)
	}

	return nil
}
