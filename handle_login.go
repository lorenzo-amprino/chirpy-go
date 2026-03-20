package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"encoding/json"
	"net/http"
	"time"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	req := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		responseWithError(w, 400, "Something went wrong")
		return
	}

	user, err := cfg.queries.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		responseWithError(w, 401, "Incorrect email or password")
		return
	}

	if ok, err := auth.CheckPasswordHash(req.Password, user.HashedPassword); err != nil || !ok {
		responseWithError(w, 401, "Incorrect email or password")
		return
	}

	type response struct {
		User
	}

	expiration := time.Duration(1) * time.Hour

	if expiration > time.Duration(1)*time.Hour {
		expiration = time.Duration(1) * time.Hour
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expiration)
	refreshToken, err := auth.MakeRefreshToken()

	cfg.queries.SaveRefreshToken(r.Context(), database.SaveRefreshTokenParams{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	})

	responseWithJson(w, 200, response{User: User{
		ID:           user.ID.String(),
		CreatedAt:    user.CreatedAt.String(),
		UpdatedAt:    user.UpdatedAt.String(),
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken,
		IsChirpyRed:  user.IsChirpyRed.Bool,
	}})
}
