package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

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
		json.NewEncoder(w).Encode("Congrats! Your API is now setup!")
	})

	// pathPrefix untuk membuat route baru. Subrouter untuk menguji route pada pathPrefix.
	routes.RouteInit(r.PathPrefix("/api/v1").Subrouter())

	// route untuk menginisialisasi folder dengan file, image css, js agar dapat diakses kedalam project
	r.PathPrefix("/uploads").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	// Config CORS (fitur yang berfungsi untuk memberikan akses dari frontend)
	var allowedHeaders = handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	var allowedMethods = handlers.AllowedMethods([]string{"GET", "POST", "PATCH", "DELETE", "PUT", "HEAD"})
	var allowedOrigins = handlers.AllowedOrigins([]string{"*"})

	var PORT = os.Getenv("PORT")

	fmt.Println("Your server at http://localhost:5000")
	http.ListenAndServe(":"+PORT, handlers.CORS(allowedHeaders, allowedMethods, allowedOrigins)(r))
}

// PORT=5000
// SERVER_KEY=SB-Mid-server-CBYg0a0CWSxQrUrIYbcaHJvM
// CLIENT_KEY=SB-Mid-client-xBHWdiuU4aVE9vOq
// SYSTEM_EMAIL=rafialfian770@gmail.com
// SYSTEM_PASSWORD=ffcjtorkjkynndrk
// CLOUD_NAME=dixxnrj9b
// API_KEY=899448373737227
// API_SECRET=uls7g_goNxDCJEqEYyLmkG1XC-g
