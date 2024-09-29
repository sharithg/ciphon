package handlers

import (
	"net/http"
	"slices"
	"strings"

	"github.com/sharithg/siphon/models"
	"github.com/sharithg/siphon/storage"
)

type Env struct {
	Nodes   models.NodeModel
	Storage *storage.Minio
}

var originAllowlist = []string{
	"http://127.0.0.1:9999",
	"http://cats.com",
	"http://safe.frontend.net",
}

var methodAllowlist = []string{"GET", "POST", "DELETE", "OPTIONS"}

func isPreflight(r *http.Request) bool {
	return r.Method == "OPTIONS" &&
		r.Header.Get("Origin") != "" &&
		r.Header.Get("Access-Control-Request-Method") != ""
}

func CheckCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isPreflight(r) {
			origin := r.Header.Get("Origin")
			method := r.Header.Get("Access-Control-Request-Method")
			if slices.Contains(originAllowlist, origin) && slices.Contains(methodAllowlist, method) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(methodAllowlist, ", "))
				w.Header().Add("Vary", "Origin")
			}
		} else {
			origin := r.Header.Get("Origin")
			if slices.Contains(originAllowlist, origin) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Add("Vary", "Origin")
			}
		}
		next.ServeHTTP(w, r)
	})
}
