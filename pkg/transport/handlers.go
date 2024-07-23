package transport

import (
	"net/http"
)

func healthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		health := struct {
			healthy bool
		}{}

		health.healthy = true

		ResponseWithJSON(w, http.StatusOK, health)
	}
}
