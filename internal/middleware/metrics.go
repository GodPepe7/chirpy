package middleware

import (
	"html/template"
	"net/http"
)

type ApiConfig struct {
	FileserverHits int
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		cfg.FileserverHits++
		next.ServeHTTP(rw, req)
	})
}

var tmpl = template.Must(template.ParseFiles("metrics.html"))

func (cfg *ApiConfig) MetricsHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "text/html")
	rw.WriteHeader(200)
	err := tmpl.Execute(rw, cfg)
	if err != nil {
		rw.WriteHeader(500)
	}
}

func (cfg *ApiConfig) ResetHandler(rw http.ResponseWriter, req *http.Request) {
	cfg.FileserverHits = 0
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(200)
	rw.Write([]byte(http.StatusText(http.StatusOK)))
}
