package api

import "net/http"

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
