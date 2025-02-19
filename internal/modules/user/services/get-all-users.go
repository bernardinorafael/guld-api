package usersvc

import (
	"context"
	"fmt"
	"strings"

	. "github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/internal/modules/user"
	"github.com/bernardinorafael/pkg/logger"
	"github.com/bernardinorafael/pkg/pagination"
)

func (s svc) GetAll(ctx context.Context, dto user.UserSearchParams) (*pagination.Paginated[user.Entity], error) {
	safeSort := map[string]bool{
		"full_name": true,
		"username":  true,
		"created":   true,
	}
	// Ignoring `-` preffix on verify sort opts
	sort := strings.TrimPrefix(dto.Sort, "-")
	if !safeSort[sort] {
		msg := "invalid sort params"
		s.log.Errorw(ctx, msg, logger.Err(fmt.Errorf("invalid sort params: %s", dto.Sort)))
		return nil, NewValidationFieldError(
			msg,
			fmt.Errorf("invalid sort params: %s", dto.Sort),
			[]Field{
				{Field: "sort", Msg: "invalid sort params"},
			},
		)
	}

	users, totalItems, err := s.userRepo.GetAll(ctx, dto)
	if err != nil {
		msg := "failed to retrieve users"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return nil, NewBadRequestError(msg, err)
	}

	paginated := pagination.New(users, totalItems, dto.Page, dto.Limit)
	return &paginated, nil
}
