package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	"api/internal/coin"
	"api/internal/users"
)

func main() {
	r := mux.NewRouter()

	// mount feature areas under one API
	coin.RegisterRoutes(r.PathPrefix("/api/go/coin").Subrouter())
	users.RegisterRoutes(r.PathPrefix("/api/go/users").Subrouter())

	// optional middlewares
	handler := jsonContentTypeMiddleware(enableCORS(r))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Println("API listening on :" + port)
	log.Fatal(srv.ListenAndServe())
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow any Origin
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Check if the requests is for CORS preflight
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Pass down the request to the next middleware (or final handler)
		next.ServeHTTP(w, r)
	})

}

// handle JSON objects
func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Set JSON Content Type
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
