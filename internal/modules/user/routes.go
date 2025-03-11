package user

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bernardinorafael/internal/_shared/dto"
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
	m := middleware.NewWithAuth(c.log, c.secretKey)

	r.Route("/api/v1/users", func(r chi.Router) {
		r.Use(m.WithAuth)

		// Users
		r.Get("/", c.getAllUsers)
		r.Post("/", c.create)
		r.Get("/{userId}", c.getUser)
		r.Delete("/{userId}", c.delete)
		r.Patch("/{userId}/toggle-lock", c.toggleLock)
		r.Patch("/{userId}/profile", c.updateProfile)
		r.Patch("/{userId}/profile/avatar", c.uploadAvatar)

		// Emails
		// TODO: Move this to a dedicated emails router
		r.Get("/{userId}/emails", c.getEmails)
		r.Patch("/{userId}/emails/{emailId}/set-primary", c.setPrimaryEmail)
		r.Get("/emails/{emailId}", c.getEmail)
	})
}

func (c controller) uploadAvatar(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")

	if err := r.ParseMultipartForm(3 << 20); err != nil {
		NewHttpError(w, fmt.Errorf("maximum file size is 3MB"))
		return
	}

	f, fh, err := r.FormFile("avatar")
	if err != nil {
		NewHttpError(w, fmt.Errorf("error processing file: %w", err))
		return
	}
	defer f.Close()

	err = c.svc.UpdateAvatar(c.ctx, userId, f, fh.Filename)
	if err != nil {
		NewHttpError(w, err)
		return
	}

	util.WriteSuccessResponse(w, http.StatusOK)
}

func (c controller) updateProfile(w http.ResponseWriter, r *http.Request) {
	var body UpdateProfileDTO
	var userId = chi.URLParam(r, "userId")

	err := util.ReadRequestBody(w, r, &body)
	if err != nil {
		NewHttpError(w, err)
		return
	}

	err = c.svc.UpdateProfile(c.ctx, userId, body)
	if err != nil {
		NewHttpError(w, err)
		return
	}

	util.WriteSuccessResponse(w, http.StatusOK)
}

func (c controller) getEmail(w http.ResponseWriter, r *http.Request) {
	email, err := c.svc.FindEmail(r.Context(), chi.URLParam(r, "emailId"))
	if err != nil {
		NewHttpError(w, err)
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, map[string]any{"email": email})
}

func (c controller) setPrimaryEmail(w http.ResponseWriter, r *http.Request) {
	if err := c.svc.SetPrimaryEmail(
		c.ctx,
		chi.URLParam(r, "userId"),
		chi.URLParam(r, "emailId"),
	); err != nil {
		NewHttpError(w, err)
		return
	}

	util.WriteSuccessResponse(w, http.StatusOK)
}

func (c controller) getEmails(w http.ResponseWriter, r *http.Request) {
	emails, err := c.svc.FindAllEmails(c.ctx, chi.URLParam(r, "userId"))
	if err != nil {
		NewHttpError(w, err)
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, map[string]any{"emails": emails})
}

func (c controller) toggleLock(w http.ResponseWriter, r *http.Request) {
	if err := c.svc.ToggleLock(c.ctx, chi.URLParam(r, "userId")); err != nil {
		NewHttpError(w, err)
		return
	}

	util.WriteSuccessResponse(w, http.StatusOK)
}

func (c controller) delete(w http.ResponseWriter, r *http.Request) {
	if err := c.svc.Delete(c.ctx, chi.URLParam(r, "userId")); err != nil {
		NewHttpError(w, err)
		return
	}

	util.WriteSuccessResponse(w, http.StatusOK)
}

func (c controller) create(w http.ResponseWriter, r *http.Request) {
	var body UserRegisterDTO

	if err := util.ReadRequestBody(w, r, &body); err != nil {
		NewHttpError(w, err)
		return
	}

	userId, err := c.svc.Create(c.ctx, body)
	if err != nil {
		NewHttpError(w, err)
		return
	}

	util.WriteJSONResponse(w, http.StatusCreated, map[string]any{
		"user_id": userId,
	})
}

func (c controller) getUser(w http.ResponseWriter, r *http.Request) {
	res, err := c.svc.FindByID(c.ctx, chi.URLParam(r, "userId"))
	if err != nil {
		NewHttpError(w, err)
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"user":   res.User,
		"emails": res.Emails,
		"meta":   res.Meta,
	})
}

func (c controller) getAllUsers(w http.ResponseWriter, r *http.Request) {
	var p dto.SearchParams

	p.Query = util.ReadQueryString(r.URL.Query(), "q", "")
	p.Limit = util.ReadQueryInt(r.URL.Query(), "limit", 25)
	p.Page = util.ReadQueryInt(r.URL.Query(), "page", 1)
	p.Sort = util.ReadQueryString(r.URL.Query(), "sort", "created")

	res, err := c.svc.GetAll(r.Context(), p)
	if err != nil {
		NewHttpError(w, err)
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"data": res.Data,
		"meta": res.Meta,
	})
}
