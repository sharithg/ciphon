package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"

	"github.com/google/go-github/v65/github"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/pkg/errors"
)

type GhWebhookHandler struct {
	githubapp.ClientCreator

	Preamble string

	app *Application
}

func NewGhWebhookHandler(cc githubapp.ClientCreator, preamble string, app *Application) *GhWebhookHandler {
	return &GhWebhookHandler{
		app:           app,
		ClientCreator: cc,
		Preamble:      preamble,
	}
}

func (h *GhWebhookHandler) Handles() []string {
	return []string{"push"}
}

func readCloserToString(rc io.ReadCloser) (string, error) {
	defer rc.Close()

	content, err := io.ReadAll(rc)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func (h *GhWebhookHandler) Handle(ctx context.Context, eventType, deliveryID string, payload []byte) error {
	var event github.PushEvent

	fmt.Println("handling event: ", eventType)
	if err := json.Unmarshal(payload, &event); err != nil {
		return wrappedErrorWithLog(err, "failed to parse issue comment event payload")
	}

	owner := event.Sender.Login
	repo := event.Repo.Name

	headCommit := event.Ref

	if headCommit == nil {
		return errorWithLog("head commit null")
	}

	opts := &github.RepositoryContentGetOptions{
		Ref: *headCommit,
	}

	rc, _, err := h.app.Github.Client.Repositories.DownloadContents(ctx, *owner, *repo, ".siphon/pipeline.yaml", opts)

	if err != nil {
		return wrappedErrorWithLog(err, "failed to read config file")
	}

	config, err := readCloserToString(rc)

	if err != nil {
		return wrappedErrorWithLog(err, "failed to read config contents")
	}

	fmt.Println(config)

	return nil
}

func wrappedErrorWithLog(err error, message string) error {
	e := errors.Wrap(err, message)
	slog.Error("error handling gh wh event", "error", e)
	return e
}

func errorWithLog(message string) error {
	e := errors.New(message)
	slog.Error("error handling gh wh event", "error", e)
	return e
}
