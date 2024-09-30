package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/rs/cors"

	"github.com/sharithg/siphon/github"
	"github.com/sharithg/siphon/handlers"
	"github.com/sharithg/siphon/models"
	"github.com/sharithg/siphon/storage"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUrl := os.Getenv("PG_URL")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Failed to open db conn: ", err)
	}

	models.Migrate(dbUrl)

	client := storage.SetupMinio()

	env := &handlers.Env{
		Nodes:   models.NodeModel{DB: db},
		Storage: client,
	}

	router := http.NewServeMux()
	gh := github.New()

	router.Handle(githubapp.DefaultWebhookRoute, gh.Handler)

	router.HandleFunc("POST /node", env.AddNode)
	router.HandleFunc("GET /nodes", env.GetNodes)
	router.HandleFunc("GET /github/webooks", env.HandleGhWebhook)

	handler := cors.Default().Handler(router)

	if err = http.ListenAndServe(":8000", handler); err != nil {
		log.Fatal("error running server: ", err)
	}
}
