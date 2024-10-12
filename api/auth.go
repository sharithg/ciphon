package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sharithg/siphon/internal/auth"
	"github.com/sharithg/siphon/internal/storage"
)

type RefreshToken struct {
	Token string `json:"token"`
}

func (app *Application) githubLoginHandler(w http.ResponseWriter, r *http.Request) {
	githubClientID := app.Config.Github.AppConfig.OAuth.ClientID

	redirectURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s",
		githubClientID,
		app.Config.Github.GithubCallbackUrl,
	)

	http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
}

func (app *Application) githubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	key := fmt.Sprintf("githubcode:%s", code)

	val, err := app.Cache.Get(r.Context(), key).Result()

	if err != redis.Nil && err != nil {
		app.internalServerError(w, r, err)
		return
	}

	var githubAccessToken string
	if err != redis.Nil && val != "" {
		githubAccessToken = val
	} else {
		githubAccessToken, err = app.Auth.GetGithubAccessToken(code, app.Config.Github.AppConfig.OAuth.ClientID, app.Config.Github.AppConfig.OAuth.ClientSecret)

		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		err = app.Cache.Set(r.Context(), key, githubAccessToken, 10*time.Second).Err()

		if err != nil {
			app.internalServerError(w, r, err)
			return
		}
	}

	githubData, err := app.Auth.GetGithubData(githubAccessToken)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.loggedinHandler(w, r, githubData)
}

func (app *Application) refreshTokens(w http.ResponseWriter, r *http.Request) {
	var payload RefreshToken
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	token, err := app.Auth.RefreshToken(payload.Token)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusCreated, token)
}

func (app *Application) loggedinHandler(w http.ResponseWriter, r *http.Request, githubData []byte) {
	if string(githubData) == "" {
		fmt.Fprintf(w, "UNAUTHORIZED!")
		return
	}

	user, err := app.createUserIfNotExists(r.Context(), githubData)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	tokens, err := app.Auth.CreateToken(user.ID)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, tokens)
}

func (app *Application) createUserIfNotExists(ctx context.Context, ghResp []byte) (*storage.User, error) {
	gu, err := auth.LoadGithubUser(ghResp)

	if err != nil {
		return nil, err
	}

	if gu.Login == "" {
		return nil, errors.New("user login not found")
	}

	externalId := strconv.Itoa(gu.ID)

	user := storage.User{
		Username:   gu.Login,
		Email:      gu.Email,
		ExternalId: externalId,
		AuthType:   "github",
	}

	gh := storage.GitHubUserInfo{
		Data: *gu,
	}

	existingUser, err := app.Store.UsersStore.GetByExternalId(ctx, externalId)

	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return existingUser, nil
	}

	newUser, err := app.Store.UsersStore.Create(ctx, user, gh)

	if err != nil {
		return nil, err
	}

	return newUser, nil
}
