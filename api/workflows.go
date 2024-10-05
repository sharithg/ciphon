package api

import (
	"net/http"

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

func (app *Application) triggerWorkflow(w http.ResponseWriter, r *http.Request) {
	workflowId := chi.URLParam(r, "workflowId")

	wm := workflow.New(app.Store, app.Config.Github.AppConfig.OAuth.ClientID)

	go wm.TriggerWorkflow(r.Context(), workflowId)
}
