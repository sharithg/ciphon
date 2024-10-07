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
	"github.com/go-redis/redis/v8"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/sharithg/siphon/internal/env"
	"github.com/sharithg/siphon/internal/repo"
	"github.com/sharithg/siphon/internal/storage"
	"github.com/sharithg/siphon/internal/storage/minio"
)

type Config struct {
	Github repo.GithubConfig
	Db     DbConfig
	Cache  CacheConfig
	Addr   string
	Env    string
}

type HTTPConfig struct {
	Address string
	Port    int
}

type DbConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

type CacheConfig struct {
	Addr string
}

type Application struct {
	Config       Config
	Store        *storage.Storage
	MinioStorage *minio.Storage
	Github       *repo.Github
	Wh           *GhWebhookHandler
	Cache        *redis.Client
}

func (app *Application) Mount() http.Handler {
	r := chi.NewRouter()

	handler := NewGhWebhookHandler(app.Github.ClientCreator, app.Config.Github.PullRequestPreamble, app)

	webhookHandler := githubapp.NewDefaultEventDispatcher(app.Config.Github.AppConfig, handler)

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ignoredPaths := map[string]struct{}{
				"/v1/nodes": {},
			}

			if _, ok := ignoredPaths[r.URL.Path]; ok {
				next.ServeHTTP(w, r)
				return
			}

			middleware.Logger(next).ServeHTTP(w, r)
		})
	})

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

	r.Handle(githubapp.DefaultWebhookRoute, webhookHandler)

	r.Route("/v1", func(r chi.Router) {
		r.Route("/nodes", func(r chi.Router) {
			r.Get("/", app.getNodesHandler)
			r.Post("/", app.createNodeHandler)
			r.Post("/{nodeId}", app.installToolsForNode)
		})
		r.Route("/repos", func(r chi.Router) {
			r.Get("/", app.getReposHandler)
			r.Post("/connect", app.connectRepoHandler)
			r.Get("/new", app.getNewReposHandler)
		})

		r.HandleFunc("/sse/steps/run-events/{stepId}", app.stepEventsHandler)
		r.HandleFunc("/sse/workflows/run-events", app.eventsHandler)

		r.Route("/workflows", func(r chi.Router) {
			r.Get("/", app.getWorkflows)
			r.Post("/trigger/{workflowId}", app.triggerWorkflow)
			r.Route("/{workflowId}", func(r chi.Router) {
				r.Get("/jobs", app.getJobs)
				r.Route("/jobs/{jobId}", func(r chi.Router) {
					r.Get("/steps", app.getSteps)
					r.Get("/steps/{stepId}/output", app.getStepOutput)
				})
			})
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
