package agent

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sharithg/siphon/internal/config"
	"github.com/sharithg/siphon/internal/docker"
	"github.com/sharithg/siphon/internal/runner"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Optionally handle origin checks here if needed
		return true
	},
}

type Application struct {
	Config      Config
	Docker      *docker.Docker
	AgentConfig *config.AgentConfig
	Runner      *runner.Runner
}

type Config struct {
	Addr string
	Env  string
}

func authMiddleware(token string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("X-Ciphon-Auth")
		if authToken != token {
			http.Error(w, "Forbidden: Invalid auth", http.StatusForbidden)
			return
		}
		next(w, r)
	}
}

func (app *Application) Run() error {
	log.Printf("Server has started on %s, env: %s", app.Config.Addr, app.Config.Env)

	http.HandleFunc("/ws", authMiddleware(app.AgentConfig.Token, app.serveWs))
	log.Fatal(http.ListenAndServe(app.Config.Addr, nil))
	return nil
}
