package sqlconnect

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

func ConnectDB() (*sql.DB, error) {
	dbURL := os.Getenv("DATABASE_URL")
	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		return nil, err
	}
	fmt.Println("DB CONNECTED SUCCESSFULLY.")
	return db, nil
}
