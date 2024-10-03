package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/google/go-github/v65/github"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/sharithg/siphon/internal/env"
	"github.com/sharithg/siphon/internal/storage"
	"github.com/sharithg/siphon/internal/storage/minio"
)

type Config struct {
	Github GithubConfig
	Db     DbConfig
	Addr   string
	Env    string
}

type HTTPConfig struct {
	Address string
	Port    int
}

type GithubConfig struct {
	AppConfig           githubapp.Config
	InstallationId      int64
	PullRequestPreamble string
}

type GithubOAuth struct {
	ClientID     string `yaml:"client_id" json:"clientId"`
	ClientSecret string `yaml:"client_secret" json:"clientSecret"`
}

type GhApplicationConfig struct {
	PullRequestPreamble string
}

type DbConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

type Application struct {
	Config       Config
	Store        *storage.Storage
	MinioStorage *minio.Storage
	GithubClient *github.Client
}

func (app *Application) Mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{env.GetString("CORS_ALLOWED_ORIGIN", false, "http://localhost:5173")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Route("/nodes", func(r chi.Router) {
			r.Get("/", app.getNodesHandler)
			r.Post("/", app.createNodeHandler)
		})
		r.Route("/repos", func(r chi.Router) {
			r.Get("/", app.getReposHandler)
		})
	})

	return r
}

func (app *Application) Run(mux http.Handler) error {

	srv := &http.Server{
		Addr:         app.Config.Addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		slog.Info("signal caught", "signal", s.String())

		shutdown <- srv.Shutdown(ctx)
	}()

	slog.Info("server has started", "addr", app.Config.Addr, "env", app.Config.Env)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}

	slog.Info("server has stopped", "addr", app.Config.Addr, "env", app.Config.Env)

	return nil
}
