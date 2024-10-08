package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/go-github/v65/github"
	"github.com/sharithg/siphon/internal/storage"
)

type GithubRepoResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	LastUpdated string `json:"lastUpdated"`
	Owner       string `json:"owner"`
}

type ConnectRepoRequest struct {
	Name  string `json:"name"`
	Owner string `json:"owner"`
}

func (app *Application) getNewReposHandler(w http.ResponseWriter, r *http.Request) {
	opt := &github.ListOptions{}
	ctx := context.Background()
	repos, _, err := app.Github.Client.Apps.ListRepos(ctx, opt)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	existingRepos, err := app.Store.Repos.All()

	existingRepoSet := make(map[int64]struct{})
	for _, repo := range existingRepos {
		existingRepoSet[repo.RepoID] = struct{}{}
	}

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	var repoList []GithubRepoResponse
	for _, repo := range repos.Repositories {

		if _, exists := existingRepoSet[repo.GetID()]; exists {
			continue
		}

		repoList = append(repoList, GithubRepoResponse{
			ID:          repo.GetID(),
			Name:        repo.GetName(),
			Description: repo.GetDescription(),
			LastUpdated: repo.GetUpdatedAt().String(),
			Owner:       *repo.GetOwner().Login,
		})
	}

	if err := app.jsonResponse(w, http.StatusOK, repoList); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *Application) getReposHandler(w http.ResponseWriter, r *http.Request) {
	repos, err := app.Store.Repos.All()

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, repos); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *Application) connectRepoHandler(w http.ResponseWriter, r *http.Request) {
	var payload ConnectRepoRequest
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	repo, _, err := app.Github.Client.Repositories.Get(r.Context(), payload.Owner, payload.Name)

	if err != nil {
		app.badRequestResponse(w, r, errors.New("error fetching repo"))
		return
	}

	if repo == nil {
		app.notFoundResponse(w, r, errors.New("repo not found"))
		return
	}

	b, err := json.Marshal(repo)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	newRepo := storage.CreateRepo{
		RepoID:        *repo.ID,
		Name:          *repo.Name,
		Owner:         payload.Owner,
		Description:   *repo.Description,
		URL:           *repo.HTMLURL,
		RepoCreatedAt: repo.CreatedAt.Time,
		RawData:       string(b),
	}

	id, err := app.Store.Repos.Create(newRepo)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, id)
}
