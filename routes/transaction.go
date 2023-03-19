package routes

import (
	"waysbook/handlers"
	"waysbook/pkg/middleware"
	"waysbook/pkg/mysql"
	"waysbook/repositories"

	"github.com/gorilla/mux"
)

func TransactionRoutes(r *mux.Router) {
	transactionRepository := repositories.RepositoryTransaction(mysql.DB)
	h := handlers.HandlerTransaction(transactionRepository)

	r.HandleFunc("/transactions-admin", middleware.AuthAdmin(h.FindTransactions)).Methods("GET")
	r.HandleFunc("/transactions", middleware.Auth(h.FindTransactionsByUser)).Methods("GET")
	r.HandleFunc("/transaction/{id}", middleware.Auth(h.GetDetailTransaction)).Methods("GET")
	r.HandleFunc("/transaction", middleware.Auth(h.CreateTransaction)).Methods("POST")
	// r.HandleFunc("/transaction/{id}", middleware.Auth(h.UpdateTransactionStatus)).Methods("PATCH")
	r.HandleFunc("/transaction-admin/{id}", middleware.AuthAdmin(h.UpdateTransactionByAdmin)).Methods("PATCH")
	r.HandleFunc("/notification", h.Notification).Methods("POST")
}
