package main

import (
	"log"
	"net/http"
)

func main() {
	serveMux := http.NewServeMux()
	serveMux.Handle("/", http.FileServer(http.Dir(".")))
	server := http.Server{Addr: ":8080", Handler: serveMux}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
