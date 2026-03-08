package main

import (
	"chirpy/internal/database"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (c *apiConfig) handleCreateChirp(w http.ResponseWriter, r *http.Request) {

	req := struct {
		Body   string `json:"body"`
		UserID string `json:"user_id"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		responseWithError(w, "Something went wrong")
		return
	}
	if len(req.Body) > 140 {
		responseWithError(w, "Chirp is too long")
		return
	}
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		responseWithError(w, "Invalid user ID")
		return
	}

	cleaned := cleanMessage(req.Body)

	chirp, err := c.queries.CreateChirp(r.Context(), database.CreateChirpParams{Body: cleaned, UserID: userID})
	if err != nil {
		responseWithError(w, "Something went wrong")
		return
	}

	type returnVals struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt string    `json:"created_at"`
		UpdatedAt string    `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	responseWithJson(w, 201, returnVals{
		Id:        chirp.ID,
		CreatedAt: chirp.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: chirp.UpdatedAt.Format("2006-01-02 15:04:05"),
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}
