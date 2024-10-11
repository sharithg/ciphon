package agent

import (
	"log"
	"net/http"

	"github.com/sharithg/siphon/internal/runner"
)

func (app *Application) serveWs(w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		return
	}
	defer ws.Close()

	log.Printf("Client connected from %s", ws.RemoteAddr())

	for {
		var event runner.Commands
		err := ws.ReadJSON(&event)
		if err != nil {
			log.Printf("Receive failed: %s", err.Error())
			break
		}
		switch event.Type {
		case "run_command":
			go func() {
				if err := app.Runner.RunCommands(ws, event); err != nil {
					log.Println("Error writing to stdin:", err)
				}
			}()
		default:
			log.Printf("Unknown event type: %s", event.Type)
		}
	}
}
