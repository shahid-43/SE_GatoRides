package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"backend/config"
	"backend/routes"
)

func main() {
	config.ConnectDB()
	router := routes.SetupRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server running on port", port)
	fmt.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":"+port, router))

}
