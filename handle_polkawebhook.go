package main

import (
	"chirpy/internal/auth"
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type request struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) handlePolkaWebhook(w http.ResponseWriter, r *http.Request) {

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil || apiKey != cfg.polka_api_key {
		responseWithError(w, 401, "Unathorized")
		return
	}

	var req request
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		responseWithError(w, 400, "Bad Request")
		return
	}

	if req.Event != "user.upgraded" {
		responseWithJson(w, 204, "")
	}

	_, err = cfg.queries.UpgradeToChirpyRed(context.Background(), uuid.MustParse(req.Data.UserID))
	if err != nil {
		responseWithError(w, 404, "user not found")
	}

	responseWithJson(w, 204, "")

}
