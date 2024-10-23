package main

import (
	"log"
	"log/slog"
	"os"
	"path"

	"github.com/joho/godotenv"
	"github.com/sharithg/siphon/agent"
	"github.com/sharithg/siphon/agent/cli"
	"github.com/sharithg/siphon/agent/docker"
	"github.com/sharithg/siphon/internal/config"
	"github.com/sharithg/siphon/internal/env"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		slog.Warn("error loading .env file", "err", err)
	}

	cfg := agent.Config{
		Addr: env.GetString("AGENT_ADDR", false, ":8888"),
		Env:  env.GetString("GOENV", false, "local"),
	}

	cli, err := cli.New()

	if err != nil {
		log.Fatalf("error creating docker client: %s", err)
	}

	dock := docker.New()

	homeDir, err := os.UserHomeDir()

	if err != nil {
		log.Fatalf("error getting user home dir: %s", err)
	}

	defaultConfig := path.Join(homeDir, ".ciphon/agent.yaml")

	agentConfig, err := config.LoadAgentConfig(env.GetString("AGENT_CONFIG_PATH", false, defaultConfig))

	if err != nil {
		log.Fatalf("error loading agent config: %s", err)
	}

	app := &agent.Application{
		Config:      cfg,
		Cli:         cli,
		AgentConfig: agentConfig,
		Docker:      dock,
	}

	log.Fatal(app.Run())
}
