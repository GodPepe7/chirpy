package utils

import (
	"encoding/json"
	"net/http"
	"strings"
)

func RespondWithError(rw http.ResponseWriter, code int, msg string) {
	type errorResponse struct {
		Error string `json:"error"`
	}

	rw.WriteHeader(code)
	errRes := errorResponse{Error: msg}
	response, _ := json.Marshal(errRes)
	rw.Write(response)
}

func RespondWithJSON(rw http.ResponseWriter, code int, payload interface{}) {
	rw.WriteHeader(code)
	response, _ := json.Marshal(payload)
	rw.Write(response)
}

func ReplaceBadWords(input string) string {
	profane := map[string]string{
		"kerfuffle": "****",
		"sharbert":  "****",
		"fornax":    "****",
	}

	split := strings.Split(input, " ")
	for i, word := range split {
		s := strings.ToLower(word)
		censored, ok := profane[s]
		if ok {
			split[i] = censored
		}
	}
	return strings.Join(split, " ")
}
