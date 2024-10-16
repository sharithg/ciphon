package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sharithg/siphon/internal/repository"
	"github.com/sharithg/siphon/internal/workflow"
)

func (app *Application) getWorkflows(w http.ResponseWriter, r *http.Request) {
	workflows, err := app.Repository.GetWorkflowRuns(r.Context())

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, workflows); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *Application) getJobs(w http.ResponseWriter, r *http.Request) {
	workflowId := chi.URLParam(r, "workflowId")

	id, err := uuid.Parse(workflowId)

	if err != nil {
		app.badRequestResponse(w, r, errors.New("invalid workflow id"))
		return
	}

	jobs, err := app.Repository.GetJobsByWorkflowId(r.Context(), id)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, jobs); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *Application) getSteps(w http.ResponseWriter, r *http.Request) {
	jobId := chi.URLParam(r, "jobId")
	id, err := uuid.Parse(jobId)

	if err != nil {
		app.badRequestResponse(w, r, errors.New("invalid job id"))
		return
	}

	steps, err := app.Repository.GetStepsByJobId(r.Context(), id)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, steps); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *Application) getStepOutput(w http.ResponseWriter, r *http.Request) {
	stepId := chi.URLParam(r, "stepId")

	id, err := uuid.Parse(stepId)

	if err != nil {
		app.badRequestResponse(w, r, errors.New("invalid step id"))
		return
	}

	stepOutputs, err := app.Repository.GetCommandOutputsByStepId(r.Context(), id)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, stepOutputs); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *Application) triggerWorkflow(w http.ResponseWriter, r *http.Request) {
	workflowId := chi.URLParam(r, "workflowId")

	id, err := uuid.Parse(workflowId)

	if err != nil {
		app.badRequestResponse(w, r, errors.New("invalid workflow id"))
		return
	}

	wm := workflow.New(app.Repository, app.Config.Github.AppConfig.OAuth.ClientID, app.Cache)

	if err := app.Repository.ResetWorkflowRun(r.Context(), app.Pool, id); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.updateWorkflowStatus(r.Context(), id, "running"); err != nil {
		slog.Error("error updating workflow status: %w", "err", err)
	}

	go func() {
		ctx := context.Background()
		start := time.Now()

		err := wm.TriggerWorkflow(ctx, id)
		if err != nil {
			if err := app.updateWorkflowStatusWithDuration(ctx, id, "failed", start); err != nil {
				slog.Error("error updating workflow status: %w", "err", err)
			}
		}

		if err := app.updateWorkflowStatusWithDuration(ctx, id, "success", start); err != nil {
			slog.Error("error updating workflow status: %w", "err", err)
		}
	}()
}

func (app *Application) updateWorkflowStatus(ctx context.Context, workflowId uuid.UUID, status string) error {
	eventPayload := workflow.WorkflowRun{Id: workflowId, Status: status, Type: "workflow"}
	if err := app.Cache.Publish(ctx, "workflow_run", eventPayload).Err(); err != nil {
		return err
	}
	if err := app.Repository.UpdateWorkflowRunStatus(ctx, repository.UpdateWorkflowRunStatusParams{
		ID:     workflowId,
		Status: &status,
	}); err != nil {
		return err
	}
	return nil
}

func (app *Application) updateWorkflowStatusWithDuration(ctx context.Context, workflowId uuid.UUID, status string, start time.Time) error {
	duration := time.Since(start)
	secs := duration.Seconds()
	eventPayload := workflow.WorkflowRun{Id: workflowId, Status: status, Type: "workflow"}

	if err := app.Cache.Publish(ctx, "workflow_run", eventPayload).Err(); err != nil {
		return err
	}
	if err := app.updateWorkflowStatus(ctx, workflowId, status); err != nil {
		return err
	}
	if err := app.Repository.UpdateWorkflowRunDuration(ctx, repository.UpdateWorkflowRunDurationParams{
		Duration: &secs,
		ID:       workflowId,
	}); err != nil {
		return err
	}
	return nil
}
