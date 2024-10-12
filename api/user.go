package api

import (
	"errors"
	"fmt"
	"net/http"
)

func (app *Application) getUser(w http.ResponseWriter, r *http.Request) {
	userId, ok := app.getUserFromContext(r)

	if !ok {
		app.badRequestResponse(w, r, errors.New("user not found in request"))
		return
	}

	user, err := app.Store.UsersStore.GetById(r.Context(), userId)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if user == nil {
		app.badRequestResponse(w, r, fmt.Errorf("user not found for id: %s", userId))
		return
	}

	app.jsonResponse(w, http.StatusOK, user)
}
