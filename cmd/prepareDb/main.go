package main

import (
	"fiveLettersHelper/internal/db"
	"log"
)

func main() {
	db, err := db.GetDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

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

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
	}

}
