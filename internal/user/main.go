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

func NewCLIUser(db *sql.DB) (User, error) {
	user := User{Id: uuid.NewString(), Type: "cli"}
	_, err := db.Exec("INSERT INTO user (id, type) VALUES (?, ?)", user.Id, user.Type)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func CreateAndGetCLIUser(db *sql.DB) (User, error) {
	// TODO: search fo CLI user
	return NewCLIUser(db)
}
