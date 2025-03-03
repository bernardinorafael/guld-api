package usersvc

import (
	"github.com/bernardinorafael/internal/mailer"
	"github.com/bernardinorafael/internal/modules/email"
	"github.com/bernardinorafael/internal/modules/user"
	"github.com/bernardinorafael/internal/uploader"
	"github.com/bernardinorafael/pkg/logger"
)

type svc struct {
	log          logger.Logger
	userRepo     user.RepositoryInterface
	emailService email.ServiceInterface
	mailer       mailer.Mailer
	uploader     uploader.UploaderInterface
}

func New(
	log logger.Logger,
	userRepo user.RepositoryInterface,
	emailService email.ServiceInterface,
	mailer mailer.Mailer,
	uploader uploader.UploaderInterface,
) user.ServiceInterface {
	return &svc{
		log:          log,
		userRepo:     userRepo,
		emailService: emailService,
		mailer:       mailer,
		uploader:     uploader,
	}
}
