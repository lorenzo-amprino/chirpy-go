package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"context"
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handleUpdateUser(w http.ResponseWriter, r *http.Request) {

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

	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req request
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		responseWithError(w, 400, "Bad Request")
		return
	}

	if req.Email == "" || req.Password == "" {
		responseWithError(w, 400, "Email and password are required")
		return
	}

	p, err := auth.HashPassword(req.Password)

	updated, err := cfg.queries.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             userID,
		Email:          req.Email,
		HashedPassword: p,
	})
	if err != nil {
		responseWithError(w, 400, "Something went wrong")
		return
	}

	user, err := cfg.queries.GetUserByEmail(context.Background(), updated.Email)
	if err != nil {
		responseWithError(w, 400, "Something went wrong")
		return
	}

	type response struct {
		ID          string `json:"id"`
		CreatedAt   string `json:"created_at"`
		UpdatedAt   string `json:"updated_at"`
		Email       string `json:"email"`
		IsChirpyRed bool   `json:"is_chirpy_red"`
	}

	responseWithJson(w, 200, response{
		ID:          userID.String(),
		CreatedAt:   user.CreatedAt.String(),
		UpdatedAt:   user.UpdatedAt.String(),
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed.Bool,
	})

}
