package api

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (app *Application) parseUUIDParam(w http.ResponseWriter, r *http.Request, paramName string) (uuid.UUID, bool) {
	paramValue := chi.URLParam(r, paramName)
	id, err := uuid.Parse(paramValue)
	if err != nil {
		app.badRequestResponse(w, r, errors.New("invalid "+paramName))
		return uuid.Nil, false
	}
	return id, true
}
