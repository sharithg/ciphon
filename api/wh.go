package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/google/go-github/v65/github"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/pkg/errors"
	"github.com/sharithg/siphon/internal/parser"
	"github.com/sharithg/siphon/internal/storage"
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

	headCommit := event.HeadCommit

	if headCommit == nil {
		return errorWithLog("head commit null")
	}

	opts := &github.RepositoryContentGetOptions{
		Ref: *headCommit.ID,
	}

	rc, _, err := h.app.Github.Client.Repositories.DownloadContents(ctx, *owner, *repo, ".siphon/pipeline.yaml", opts)

	if err != nil {
		return wrappedErrorWithLog(err, "failed to read config file")
	}

	config, err := readCloserToString(rc)

	if err != nil {
		return wrappedErrorWithLog(err, "failed to read config contents")
	}

	go h.handlePushEvent(ctx, event, config)

	return nil
}

func (h *GhWebhookHandler) handlePushEvent(ctx context.Context, event github.PushEvent, configStr string) {

	headCommit := event.HeadCommit

	if headCommit == nil {
		slog.Error("error parsing event, HeadCommit is nil")
		return
	}
	pipelineRun := storage.PipelineRun{
		CommitSHA:  *headCommit.ID,
		ConfigFile: configStr,
		Branch:     strings.Replace(*event.Ref, "refs/heads/", "", -1),
		Status:     "received",
		RepoId:     *event.Repo.ID,
	}

	config, err := parser.ParseConfig(configStr)

	if err != nil {
		slog.Error("error parsing config", "error", err)
		return
	}

	if err = config.ValidateWorkflows(); err != nil {
		slog.Error("error validating config", "error", err)
		return
	}

	pipelineId, err := h.app.Store.PipelineRunsStore.Create(ctx, pipelineRun)

	if err != nil {
		slog.Error("error creating pipeline run", "error", err)
		return
	}

	for name := range config.Workflows {
		workflowRun := storage.WorkflowRun{
			Name:          name,
			PipelineRunID: pipelineId,
		}

		workflowId, err := h.app.Store.WorkflowRunsStore.Create(ctx, workflowRun)

		if err != nil {
			slog.Error("error creating workflow", "error", err)
			return
		}

		jobs, err := config.GetWorkflowJobs(name)

		if err != nil {
			slog.Error("error getting jobs for workflow", "error", err)
			return
		}

		for _, job := range jobs {
			jobRun := storage.JobRun{
				WorkflowID: workflowId,
				Name:       job.Name,
				Docker:     job.Docker,
				Node:       job.Node,
			}
			jobId, err := h.app.Store.JobRunsStore.Create(ctx, jobRun)

			if err != nil {
				slog.Error("error creating job", "error", err)
				return
			}

			for idx, step := range job.Steps {
				var restoreCache []string
				var paths []string

				if step.Step.RestoreCache != nil {
					restoreCache = step.Step.RestoreCache.Keys
				}
				if step.Step.SaveCache != nil {
					paths = step.Step.SaveCache.Paths
				}

				stepRun := storage.StepRun{
					JobID:     jobId,
					Type:      step.Step.Type,
					Name:      step.Step.Name,
					Command:   step.Step.Command,
					Keys:      restoreCache,
					Paths:     paths,
					StepOrder: idx,
				}

				_, err := h.app.Store.StepRunsStore.Create(ctx, stepRun)

				if err != nil {
					slog.Error("error creating step", "error", err)
					return
				}
			}
		}

	}
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
