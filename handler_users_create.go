package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"encoding/json"
	"net/http"
)

type User struct {
	ID           string `json:"id"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	Email        string `json:"email"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	IsChirpyRed  bool   `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	req := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	type response struct {
		User
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		responseWithError(w, 400, "Something went wrong")
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		responseWithError(w, 400, "Something went wrong")
		return
	}

	user, err := cfg.queries.CreateUser(r.Context(), database.CreateUserParams{
		Email:          req.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		responseWithError(w, 400, "Something went wrong")
		return
	}

	responseWithJson(w, 201, response{User: User{
		ID:          user.ID.String(),
		CreatedAt:   user.CreatedAt.String(),
		UpdatedAt:   user.UpdatedAt.String(),
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed.Bool,
	}})

}
