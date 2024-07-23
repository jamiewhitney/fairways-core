package transport

import (
	"encoding/json"
	"net/http"
)

func ResponseWithJSON(w http.ResponseWriter, code int, payload interface{}) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		return
	}

	return
}

func RespondWithError(w http.ResponseWriter, code int, err error) {
	ResponseWithJSON(w, code, map[string]string{"error": err.Error()})
}
