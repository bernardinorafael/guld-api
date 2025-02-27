package usersvc

import (
	"context"
	"fmt"
	"strings"

	"github.com/bernardinorafael/internal/_shared/dto"
	. "github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/internal/infra/http/middleware"
	"github.com/bernardinorafael/internal/modules/user"
	"github.com/bernardinorafael/pkg/logger"
	"github.com/bernardinorafael/pkg/pagination"
)

func (s svc) GetAll(
	ctx context.Context,
	dto dto.SearchParams,
) (*pagination.Paginated[user.EntityWithTeam], error) {
	safeSort := map[string]bool{
		"full_name": true,
		"username":  true,
		"created":   true,
	}
	// Ignoring `-` preffix on verify sort opts
	sort := strings.TrimPrefix(dto.Sort, "-")
	if !safeSort[sort] {
		return nil, NewValidationFieldError(
			"invalid sort params",
			fmt.Errorf("invalid sort params: %s", dto.Sort),
			[]Field{
				{Field: "sort", Msg: "invalid sort params"},
			},
		)
	}

	rawUsers, totalItems, err := s.userRepo.GetAll(ctx, dto)
	if err != nil {
		s.log.Errorw(ctx, "failed to retrieve users", logger.Err(err))
		return nil, NewBadRequestError("failed to retrieve users", err)
	}

	userInCtx, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return nil, NewBadRequestError("failed to retrieve user id", err)
	}

	// Filter out the user in context
	// TODO: Move this to repository
	users := make([]user.EntityWithTeam, 0)
	for _, u := range rawUsers {
		if u.ID != userInCtx {
			users = append(users, u)
		}
	}
	totalItems--

	// This is a workaround to avoid null json on team field
	for i := range users {
		if users[i].Team.ID == nil {
			users[i].Team = nil
		}
	}

	paginated := pagination.New(users, totalItems, dto.Page, dto.Limit)
	return &paginated, nil
}
