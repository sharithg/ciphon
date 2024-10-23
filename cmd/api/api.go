package main

import (
	"log"
	"log/slog"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/sharithg/siphon/api"
	"github.com/sharithg/siphon/internal/auth"
	"github.com/sharithg/siphon/internal/db"
	"github.com/sharithg/siphon/internal/env"
	"github.com/sharithg/siphon/internal/repo"
	"github.com/sharithg/siphon/internal/repository"
	"github.com/sharithg/siphon/internal/service"
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
		Cache: api.CacheConfig{
			Addr: env.GetString("CACHE_URL", false, "localhost:6379"),
		},
	}

	pool, err := db.New(cfg.Db.Addr)

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	if err != nil {
		log.Fatal("error configuring db: ", err)
	}

	defer pool.Close()

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

	repository := repository.New(pool)
	minioStorage := minio.NewStorage(minioClient)

	if err = minioStorage.SetupBuckets(); err != nil {
		log.Fatal("error setting up minio buckets", err)
	}

	auth := auth.New(env.GetString("JWT_SECRET_KEY", true, ""), time.Hour*24, time.Hour*24*7)

	service := service.NewService(repository)

	app := &api.Application{
		Config:       cfg,
		Repository:   repository,
		MinioStorage: minioStorage,
		Github:       ghClient,
		Cache:        redisClient,
		Auth:         auth,
		Pool:         pool,
		Service:      service,
	}

	mux := app.Mount()

	log.Fatal(app.Run(mux))
}
