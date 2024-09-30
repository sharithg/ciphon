package handlers

import "net/http"

func (env *Env) HandleGhWebhook(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Node added successfully"))
}
