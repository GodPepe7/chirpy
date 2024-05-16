package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/godpepe7/chirpy/internal/db"
	"github.com/godpepe7/chirpy/internal/middleware"
)

type parameters struct {
	Body string `json:"body"`
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	serveMux := http.NewServeMux()
	fsHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	apiConfig := middleware.ApiConfig{FileserverHits: 0}

	serveMux.Handle("/app/*", apiConfig.MiddlewareMetricsInc(fsHandler))
	serveMux.HandleFunc("GET /api/healthz", healthzHandler)
	serveMux.HandleFunc("GET /api/reset", apiConfig.ResetHandler)
	serveMux.HandleFunc("GET /api/chirps", getChirpHandler)
	serveMux.HandleFunc("POST /api/chirps", postChirpHandler)
	serveMux.HandleFunc("GET /admin/metrics", apiConfig.MetricsHandler)

	server := &http.Server{Addr: ":" + port, Handler: serveMux}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func healthzHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(200)
	rw.Write([]byte(http.StatusText(http.StatusOK)))
}

func getChirpHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

}

func postChirpHandler(rw http.ResponseWriter, req *http.Request) {
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

	cleanedString := replaceBadWords(params.Body)
	db, err := db.NewDB("database")
	if err != nil {
		fmt.Println(err)
		respondWithError(rw, 500, "Something went wrong")
	}
	chirp, err := db.CreateChirp(cleanedString)
	if err != nil {
		fmt.Println(err)
		respondWithError(rw, 500, "Something went wrong")
	}
	respondWithJSON(rw, 200, chirp)
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
