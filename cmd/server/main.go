package main

import (
	"ecommerce-go/internal/api/router"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err.Error())
	}
	port := os.Getenv("HTTP_PORT")

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
