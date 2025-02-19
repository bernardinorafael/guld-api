package usersvc

import (
	"github.com/bernardinorafael/internal/mailer"
	"github.com/bernardinorafael/internal/modules/email"
	"github.com/bernardinorafael/internal/modules/phone"
	"github.com/bernardinorafael/internal/modules/user"
	"github.com/bernardinorafael/pkg/logger"
)

const (
	maxEmailAndPhoneByUser = 3
)

type svc struct {
	log          logger.Logger
	userRepo     user.RepositoryInterface
	emailService email.ServiceInterface
	phoneRepo    phone.RepositoryInterface
	mailer       mailer.Mailer
}

func New(
	log logger.Logger,
	userRepo user.RepositoryInterface,
	emailService email.ServiceInterface,
	phoneRepo phone.RepositoryInterface,
	mailer mailer.Mailer,
) user.ServiceInterface {
	return &svc{
		log:          log,
		userRepo:     userRepo,
		emailService: emailService,
		phoneRepo:    phoneRepo,
		mailer:       mailer,
	}
}
