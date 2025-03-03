package usersvc

import (
	"context"
	"fmt"
	"io"

	. "github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/internal/_shared/util"
	"github.com/bernardinorafael/internal/modules/user"
	"github.com/bernardinorafael/pkg/logger"
)

// TODO: Salva no banco apenas o ID do avatar, não o URL. Posteriormente construir a URL do avatar
// Com base nas informações do bucket e do ID do avatar

func (s svc) UpdateAvatar(ctx context.Context, userId string, file io.Reader, filename string) error {
	u, err := s.userRepo.FindByID(ctx, userId)
	if err != nil {
		return NewBadRequestError("failed to retrieve user", err)
	}
	if u == nil {
		return NewNotFoundError("user not found", nil)
	}

	key := fmt.Sprintf("%s/%s", u.ID, util.GenID("avatar"))
	out, err := s.uploader.UploadFile(ctx, file, key)
	if err != nil {
		s.log.Errorw(ctx, "error on upload avatar", logger.Err(err))
	}

	userEntity, err := user.NewFromEntity(*u)
	if err != nil {
		return NewValidationFieldError("error on init user entity", err, nil)
	}
	userEntity.ChangeProfilePicture(out.Location)
	toStore := userEntity.Store()

	err = s.userRepo.Update(ctx, toStore)
	if err != nil {
		return NewBadRequestError("error on update profile", err)
	}

	return nil
}
