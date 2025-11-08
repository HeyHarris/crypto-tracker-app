package users

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type User struct {
	Id              int       `json:"id"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	CreateTimestamp time.Time `json:"createTimestamp"`
}

// RegisterRoutes mounts users endpoints on the passed subrouter.
// Example mount from main(): users.RegisterRoutes(r.PathPrefix("/api/go/users").Subrouter())
func RegisterRoutes(r *mux.Router) {
	dsn := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	// Ensure table exists
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT,
		createTimestamp TIMESTAMP DEFAULT NOW()
	)`); err != nil {
		panic(err)
	}

	// Optional readiness loop
	for i := 0; i < 30; i++ {
		if err = db.Ping(); err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	if err != nil {
		panic(err)
	}

	// Routes
	r.HandleFunc("", getUsers(db)).Methods("GET")     // GET /api/go/users
	r.HandleFunc("/{id}", getUser(db)).Methods("GET") // GET /api/go/users/123
	r.HandleFunc("", createUser(db)).Methods("POST")  // POST /api/go/users
}

// get all users
func getUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		rows, err := db.Query("SELECT * FROM users")
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		usersList := []User{}
		for rows.Next() {
			var user User
			if err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.CreateTimestamp); err != nil {
				log.Fatal(err)
			}

			usersList = append(usersList, user)

		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(usersList)
	}
}

// get user
func getUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pathVariables := mux.Vars(r)
		idStr := pathVariables["id"]

		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid id format — must be an integer", http.StatusBadRequest)
			return
		}
		var userTable User
		err = db.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&userTable.Id, &userTable.Name, &userTable.Email, &userTable.CreateTimestamp)
		if err != nil {
			var errorString string = "User with Id of " + strconv.Itoa(id) + " not found in our records!"
			http.Error(w, errorString, http.StatusNotFound) // better way, you can also put custom error message
			// w.WriteHeader(http.StatusNotFound) //Will Add a HTTP response code in header
			return
		}
		json.NewEncoder(w).Encode(userTable)
	}
}

// create user
func createUser(db *sql.DB) http.HandlerFunc {
	type usersRequest struct {
		Name  any `json:"name"`
		Email any `json:"email"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var body usersRequest
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&body); err != nil {
			http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Validate that Name and Email are Strings
		var name string
		switch v := body.Name.(type) {
		case string:
			name = v
		case float64: // JSON numbers decode to float64
			http.Error(w, "Name must be a string", http.StatusBadRequest)
			return
		case nil:
			http.Error(w, "Name is Required", http.StatusBadRequest)
			return
		default:
			http.Error(w, "Name must be a string", http.StatusBadRequest)
			return
		}

		// Coerce name to string (or reject if you prefer)
		var email string
		switch v := body.Email.(type) {
		case string:
			email = v
		case float64: // JSON numbers decode to float64
			http.Error(w, "email must be a string", http.StatusBadRequest)
			return
		case nil:
			http.Error(w, "Email is Required", http.StatusBadRequest)
			return
		default:
			http.Error(w, "Email must be a string", http.StatusBadRequest)
			return
		}

		if name == "" || email == "" {
			http.Error(w, "Name and Email are required", http.StatusBadRequest)
			return
		}

		var user User
		user.Name = name
		user.Email = email

		if err := db.QueryRow(
			"INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id, createTimestamp",
			user.Name, user.Email,
		).Scan(&user.Id, &user.CreateTimestamp); err != nil {
			log.Printf("insert error: %v", err)
			http.Error(w, "failed to create user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(user)
	}
}
