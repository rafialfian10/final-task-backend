package routes

import (
	"github.com/gorilla/mux"
)

// membuat function RouteInit untuk membuat route ke masing-masing route
func RouteInit(r *mux.Router) {
	UserRoutes(r)
	AuthRoutes(r)
	BookRoutes(r)
	CartRoutes(r)
	TransactionRoutes(r)
}
