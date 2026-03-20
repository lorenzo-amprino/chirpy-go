package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"context"
	"encoding/json"
	"net/http"
	"sort"

	"github.com/google/uuid"
)

func (c *apiConfig) handleCreateChirp(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		responseWithError(w, 401, "Unauthorized")
		return
	}

	userID, err := auth.ValidateJWT(token, c.jwtSecret)
	if err != nil {
		responseWithError(w, 401, "Unauthorized")
		return
	}

	req := struct {
		Body   string `json:"body"`
		UserID string `json:"user_id"`
	}{}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		responseWithError(w, 400, "Something went wrong")
		return
	}
	if len(req.Body) > 140 {
		responseWithError(w, 400, "Chirp is too long")
		return
	}

	cleaned := cleanMessage(req.Body)

	chirp, err := c.queries.CreateChirp(r.Context(), database.CreateChirpParams{Body: cleaned, UserID: userID})
	if err != nil {
		responseWithError(w, 400, "Something went wrong")
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

func (c *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {

	authorId := r.URL.Query().Get("author_id")
	sorting := r.URL.Query().Get("sort")
	var chirps []database.Chirp
	var err error
	if authorId != "" {
		chirps, err = c.queries.GetChirpsByAuthorId(context.Background(), uuid.MustParse(authorId))
		if err != nil {
			responseWithError(w, 400, "Something went wrong")
			return
		}
	} else {
		chirps, err = c.queries.GetAllChirps(r.Context())
		if err != nil {
			responseWithError(w, 400, "Something went wrong")
			return
		}
	}

	if sorting == "desc" {
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.After(chirps[j].CreatedAt) })
	}
	if sorting == "asc" {
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.Before(chirps[j].CreatedAt) })
	}
	type returnVals struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt string    `json:"created_at"`
		UpdatedAt string    `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	var response []returnVals
	for _, chirp := range chirps {
		response = append(response, returnVals{
			Id:        chirp.ID,
			CreatedAt: chirp.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: chirp.UpdatedAt.Format("2006-01-02 15:04:05"),
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	responseWithJson(w, 200, response)
}

func (c *apiConfig) getChirpHandler(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	if id == "" {
		responseWithError(w, 400, "Missing chirp ID")
		return
	}
	uuid, err := uuid.Parse(id)
	if err != nil {
		responseWithError(w, 400, "Invalid chirp ID")
		return
	}

	chirp, err := c.queries.GetChirpByID(r.Context(), uuid)
	if err != nil {
		responseWithError(w, 404, "Chirp not found")
		return
	}

	type returnVals struct {
		Id        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Body      string `json:"body"`
		UserID    string `json:"user_id"`
	}

	responseWithJson(w, 200, returnVals{
		Id:        chirp.ID.String(),
		CreatedAt: chirp.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: chirp.UpdatedAt.Format("2006-01-02 15:04:05"),
		Body:      chirp.Body,
		UserID:    chirp.UserID.String(),
	})

}
