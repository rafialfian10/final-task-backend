package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"waysbook/database"
	"waysbook/pkg/mysql"
	"waysbook/routes"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	// ENV config
	errEnv := godotenv.Load()
	if errEnv != nil {
		fmt.Println("Failed to load .env file")
	}

	// Database Init
	mysql.DatabaseInit()

	// Run Migration
	database.RunMigration()

	// Mux Router
	r := mux.NewRouter()

	// Routes
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Congrats! Your Dumbass API is now setup!")
	})

	routes.RouteInit(r.PathPrefix("/api/v1").Subrouter())

	// path file
	r.PathPrefix("/uploads").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	// Config CORS
	var allowedHeaders = handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	var allowedMethods = handlers.AllowedMethods([]string{"GET", "POST", "PATCH", "DELETE", "PUT", "HEAD"})
	var allowedOrigins = handlers.AllowedOrigins([]string{"*"})

	// var port = os.Getenv("PORT")
	// var port = 5000
	fmt.Println("Your server at http://localhost:5000")
	http.ListenAndServe("localhost:5000", handlers.CORS(allowedHeaders, allowedMethods, allowedOrigins)(r))
}
