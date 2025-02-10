package routes

import (
	"backend/controllers"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/signup", controllers.Signup).Methods("POST")
	router.HandleFunc("/verify-login", controllers.Login).Methods("GET")
	router.HandleFunc("/login", controllers.Login).Methods("POST")

	return router
}
