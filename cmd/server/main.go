package main

import (
	"database/sql"
	"ecommerce-go/internal/api/router"
	"ecommerce-go/internal/app"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
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
	dbURL := os.Getenv("DATABASE_URL")

	// Connect DB
	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		log.Fatal("fail to connect to database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("fail to ping database:", err)
	}

	a := &app.Application{
		Db:     db,
		Logger: log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime),
	}

	mux := router.Router(a)

	server := http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%v", port),
		Handler: mux,
	}
	fmt.Println("SERVER RUNNING ON PORT: ", port)

	err = server.ListenAndServe()

	if err != nil {
		log.Fatal("Error starting the server", err)
	}
}
