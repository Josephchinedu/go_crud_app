package main

import (
	"fmt"
	"go-postgres/router"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	r := router.Router()
	// fs := http.FileServer(http.Dir("build"))
	// http.Handle("/", fs)
	fmt.Println("Starting server on the port 8080...")

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	port := ":" + os.Getenv("PORT")

	log.Fatal(http.ListenAndServe(port, r))
}
