package routes

import (
	"waysbook/handlers"
	"waysbook/pkg/middleware"
	"waysbook/pkg/mysql"
	"waysbook/repositories"

	"github.com/gorilla/mux"
)

func BookRoutes(r *mux.Router) {
	bookRepository := repositories.RepositoryBook(mysql.DB)
	h := handlers.HanlderBook(bookRepository)

	r.HandleFunc("/books", h.FindBooks).Methods("GET")
	r.HandleFunc("/books-promo", h.FindBookPromo).Methods("GET")
	r.HandleFunc("/book/{id}", h.GetBook).Methods("GET")
	r.HandleFunc("/book", middleware.AuthAdmin(middleware.UploadFilePdf(middleware.UploadFileImage(h.CreateBook)))).Methods("POST")
	r.HandleFunc("/book-promo/{id}", middleware.AuthAdmin(h.UpdateBookPromo)).Methods("PATCH")
	r.HandleFunc("/book/{id}", middleware.AuthAdmin(h.UpdateBook)).Methods("PATCH")
	r.HandleFunc("/book/{id}", middleware.AuthAdmin(h.DeleteBook)).Methods("DELETE")
}
