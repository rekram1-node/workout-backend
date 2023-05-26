package main

import (
	"os"

	env "github.com/caarlos0/env/v6"
	"github.com/go-chi/chi"
	"github.com/rekram1-node/httptemplate"
	"github.com/rekram1-node/workout-backend/handlers"
	"github.com/rekram1-node/workout-backend/middleware"
	"github.com/rekram1-node/workout-backend/repository"
	"github.com/rs/zerolog"
)

type config struct {
	PgURI     string `env:"PG_URI,required"`
	JWTSecret string `env:"JWT_SECRET,required"`
}

func main() {
	cfg := config{}
	logger := zerolog.New(os.Stdout)

	if err := env.Parse(&cfg); err != nil {
		logger.Fatal().Err(err).Msg("failed to read configuration")
	}

	db, err := repository.New(cfg.PgURI)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to database")
	}

	app, err := httptemplate.New("workout-backend")
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create template application")
	}

	jwt := &middleware.JWTAuthentication{
		SecretKey: cfg.JWTSecret,
	}

	app.Router.Route("/client-services", func(r chi.Router) {
		r.Route("/user", func(usr chi.Router) {
			usr.Post("/signin", handlers.LoginHandler(db, cfg.JWTSecret))
			usr.Post("/", handlers.UserCreate(db, cfg.JWTSecret))
			usr.With(jwt.Authentication).Get("/", handlers.UserRead(db))
			usr.With(jwt.Authentication).Put("/", handlers.UserUpdate(db))
			usr.With(jwt.Authentication).Delete("/", handlers.UserDelete(db))
		})
		r.Route("/meso", func(meso chi.Router) {
			meso.With(jwt.Authentication).Post("/", handlers.MesoCreate(db))
			meso.With(jwt.Authentication).Get("/", handlers.MesoRead(db))
			meso.With(jwt.Authentication).Get("/top", handlers.MesosRead(db))
			meso.With(jwt.Authentication).Put("/", handlers.UpdateMeso(db))
			meso.With(jwt.Authentication).Delete("/", handlers.DeleteMeso(db))
		})
	})

	app.Start()
}
