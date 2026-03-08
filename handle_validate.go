package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handleValidateChirp(w http.ResponseWriter, r *http.Request) {

	req := struct {
		Body string `json:"body"`
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

	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	responseWithJson(w, 200, returnVals{CleanedBody: cleanMessage(req.Body)})

}

func cleanMessage(message string) string {
	s := strings.Split(message, " ")
	for i, word := range s {
		if containsBannedWord(word) {
			s[i] = "****"
		}
	}
	return strings.Join(s, " ")

}

func containsBannedWord(word string) bool {
	for _, banned := range BANNED_WORDS {
		if strings.EqualFold(word, banned) {
			return true
		}
	}
	return false
}

var BANNED_WORDS = []string{"kerfuffle", "sharbert", "fornax"}
