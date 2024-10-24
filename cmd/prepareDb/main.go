package main

import (
	"fiveLettersHelper/internal/dbUtils"
	"log"
)

func main() {
	db, err := dbUtils.GetDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = dbUtils.PrepareDB(db)
	if err != nil {
		log.Fatal("Error preparing DB:", err)
	}
}
