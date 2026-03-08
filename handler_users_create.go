package main

import (
	"encoding/json"
	"net/http"
)

type User struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Email     string `json:"email"`
}

func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	req := struct {
		Email string `json:"email"`
	}{}

	type response struct {
		User
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		responseWithError(w, "Something went wrong")
		return
	}

	user, err := cfg.queries.CreateUser(r.Context(), req.Email)
	if err != nil {
		responseWithError(w, "Something went wrong")
		return
	}

	responseWithJson(w, 201, response{User: User{
		ID:        user.ID.String(),
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
		Email:     user.Email,
	}})

}
