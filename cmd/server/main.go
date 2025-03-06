package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/bernardinorafael/internal/_shared/envconf"
	"github.com/bernardinorafael/internal/_shared/loggerconf"
	"github.com/bernardinorafael/internal/infra/database/pg"
	"github.com/bernardinorafael/internal/infra/http/middleware"
	"github.com/bernardinorafael/internal/mailer"
	"github.com/bernardinorafael/internal/modules/account"
	"github.com/bernardinorafael/internal/modules/email"
	"github.com/bernardinorafael/internal/modules/org"
	"github.com/bernardinorafael/internal/modules/permission"
	"github.com/bernardinorafael/internal/modules/role"
	"github.com/bernardinorafael/internal/modules/team"
	"github.com/bernardinorafael/internal/modules/user"
	userrepo "github.com/bernardinorafael/internal/modules/user/repository"
	usersvc "github.com/bernardinorafael/internal/modules/user/services"
	"github.com/bernardinorafael/internal/uploader"
	"github.com/bernardinorafael/pkg/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func main() {
	ctx := context.Background()

	r := chi.NewRouter()
	r.Use(middleware.WithRecoverPanic)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	env, err := envconf.New()
	if err != nil {
		slog.Error("failed load env", "error", err)
		panic(err)
	}
	// NOTE: To turn off logging, set env.Debug to false
	// NOTE: Passing empty string as name to avoid logging the name of the service
	log := loggerconf.New("", env.Debug)

	db, err := pg.NewDatabase(log, env.DSN)
	if err != nil {
		log.Errorw(ctx, "failed to connect to database", logger.Err(err))
		panic(err)
	}
	defer db.Close()
	log.Info(ctx, "Database connected")

	uploader := uploader.NewUploader(ctx, log)

	mailer := mailer.New(ctx, log, mailer.Config{
		APIKey:           env.ResendAPIKey,
		MaxRetries:       3,
		RetryDelay:       time.Second * 2,
		OperationTimeout: time.Second * 10,
	})

	// Repositories
	permissionRepo := permission.NewRepository(db.GetDB())
	userRepo := userrepo.New(db.GetDB())
	emailRepo := email.NewRepository(db.GetDB())
	roleRepo := role.NewRepository(db.GetDB())
	teamRepo := team.NewRepository(db.GetDB())
	accRepo := account.NewRepository(db.GetDB())
	orgRepo := org.NewRepo(db.GetDB())

	// Services
	emailService := email.NewService(log, emailRepo, mailer)

	permissionService := permission.NewService(log, permissionRepo)
	accService := account.NewService(ctx, log, accRepo, mailer, env.JWTSecret)
	roleService := role.NewService(log, roleRepo)
	teamService := team.NewService(log, teamRepo)
	userService := usersvc.New(log, userRepo, emailService, mailer, uploader)
	orgService := org.NewService(log, orgRepo)

	// Controllers
	email.NewController(ctx, log, emailService, env.JWTSecret).RegisterRoute(r)
	team.NewController(ctx, log, teamService, env.JWTSecret).RegisterRoute(r)
	user.NewController(ctx, log, userService, env.JWTSecret).RegisterRoute(r)
	role.NewController(ctx, log, roleService, env.JWTSecret).RegisterRoute(r)
	account.NewController(ctx, log, accService, env.JWTSecret).RegisterRoute(r)
	org.NewController(ctx, log, orgService, env.JWTSecret).RegisterRoute(r)
	permission.NewController(ctx, log, permissionService, env.JWTSecret).RegisterRoute(r)

	log.Info(ctx, "Server started")
	err = http.ListenAndServe(":"+env.Port, r)
	if err != nil {
		log.Errorw(ctx, "failed to start server", logger.Err(err))
		os.Exit(1)
	}
}
