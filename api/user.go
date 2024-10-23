package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (app *Application) getUser(w http.ResponseWriter, r *http.Request) {
	userIdStr, ok := app.getUserFromContext(r)

	if !ok {
		app.badRequestResponse(w, r, errors.New("user not found in request"))
		return
	}

	userId, err := uuid.Parse(userIdStr)

	if err != nil {
		app.badRequestResponse(w, r, errors.New("invalid user id"))
		return
	}

	user, err := app.Repository.GetUserById(r.Context(), userId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.badRequestResponse(w, r, fmt.Errorf("user not found for id: %s", userId))
			return
		}

		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, user)
}
