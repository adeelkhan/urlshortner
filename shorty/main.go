package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/adeelkhan/shorty/shortener"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/s", shortener.GetShortUrl)
	r.HandleFunc("/l", shortener.GetLongUrl)

	serverAddress := "localhost:8090"

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)
	srv := &http.Server{
		Handler:      handler,
		Addr:         serverAddress,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Printf("Server starting at...%s\n", serverAddress)
	log.Fatal(srv.ListenAndServe(), handler)
}
