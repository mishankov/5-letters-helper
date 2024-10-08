package user

import (
	"database/sql"

	"github.com/google/uuid"
)

type User struct {
	Id         string
	Type       string
	Identifier string
}

func NewCLIUser(db *sql.DB) User {
	// TODO: insert does not work
	user := User{Id: uuid.NewString(), Type: "cli"}
	db.Exec("INSERT INTO user (id, type) VALUES (?, ? ,?)", user.Id, user.Type)
	return user
}

func CreateAndGetCLIUser(db *sql.DB) User {
	// TODO: search fo CLI user
	return NewCLIUser(db)
}
