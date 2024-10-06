package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sharithg/siphon/internal/workflow"
)

func (app *Application) getWorkflows(w http.ResponseWriter, r *http.Request) {
	workflows, err := app.Store.WorkflowRunsStore.GetWorkflowRuns()

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

	jobs, err := app.Store.JobRunsStore.GetByWorkflowId(workflowId)

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

	jobs, err := app.Store.StepRunsStore.GetByJobId(jobId)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, jobs); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *Application) triggerWorkflow(w http.ResponseWriter, r *http.Request) {
	workflowId := chi.URLParam(r, "workflowId")

	wm := workflow.New(app.Store, app.Config.Github.AppConfig.OAuth.ClientID, app.Cache)

	if err := app.Store.WorkflowRunsStore.UpdateAllStatuses(workflowId); err != nil {
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

func (app *Application) eventsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	ctx := r.Context()

	subscriber := app.Cache.Subscribe(ctx, "workflow_run", "job_run", "step_run")

	for {
		select {
		case <-ctx.Done():
			// If the client disconnects, stop sending events.
			fmt.Println("Client disconnected")
			return
		default:
			msg, err := subscriber.ReceiveMessage(ctx)

			if err != nil {
				http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
				return
			}

			fmt.Println("Recived message: ", msg.Payload)

			fmt.Fprintf(w, "data: %s\n\n", msg.Payload)
			flusher.Flush()
		}
	}
}

func (app *Application) updateWorkflowStatus(ctx context.Context, workflowId, status string) error {
	eventPayload := workflow.WorkflowRun{Id: workflowId, Status: status, Type: "workflow"}
	if err := app.Cache.Publish(ctx, "workflow_run", eventPayload).Err(); err != nil {
		return err
	}
	if err := app.Store.WorkflowRunsStore.UpdateStatus(workflowId, status); err != nil {
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
	if err := app.Store.WorkflowRunsStore.UpdateDuration(workflowId, secs); err != nil {
		return err
	}
	return nil
}
