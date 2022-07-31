package main

import (
	"fmt"
	"go-postgres/router"
	"log"
	"net/http"
	"os"
)

func main() {

	r := router.Router()
	// fs := http.FileServer(http.Dir("build"))
	// http.Handle("/", fs)
	fmt.Println("Starting server on the port 8080...")

	port := ":" + os.Getenv("PORT")

	log.Fatal(http.ListenAndServe(port, r))
}
