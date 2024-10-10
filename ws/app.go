package ws

import (
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sharithg/siphon/internal/config"
	"github.com/sharithg/siphon/internal/docker"
	storage "github.com/sharithg/siphon/internal/storage/kv"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Optionally handle origin checks here if needed
		return true
	},
}

type Application struct {
	Config      Config
	Store       *storage.KvStorage
	Docker      *docker.Docker
	AgentConfig *config.AgentConfig
}

type Config struct {
	Addr string
	Env  string
}

func ipWhitelisted(r *http.Request, whitelist []string) bool {
	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil {
		return false
	}
	for _, allowedIP := range whitelist {
		if allowedIP == remoteIP {
			return true
		}
	}
	return false
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
