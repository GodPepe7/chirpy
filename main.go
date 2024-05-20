package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/godpepe7/chirpy/internal/db"
	"github.com/godpepe7/chirpy/internal/handler"
	"github.com/joho/godotenv"
)

func main() {
	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if *dbg {
		db.RemoveDBFile()
	}

	// by default, godotenv will look for a file named .env in the current directory
	godotenv.Load()
	jwtSecret := os.Getenv("JWT_SECRET")

	const filepathRoot = "./html"
	const port = "8080"

	serveMux := http.NewServeMux()
	fsHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	apiConfig := handler.ApiConfig{FileserverHits: 0, JwtSecret: jwtSecret}
	database, err := db.NewDB("database")
	if err != nil {
		log.Fatal(err)
	}
	handler := handler.NewHandler(database)

	serveMux.Handle("/app/*", apiConfig.MiddlewareMetricsInc(fsHandler))
	serveMux.HandleFunc("GET /api/healthz", handler.HealthzHandler)
	serveMux.HandleFunc("GET /api/reset", apiConfig.ResetHandler)
	serveMux.HandleFunc("GET /api/chirps", handler.GetChirpHandler)
	serveMux.HandleFunc("GET /api/chirps/{id}", handler.GetChirpByIdHandler)
	serveMux.HandleFunc("POST /api/chirps", handler.PostChirpHandler)
	serveMux.HandleFunc("POST /api/users", handler.PostUserHandler)
	serveMux.HandleFunc("POST /api/login", handler.PostLoginHandler)
	serveMux.HandleFunc("GET /admin/metrics", apiConfig.MetricsHandler)

	server := &http.Server{Addr: ":" + port, Handler: serveMux}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
