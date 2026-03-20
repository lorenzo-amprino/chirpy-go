package main

import (
	"chirpy/internal/auth"
	"net/http"
	"time"
)

func (cfg *apiConfig) handleRefreshToken(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		responseWithError(w, 401, "Unauthorized")
		return
	}

	refreshToken, err := cfg.queries.GetRefreshToken(r.Context(), token)

	if err != nil || refreshToken.ExpiresAt.Before(time.Now()) || refreshToken.RevokedAt.Valid {
		responseWithError(w, 401, "Unauthorized")
		return
	}

	userID := refreshToken.UserID

	expiration := time.Duration(1) * time.Hour

	authToken, err := auth.MakeJWT(userID, cfg.jwtSecret, expiration)
	if err != nil {
		responseWithError(w, 500, "Something went wrong")
		return
	}

	type response struct {
		Token string `json:"token"`
	}

	responseWithJson(w, 200, response{Token: authToken})
}

func (cfg *apiConfig) handleRevokeToken(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		responseWithError(w, 401, "Unauthorized")
		return
	}

	refreshToken, err := cfg.queries.GetRefreshToken(r.Context(), token)
	if err != nil || refreshToken.ExpiresAt.Before(time.Now()) || refreshToken.RevokedAt.Valid {
		responseWithError(w, 401, "Unauthorized")
		return
	}

	_, err = cfg.queries.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		responseWithError(w, 500, "Something went wrong")
		return
	}

	w.WriteHeader(204)
}
