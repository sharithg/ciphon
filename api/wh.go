package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/google/go-github/v65/github"
	"github.com/google/uuid"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/pkg/errors"
	"github.com/sharithg/siphon/internal/parser"
	"github.com/sharithg/siphon/internal/repository"
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

	slog.Info("handling github webhook event", "event", eventType)

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

	go h.handlePushEvent(event, config)

	return nil
}

func (h *GhWebhookHandler) handlePushEvent(event github.PushEvent, configStr string) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipelineId, err := h.createPipelineRun(ctx, event, configStr)

	if err != nil {
		slog.Error("error creating pipeline run", "error", err)
		return
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

	for name := range config.Workflows {
		workflowRun := repository.CreateWorkflowRunParams{
			Name:          name,
			PipelineRunID: pipelineId,
		}

		workflowId, err := h.app.Repository.CreateWorkflowRun(ctx, workflowRun)

		if err != nil {
			slog.Error("error creating workflow", "error", err)
			return
		}

		jobs, err := config.GetWorkflowJobs(name)

		if err != nil {
			slog.Error("error getting jobs for workflow", "error", err)
			return
		}

		var jobRuns []repository.CreateJobRunParams
		for _, job := range jobs {
			id := uuid.New()
			jobRun := repository.CreateJobRunParams{
				ID:         id,
				WorkflowID: workflowId,
				Name:       job.Name,
				Docker:     job.Docker,
				Node:       &job.Node,
			}
			jobRuns = append(jobRuns, jobRun)
		}

		for _, job := range jobs {
			requires := config.GetJobRequires(job.Name)

			reqIds := getRequiredJobs(requires, jobRuns)

			var jobRun repository.CreateJobRunParams

			for _, run := range jobRuns {
				if run.Name == job.Name {
					jobRun = run
				}
			}

			jobRun.Requires = reqIds

			jobId, err := h.app.Repository.CreateJobRun(ctx, jobRun)

			if err != nil {
				slog.Error("error creating job", "error", err)
				return
			}

			h.createSteps(ctx, job.Steps, jobId)
		}
	}
}

func (h *GhWebhookHandler) createSteps(ctx context.Context, steps []parser.StepWrapper, jobId uuid.UUID) {
	for idx, step := range steps {
		var restoreCache []string
		var paths []string

		if step.Step.RestoreCache != nil {
			restoreCache = step.Step.RestoreCache.Keys
		}
		if step.Step.SaveCache != nil {
			paths = step.Step.SaveCache.Paths
		}

		stepOrder := int32(idx)
		stepRun := repository.CreateStepRunParams{
			JobID:     jobId,
			Type:      step.Step.Type,
			Name:      &step.Step.Name,
			Command:   &step.Step.Command,
			Keys:      restoreCache,
			Paths:     paths,
			StepOrder: stepOrder,
		}

		_, err := h.app.Repository.CreateStepRun(ctx, stepRun)

		if err != nil {
			slog.Error("error creating step", "error", err)
			return
		}
	}
}

func getRequiredJobs(requires []string, jobRuns []repository.CreateJobRunParams) []uuid.UUID {
	var reqIds []uuid.UUID
	for _, req := range requires {
		for _, jobReq := range jobRuns {
			if jobReq.Name == req {
				reqIds = append(reqIds, jobReq.ID)
			}
		}
	}
	return reqIds
}

func (h *GhWebhookHandler) createPipelineRun(ctx context.Context, event github.PushEvent, configStr string) (uuid.UUID, error) {

	headCommit := event.HeadCommit

	if headCommit == nil {
		return uuid.Nil, errors.New("error parsing event, HeadCommit is nil")
	}
	if headCommit.ID == nil {
		return uuid.Nil, errors.New("headCommit.ID is nil")
	}
	if event.Ref == nil {
		return uuid.Nil, errors.New("event.Ref is nil")
	}
	if event.Repo == nil || event.Repo.ID == nil {
		return uuid.Nil, errors.New("event.Repo.ID is nil")
	}

	pipelineRun := repository.CreatePipelineRunParams{
		CommitSha:  *headCommit.ID,
		ConfigFile: configStr,
		Branch:     strings.Replace(*event.Ref, "refs/heads/", "", -1),
		Status:     "received",
		RepoID:     *event.Repo.ID,
	}

	pipelineId, err := h.app.Repository.CreatePipelineRun(ctx, pipelineRun)

	if err != nil {
		return uuid.Nil, fmt.Errorf("error creating pipeline run: %s", err)
	}

	return pipelineId, nil
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
