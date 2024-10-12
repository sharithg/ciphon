package api

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sharithg/siphon/internal/workflow"
)

func (app *Application) getWorkflows(w http.ResponseWriter, r *http.Request) {
	workflows, err := app.Store.WorkflowRunsStore.GetWorkflowRuns(r.Context())

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

	jobs, err := app.Store.JobRunsStore.GetByWorkflowId(r.Context(), workflowId)

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

	steps, err := app.Store.StepRunsStore.GetByJobId(r.Context(), jobId)

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

	stepOutputs, err := app.Store.StepRunsStore.GetByStepID(r.Context(), stepId)

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

	wm := workflow.New(app.Store, app.Config.Github.AppConfig.OAuth.ClientID, app.Cache)

	if err := app.Store.WorkflowRunsStore.UpdateAllStatuses(r.Context(), workflowId); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.updateWorkflowStatus(r.Context(), workflowId, "running"); err != nil {
		slog.Error("error updating workflow status: %w", "err", err)
	}

	go func() {
		ctx := context.Background()
		start := time.Now()

		err := wm.TriggerWorkflow(ctx, workflowId)
		if err != nil {
			if err := app.updateWorkflowStatusWithDuration(ctx, workflowId, "failed", start); err != nil {
				slog.Error("error updating workflow status: %w", "err", err)
			}
		}

		if err := app.updateWorkflowStatusWithDuration(ctx, workflowId, "success", start); err != nil {
			slog.Error("error updating workflow status: %w", "err", err)
		}
	}()
}

func (app *Application) updateWorkflowStatus(ctx context.Context, workflowId, status string) error {
	eventPayload := workflow.WorkflowRun{Id: workflowId, Status: status, Type: "workflow"}
	if err := app.Cache.Publish(ctx, "workflow_run", eventPayload).Err(); err != nil {
		return err
	}
	if err := app.Store.WorkflowRunsStore.UpdateStatus(ctx, workflowId, status); err != nil {
		return err
	}
	return nil
}

func (app *Application) updateWorkflowStatusWithDuration(ctx context.Context, workflowId, status string, start time.Time) error {
	duration := time.Since(start)
	secs := duration.Seconds()
	eventPayload := workflow.WorkflowRun{Id: workflowId, Status: status, Type: "workflow"}

	if err := app.Cache.Publish(ctx, "workflow_run", eventPayload).Err(); err != nil {
		return err
	}
	if err := app.updateWorkflowStatus(ctx, workflowId, status); err != nil {
		return err
	}
	if err := app.Store.WorkflowRunsStore.UpdateDuration(ctx, workflowId, secs); err != nil {
		return err
	}
	return nil
}
