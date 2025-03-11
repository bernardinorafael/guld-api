package email

import (
	"context"
	"net/http"

	"github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/internal/_shared/util"
	"github.com/bernardinorafael/internal/infra/http/middleware"
	"github.com/bernardinorafael/pkg/logger"
	"github.com/go-chi/chi"
)

type controller struct {
	ctx       context.Context
	log       logger.Logger
	svc       ServiceInterface
	secretKey string
}

func NewController(
	ctx context.Context,
	log logger.Logger,
	svc ServiceInterface,
	secretKey string,
) *controller {
	return &controller{
		ctx:       ctx,
		log:       log,
		svc:       svc,
		secretKey: secretKey,
	}
}

func (c controller) RegisterRoute(r *chi.Mux) {
	m := middleware.NewWithAuth(c.log, c.secretKey)

	r.Route("/api/v1/emails", func(r chi.Router) {
		r.Use(m.WithAuth)

		r.Post("/{userId}", c.create)
		r.Delete("/{userId}/{emailId}", c.delete)
	})

	r.Route("/api/v1/emails/validations", func(r chi.Router) {
		r.Use(m.WithAuth)
		r.Post("/", c.requestValidation)
		r.Post("/{emailId}", c.validateEmail)
	})
}

func (c controller) delete(w http.ResponseWriter, r *http.Request) {
	err := c.svc.DeleteEmail(
		c.ctx,
		chi.URLParam(r, "userId"),
		chi.URLParam(r, "emailId"),
	)
	if err != nil {
		errors.NewHttpError(w, err)
		return
	}

	util.WriteSuccessResponse(w, http.StatusOK)
}

func (c controller) create(w http.ResponseWriter, r *http.Request) {
	var body CreateEmailDTO
	body.UserID = chi.URLParam(r, "userId")

	err := util.ReadRequestBody(w, r, &body)
	if err != nil {
		errors.NewHttpError(w, err)
		return
	}

	err = c.svc.AddEmail(c.ctx, body)
	if err != nil {
		errors.NewHttpError(w, err)
		return
	}

	util.WriteSuccessResponse(w, http.StatusOK)
}

func (c controller) validateEmail(w http.ResponseWriter, r *http.Request) {
	var body ValidateEmailDTO
	body.EmailID = chi.URLParam(r, "emailId")

	err := util.ReadRequestBody(w, r, &body)
	if err != nil {
		errors.NewHttpError(w, err)
		return
	}

	err = c.svc.ValidateEmail(c.ctx, body)
	if err != nil {
		errors.NewHttpError(w, err)
		return
	}

	util.WriteSuccessResponse(w, http.StatusOK)
}

func (c controller) requestValidation(w http.ResponseWriter, r *http.Request) {
	var body GenerateEmailValidationDTO

	err := util.ReadRequestBody(w, r, &body)
	if err != nil {
		errors.NewHttpError(w, err)
		return
	}

	err = c.svc.GenerateValidationCode(c.ctx, body)
	if err != nil {
		errors.NewHttpError(w, err)
		return
	}

	util.WriteSuccessResponse(w, http.StatusOK)
}
