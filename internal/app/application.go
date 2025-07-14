package app

import (
	"database/sql"
	"log"
)

type Application struct {
	Db     *sql.DB
	Logger *log.Logger
}
