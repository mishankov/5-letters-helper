package dbUtils

import (
	"database/sql"
	"os"

	_ "github.com/glebarez/go-sqlite"
)

const DB_PATH = "./db/fiveLettersHelper.db"

func GetDBFile() ([]byte, error) {
	data, err := os.ReadFile(DB_PATH)
	if err != nil {
		return []byte{}, err
	}

	return data, nil
}

func GetDB() (*sql.DB, error) {
	return sql.Open("sqlite", DB_PATH)
}

func GetTestDB() (*sql.DB, error) {
	return sql.Open("sqlite", "testDB.db")
}

func PrepareDB(db *sql.DB) error {
	sqlStmt := `
CREATE TABLE IF NOT EXISTS "game" (
	"id" TEXT NOT NULL UNIQUE,
	"user" TEXT,
	"status" TEXT,
	"created" TEXT,
	"updated" TEXT,
	PRIMARY KEY("id"),
	FOREIGN KEY ("user") REFERENCES "user"("id")
	ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE TABLE IF NOT EXISTS "user" (
	"id" TEXT NOT NULL UNIQUE,
	"type" TEXT NOT NULL,
	"identifier" TEXT,
	"created" TEXT,
	PRIMARY KEY("id")
);

CREATE TABLE IF NOT EXISTS "guess" (
	"id" TEXT NOT NULL UNIQUE,
	"game" TEXT NOT NULL,
	"number" NUMERIC NOT NULL,
	"word" TEXT,
	"result" TEXT,
	"created" TEXT,
	"updated" TEXT,
	PRIMARY KEY("id"),
	FOREIGN KEY ("game") REFERENCES "game"("id")
	ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE TABLE IF NOT EXISTS "user_settings" (
	"user" TEXT NOT NULL UNIQUE,
	"debug_mode_enabled" BOOLEAN,
	"created" TEXT,
	"updated" TEXT,
	PRIMARY KEY("user"),
	FOREIGN KEY ("user") REFERENCES "user"("id")
	ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE TABLE IF NOT EXISTS "word" (
	"word" TEXT NOT NULL UNIQUE,
	"overall_value" NUMERIC,
	PRIMARY KEY("word")
);`

	_, err := db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	return nil
}
