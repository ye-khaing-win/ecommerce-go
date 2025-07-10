package main

import (
	"ecommerce-go/internal/api/router"
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := 3000

	mux := router.Router()

	server := http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%v", port),
		Handler: mux,
	}
	fmt.Println("SERVER RUNNING ON PORT: ", port)

	err := server.ListenAndServe()

	if err != nil {
		log.Fatal("Error starting the server", err)
	}
}
