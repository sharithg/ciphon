package api

import (
	"log/slog"
	"net/http"
)

func (app *Application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error("internal error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *Application) forbiddenResponse(w http.ResponseWriter, r *http.Request, err error) {
	slog.Warn("forbidden", "method", r.Method, "path", r.URL.Path, "error", err)

	writeJSONError(w, http.StatusForbidden, "forbidden")
}

func (app *Application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	slog.Warn("bad request", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *Application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error("conflict response", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJSONError(w, http.StatusConflict, err.Error())
}

func (app *Application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	slog.Warn("not found error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJSONError(w, http.StatusNotFound, "not found")
}

func (app *Application) unauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	slog.Warn("unauthorized error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJSONError(w, http.StatusUnauthorized, err.Error())
}

func (app *Application) unauthorizedBasicErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	slog.Warn("unauthorized basic error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)

	writeJSONError(w, http.StatusUnauthorized, "unauthorized")
}

func (app *Application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request, retryAfter string) {
	slog.Warn("rate limit exceeded", "method", r.Method, "path", r.URL.Path)

	w.Header().Set("Retry-After", retryAfter)

	writeJSONError(w, http.StatusTooManyRequests, "rate limit exceeded, retry after: "+retryAfter)
}
