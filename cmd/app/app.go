package main

import (
	"log"
	"log/slog"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sharithg/siphon/api"
	"github.com/sharithg/siphon/internal/db"
	"github.com/sharithg/siphon/internal/env"
	"github.com/sharithg/siphon/internal/repo"
	"github.com/sharithg/siphon/internal/storage"
	"github.com/sharithg/siphon/internal/storage/minio"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ghCfg, err := repo.ReadGithubConfig()

	if err != nil {
		log.Fatal("error reading github config: ", err)
	}

	cfg := api.Config{
		Addr: env.GetString("ADDR", false, ":8000"),
		Db: api.DbConfig{
			Addr:         env.GetString("PG_URL", true, ""),
			MaxOpenConns: env.GetInt("PG_MAX_OPEN_CONNS", false, 30),
			MaxIdleConns: env.GetInt("PG_MAX_IDLE_CONNS", false, 30),
			MaxIdleTime:  env.GetString("PG_MAX_IDLE_TIME", false, "15m"),
		},
		Github: *ghCfg,
		Env:    env.GetString("GOENV", false, "local"),
	}

	dbClient, err := db.New(
		cfg.Db.Addr,
		cfg.Db.MaxOpenConns,
		cfg.Db.MaxIdleConns,
		cfg.Db.MaxIdleTime,
	)

	if err != nil {
		log.Fatal("error configuring db", err)
	}

	defer dbClient.Close()

	if err = db.Migrate(cfg.Db.Addr); err != nil {
		slog.Warn("running migrations", "msg", err)
	}

	minioClient, err := minio.New()

	if err != nil {
		log.Fatal("error configuring minio", err)
	}

	ghClient, err := repo.New(cfg.Github)

	if err != nil {
		log.Fatal("error configuring github client", err)
	}

	store := storage.NewStorage(dbClient)
	minioStorage := minio.NewStorage(minioClient)

	if err = minioStorage.SetupBuckets(); err != nil {
		log.Fatal("error setting up minio buckets", err)
	}

	app := &api.Application{
		Config:       cfg,
		Store:        store,
		MinioStorage: minioStorage,
		Github:       ghClient,
	}

	mux := app.Mount()

	log.Fatal(app.Run(mux))
}
