package routes

import (
	"waysbook/handlers"
	"waysbook/pkg/middleware"
	"waysbook/pkg/mysql"
	"waysbook/repositories"

	"github.com/gorilla/mux"
)

func UserRoutes(r *mux.Router) {
	userRepository := repositories.RepositoryUser(mysql.DB)
	h := handlers.HandlerUser(userRepository)

	r.HandleFunc("/users", h.FindUsers).Methods("GET")
	r.HandleFunc("/user/{id}", middleware.Auth(h.GetUser)).Methods("GET")
	r.HandleFunc("/user/{id}", middleware.Auth(middleware.UploadFilePhoto(h.UpdateUser))).Methods("PATCH")
	r.HandleFunc("/user/{id}", middleware.AuthAdmin(h.DeleteUser)).Methods("DELETE")
}
