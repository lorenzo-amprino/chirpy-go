package main

import "net/http"

func (cfg *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)

	err := cfg.queries.ResetUsers(r.Context())
	if err != nil {
		responseWithError(w, "Something went wrong")
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}
