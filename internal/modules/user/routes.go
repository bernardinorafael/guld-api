package user

import (
	"context"
	"net/http"

	. "github.com/bernardinorafael/internal/_shared/errors"
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
	m := middleware.NewWithAuth(c.ctx, c.log, c.secretKey)

	r.Route("/api/v1/users", func(r chi.Router) {
		r.Use(m.WithAuth)
		r.Get("/", c.getAllUsers)
		r.Get("/{id}", c.getUser)
		r.Post("/", c.create)
		r.Delete("/{id}", c.delete)
		r.Patch("/{userId}/toggle-lock", c.toggleLock)
	})
}

func (c controller) toggleLock(w http.ResponseWriter, r *http.Request) {
	if err := c.svc.ToggleLock(c.ctx, chi.URLParam(r, "userId")); err != nil {
		if err, ok := err.(ApplicationError); ok {
			NewHttpError(w, err)
			return
		}
		NewHttpError(w, NewInternalServerError(err))
		return
	}

	util.WriteSuccessResponse(w, http.StatusOK)
}

func (c controller) delete(w http.ResponseWriter, r *http.Request) {
	if err := c.svc.Delete(c.ctx, chi.URLParam(r, "id")); err != nil {
		if err, ok := err.(ApplicationError); ok {
			NewHttpError(w, err)
			return
		}
		NewHttpError(w, NewInternalServerError(err))
		return
	}

	util.WriteSuccessResponse(w, http.StatusOK)
}

func (c controller) create(w http.ResponseWriter, r *http.Request) {
	var body UserRegisterParams

	if err := util.ReadRequestBody(w, r, &body); err != nil {
		if err, ok := err.(ApplicationError); ok {
			NewHttpError(w, err)
			return
		}
		NewHttpError(w, NewInternalServerError(err))
		return
	}

	if err := c.svc.Create(c.ctx, body); err != nil {
		if err, ok := err.(ApplicationError); ok {
			NewHttpError(w, err)
			return
		}
		NewHttpError(w, NewInternalServerError(err))
		return
	}

	util.WriteSuccessResponse(w, http.StatusCreated)
}

func (c controller) getUser(w http.ResponseWriter, r *http.Request) {
	res, err := c.svc.FindByID(c.ctx, chi.URLParam(r, "id"))
	if err != nil {
		if err, ok := err.(ApplicationError); ok {
			NewHttpError(w, err)
			return
		}
		NewHttpError(w, NewInternalServerError(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"user":   res.User,
		"emails": res.Emails,
		"phones": res.Phones,
		"meta":   res.Meta,
	})
}

func (c controller) getAllUsers(w http.ResponseWriter, r *http.Request) {
	var p UserSearchParams

	p.Query = util.ReadQueryString(r.URL.Query(), "q", "")
	p.Limit = util.ReadQueryInt(r.URL.Query(), "limit", 25)
	p.Page = util.ReadQueryInt(r.URL.Query(), "page", 1)
	p.Sort = util.ReadQueryString(r.URL.Query(), "sort", "created")

	res, err := c.svc.GetAll(c.ctx, p)
	if err != nil {
		if err, ok := err.(ApplicationError); ok {
			NewHttpError(w, err)
			return
		}
		NewHttpError(w, NewInternalServerError(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"data": res.Data,
		"meta": res.Meta,
	})
}
