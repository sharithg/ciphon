package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/go-github/v65/github"
)

type GithubRepo struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (env *Env) GetRepos(w http.ResponseWriter, r *http.Request) {
	opt := &github.RepositoryListByUserOptions{}
	ctx := context.Background()
	repos, _, err := env.GhClient.Repositories.ListByUser(ctx, "sharithg", opt)

	if err != nil {
		log.Printf("Error getting github repos: %v", err)
		http.Error(w, "Error getting github repos", http.StatusInternalServerError)
		return
	}

	var repoList []GithubRepo
	for _, repo := range repos {
		repoList = append(repoList, GithubRepo{
			ID:   repo.GetID(),
			Name: repo.GetName(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(repoList); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
