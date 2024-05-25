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
	polkaKey := os.Getenv("POLKA_KEY")
	const filepathRoot = "./html"
	const port = "8080"

	database, err := db.NewDB("database")
	if err != nil {
		log.Fatal(err)
	}

	serveMux := http.NewServeMux()
	fsHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	apiConfig := handler.ApiConfig{FileserverHits: 0, JwtSecret: jwtSecret, DB: database, PolkaKey: polkaKey}

	serveMux.Handle("/app/*", apiConfig.MiddlewareMetricsInc(fsHandler))
	serveMux.HandleFunc("GET /api/healthz", handler.HealthzHandler)
	serveMux.HandleFunc("GET /api/reset", apiConfig.ResetHandler)
	serveMux.HandleFunc("GET /api/chirps", apiConfig.GetChirpHandler)
	serveMux.HandleFunc("GET /api/chirps/{id}", apiConfig.GetChirpByIdHandler)
	serveMux.HandleFunc("POST /api/chirps", apiConfig.PostChirpHandler)
	serveMux.HandleFunc("DELETE /api/chirps/{id}", apiConfig.DeleteChirpHandler)
	serveMux.HandleFunc("POST /api/users", apiConfig.PostUserHandler)
	serveMux.HandleFunc("PUT /api/users", apiConfig.PutUserHandler)
	serveMux.HandleFunc("POST /api/login", apiConfig.PostLoginHandler)
	serveMux.HandleFunc("POST /api/revoke", apiConfig.PostRevokeHandler)
	serveMux.HandleFunc("POST /api/refresh", apiConfig.PostRefreshHandler)
	serveMux.HandleFunc("POST /api/polka/webhooks", apiConfig.PostPolkaWebhookHandler)
	serveMux.HandleFunc("GET /admin/metrics", apiConfig.MetricsHandler)

	server := &http.Server{Addr: ":" + port, Handler: serveMux}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
