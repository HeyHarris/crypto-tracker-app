package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type User struct {
	Id              int `json:"id"`
	Name            int `json:"name"`
	Email           int `json:"email"`
	CreateTimestamp int `json:"createTimestamp"`
}

//main function

func main() {
	// connect to the PostgreSQL database
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//create table if not exists
	_, err = db.Exec("Create TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name TEXT, email TEXT, createTimestamp TIMESTAMP)")
	if err != nil {
		log.Fatal(err)
	}

	// create router
	router := mux.NewRouter()
	router.HandleFunc("/api/go/users/{id}", getUser(db)).Methods("GET")
	router.HandleFunc("/api/go/users", createUser(db)).Methods("POST")
	// router.HandleFunc("/api/go/users{id}", updateUsers(db).Methods("PUT"))

	// wrap the router with CORS and JSON content type middlewares
	enhancedRouter := enableCORS(jsonContentTypeMiddleware(router))

	//start server
	log.Fatal(http.ListenAndServe(":8000", enhancedRouter))

}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow any Origin
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Check if the requets is for CORS preflight
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

// get user
func getUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pathVariables := mux.Vars(r)
		id := pathVariables["id"]

		var userTable User
		err := db.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&userTable.Id, &userTable.Name, &userTable.Email, &userTable.CreateTimestamp)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(userTable)
	}
}

// create user
func createUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userTable User
		json.NewDecoder(r.Body).Decode(userTable)

		err := db.QueryRow("INSERT INTO users (name, email, createTimestamp) VALUES ($1, $2, $3) RETURNING id", userTable.Name, userTable.Email, time.Now()).Scan(&userTable.Id)
		if err != nil {
			log.Fatal(err)
			return
		}
		json.NewEncoder(w).Encode(userTable)
	}
}
