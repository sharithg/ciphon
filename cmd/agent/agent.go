package main

import (
	"log"
	"log/slog"
	"os"
	"path"

	"github.com/joho/godotenv"
	"github.com/sharithg/siphon/internal/config"
	"github.com/sharithg/siphon/internal/docker"
	"github.com/sharithg/siphon/internal/env"
	storage "github.com/sharithg/siphon/internal/storage/kv"
	"github.com/sharithg/siphon/ws"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		slog.Warn("error loading .env file", "err", err)
	}

	cfg := ws.Config{
		Addr: env.GetString("AGENT_ADDR", false, ":8888"),
		Env:  env.GetString("GOENV", false, "local"),
	}

	store := storage.NewKvStorage()

	dock, err := docker.New()

	if err != nil {
		log.Fatalf("error creating docker client: %s", err)
	}

	homeDir, err := os.UserHomeDir()

	if err != nil {
		log.Fatalf("error getting user home dir: %s", err)
	}

	defaultConfig := path.Join(homeDir, ".ciphon/agent.yaml")

	agentConfig, err := config.LoadAgentConfig(env.GetString("AGENT_CONFIG_PATH", false, defaultConfig))

	if err != nil {
		log.Fatalf("error loading agent config: %s", err)
	}

	app := &ws.Application{
		Config:      cfg,
		Store:       store,
		Docker:      dock,
		AgentConfig: agentConfig,
	}

	log.Fatal(app.Run())
}
