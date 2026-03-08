package main

import (
	"encoding/json"
	"net/http"
)

type ChirpValidateResponse struct {
	Error string `json:"error"`
	Valid bool   `json:"valid"`
}

func responseWithError(w http.ResponseWriter, message string) {
	res := ChirpValidateResponse{
		Error: message,
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(400)
	json.NewEncoder(w).Encode(res)
}

func responseWithJson(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}
