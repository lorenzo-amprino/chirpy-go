package main

import (
	"chirpy/internal/auth"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleDeleteChirp(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		responseWithError(w, 401, "Unauthorized")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		responseWithError(w, 401, "Unauthorized")
		return
	}

	id := r.PathValue("id")
	if id == "" {
		responseWithError(w, 400, "Missing id")
		return
	}

	chirp, err := cfg.queries.GetChirpByID(r.Context(), uuid.MustParse(id))
	if err != nil {
		responseWithError(w, 404, "Chirp not found")
		return
	}

	if chirp.UserID != userID {
		responseWithError(w, 403, "Forbidden")
		return
	}

	err = cfg.queries.DeleteChirpByID(r.Context(), uuid.MustParse(id))
	if err != nil {
		responseWithError(w, 500, "Internal Server Error")
		return
	}

	responseWithJson(w, 204, map[string]string{"message": "Chirp deleted"})

}
