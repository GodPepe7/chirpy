package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type parameters struct {
	Body string `json:"body"`
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	serveMux := http.NewServeMux()
	fsHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	apiConfig := apiConfig{FileserverHits: 0}

	serveMux.Handle("/app/*", apiConfig.middlewareMetricsInc(fsHandler))
	serveMux.HandleFunc("GET /api/healthz", healthzHandler)
	serveMux.HandleFunc("GET /api/reset", apiConfig.resetHandler)
	serveMux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)
	serveMux.HandleFunc("GET /admin/metrics", apiConfig.metricsHandler)

	server := &http.Server{Addr: ":" + port, Handler: serveMux}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func healthzHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(200)
	rw.Write([]byte(http.StatusText(http.StatusOK)))
}

func validateChirpHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(rw, 500, "Something went wrong")
		return
	}
	if len(params.Body) > 140 {
		respondWithError(rw, 400, "Chirp is too long")
		return
	}

	type cleanedResponse struct {
		CleanedBody string `json:"cleaned_body"`
	}

	cleanedString := replaceBadWords(params.Body)
	respondWithJSON(rw, 200, cleanedResponse{CleanedBody: cleanedString})
}

func respondWithError(rw http.ResponseWriter, code int, msg string) {
	type errorResponse struct {
		Error string `json:"error"`
	}

	rw.WriteHeader(code)
	errRes := errorResponse{Error: msg}
	response, _ := json.Marshal(errRes)
	rw.Write(response)
}

func respondWithJSON(rw http.ResponseWriter, code int, payload interface{}) {
	rw.WriteHeader(code)
	response, _ := json.Marshal(payload)
	rw.Write(response)
}

func replaceBadWords(input string) string {
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
