package usersvc

import (
	"github.com/bernardinorafael/internal/modules/email"
	"github.com/bernardinorafael/internal/modules/phone"
	"github.com/bernardinorafael/internal/modules/user"
	"github.com/bernardinorafael/pkg/logger"
)

type svc struct {
	log       logger.Logger
	userRepo  user.RepositoryInterface
	emailRepo email.RepositoryInterface
	phoneRepo phone.RepositoryInterface
}

func New(
	log logger.Logger,
	userRepo user.RepositoryInterface,
	emailRepo email.RepositoryInterface,
	phoneRepo phone.RepositoryInterface,
) user.ServiceInterface {
	return &svc{
		log:       log,
		userRepo:  userRepo,
		emailRepo: emailRepo,
		phoneRepo: phoneRepo,
	}
}
