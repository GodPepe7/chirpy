package main

import (
	"html/template"
	"net/http"
)

type apiConfig struct {
	FileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		cfg.FileserverHits++
		next.ServeHTTP(rw, req)
	})
}
