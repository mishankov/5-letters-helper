package db

import (
	"database/sql"

	_ "github.com/glebarez/go-sqlite"
)

func GetDB() (*sql.DB, error) {
	return sql.Open("sqlite", "./fiveLettersHelp.db")
}
