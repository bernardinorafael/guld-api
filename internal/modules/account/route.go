package account

import (
	"context"
	"net/http"

	. "github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/internal/_shared/util"
	"github.com/bernardinorafael/internal/infra/http/middleware"
	"github.com/bernardinorafael/internal/infra/token"
	"github.com/bernardinorafael/pkg/logger"
	"github.com/go-chi/chi"
)

type controller struct {
	ctx       context.Context
	log       logger.Logger
	svc       ServiceInterface
	secretKey string
}

const basePath = "/api/v1"

func NewController(
	ctx context.Context,
	log logger.Logger,
	svc ServiceInterface,
	secretKey string,
) *controller {
	return &controller{ctx, log, svc, secretKey}
}

func (c controller) RegisterRoute(r *chi.Mux) {
	m := middleware.NewWithAuth(c.log, c.secretKey)

	r.Route(basePath+"/auth", func(r chi.Router) {
		// Public
		r.Post("/login", c.login)
		r.Post("/refresh", c.renewRefreshToken)

		// Private
		r.With(m.WithAuth).Delete("/logout", c.logOut)

		// register
		// activate account
		// forgot password
		// reset password
	})

	r.Route(basePath+"/accounts", func(r chi.Router) {
		r.Use(m.WithAuth)
		r.Get("/me", c.getSigned)
		r.Post("/{userId}/change-password", c.changePassword)
	})

	r.Route(basePath+"/sessions", func(r chi.Router) {
		r.Use(m.WithAuth)

		r.Get("/", c.getAllSessions)
		r.Get("/active", c.getSession)
		r.Patch("/{sessionId}/revoke", c.revokeSession)
	})
}

func (c controller) revokeSession(w http.ResponseWriter, r *http.Request) {
	var sessionId = chi.URLParam(r, "sessionId")
	claims := r.Context().Value(middleware.AuthKey{}).(*token.AccountClaims)

	err := c.svc.RevokeSession(c.ctx, claims.Username, sessionId)
	if err != nil {
		NewHttpError(w, err)
		return
	}

	util.WriteSuccessResponse(w, http.StatusOK)
}

func (c controller) getAllSessions(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.AuthKey{}).(*token.AccountClaims)

	sessions, err := c.svc.GetAllSessions(c.ctx, claims.Username)
	if err != nil {
		NewHttpError(w, err)
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, sessions)
}

func (c controller) renewRefreshToken(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}

	err := util.ReadRequestBody(w, r, &body)
	if err != nil {
		NewHttpError(w, err)
		return
	}

	payload, err := c.svc.RenewAccessToken(c.ctx, body.RefreshToken)
	if err != nil {
		NewHttpError(w, err)
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, payload)
}

func (c controller) logOut(w http.ResponseWriter, r *http.Request) {
	if err := c.svc.Logout(r.Context()); err != nil {
		NewHttpError(w, err)
		return
	}

	util.WriteSuccessResponse(w, http.StatusOK)
}

func (c controller) changePassword(w http.ResponseWriter, r *http.Request) {
	var userId = chi.URLParam(r, "userId")
	var body struct {
		Password    string `json:"password"`
		NewPassword string `json:"new_password"`
	}

	err := util.ReadRequestBody(w, r, &body)
	if err != nil {
		NewHttpError(w, err)
		return
	}

	err = c.svc.ChangePassword(c.ctx, userId, body.Password, body.NewPassword)
	if err != nil {
		NewHttpError(w, err)
		return
	}

	util.WriteSuccessResponse(w, http.StatusOK)
}

func (c controller) getSigned(w http.ResponseWriter, r *http.Request) {
	acc, err := c.svc.GetSignedInAccount(r.Context())
	if err != nil {
		NewHttpError(w, err)
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"account": acc,
	})
}

func (c controller) login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := util.ReadRequestBody(w, r, &body); err != nil {
		NewHttpError(w, err)
		return
	}

	payload, err := c.svc.Login(
		r.Context(),
		body.Username,
		body.Password,
		r.UserAgent(),
		r.RemoteAddr,
	)
	if err != nil {
		NewHttpError(w, err)
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, payload)
}

func (c controller) getSession(w http.ResponseWriter, r *http.Request) {
	session, err := c.svc.GetSession(r.Context())
	if err != nil {
		NewHttpError(w, err)
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, session)
}
