package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sharithg/siphon/internal/repository"
	"github.com/sharithg/siphon/internal/service"
	"github.com/sharithg/siphon/internal/workflow"
)

type JobResponse struct {
	Jobs  []repository.GetJobsByWorkflowIdRow `json:"jobs"`
	Edges []service.Edge                      `json:"edges"`
}

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

	id, ok := app.parseUUIDParam(w, r, "workflowId")

	if !ok {
		app.badRequestResponse(w, r, errors.New("invalid workflow id"))
		return
	}

	jobs, dag, err := app.Service.Job.GetByWorkflowId(r.Context(), id)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	resp := JobResponse{
		Jobs:  jobs,
		Edges: dag.GetEdges(),
	}

	if err := app.jsonResponse(w, http.StatusOK, resp); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *Application) getSteps(w http.ResponseWriter, r *http.Request) {
	id, ok := app.parseUUIDParam(w, r, "jobId")

	if !ok {
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

	id, ok := app.parseUUIDParam(w, r, "stepId")

	if !ok {
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
	id, ok := app.parseUUIDParam(w, r, "workflowId")

	if !ok {
		return
	}

	wm := workflow.New(app.Repository, app.Config.Github.AppConfig.OAuth.ClientID, app.Cache, app.Service)

	if err := app.Repository.ResetWorkflowRun(r.Context(), app.Pool, id); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.updateWorkflowStatus(r.Context(), id, "running"); err != nil {
		slog.Error("error updating workflow status", "err", err)
		return
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Hour)
		defer cancel()

		start := time.Now()

		err := wm.TriggerWorkflow(ctx, id)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				slog.Warn("workflow timed out", "id", id)
				if err := app.updateWorkflowStatusWithDuration(ctx, id, "timeout", start); err != nil {
					slog.Error("error updating workflow status", "err", err)
				}
				return
			}
			if err := app.updateWorkflowStatusWithDuration(ctx, id, "failed", start); err != nil {
				slog.Error("error updating workflow status", "err", err)
			}
		} else {
			if err := app.updateWorkflowStatusWithDuration(ctx, id, "success", start); err != nil {
				slog.Error("error updating workflow status", "err", err)
			}
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
