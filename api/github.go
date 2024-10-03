package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/go-github/v65/github"
)

type GithubRepo struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	LastUpdated string `json:"lastUpdated"`
}

func (app *Application) getReposHandler(w http.ResponseWriter, r *http.Request) {
	opt := &github.ListOptions{}
	ctx := context.Background()
	repos, _, err := app.GithubClient.Apps.ListRepos(ctx, opt)

	// fmt.Println(resp.Request.Header)
	if err != nil {
		log.Printf("Error getting github repos: %v", err)
		http.Error(w, "Error getting github repos", http.StatusInternalServerError)
		return
	}

	var repoList []GithubRepo
	for _, repo := range repos.Repositories {
		repoList = append(repoList, GithubRepo{
			ID:          repo.GetID(),
			Name:        repo.GetName(),
			Description: repo.GetDescription(),
			LastUpdated: repo.GetUpdatedAt().String(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(repoList); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
