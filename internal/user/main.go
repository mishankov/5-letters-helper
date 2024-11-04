package user

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id         string
	Type       string
	Identifier sql.NullString
}

func NewCLIUser(db *sql.DB) (User, error) {
	user := User{Id: uuid.NewString(), Type: "cli"}
	_, err := db.Exec("INSERT INTO user (id, type, created) VALUES (?, ?, ?)", user.Id, user.Type, time.Now())
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func NewTelegramUser(telegramId int, db *sql.DB) (User, error) {
	user := User{Id: uuid.NewString(), Type: "telegram", Identifier: sql.NullString{String: strconv.Itoa(telegramId), Valid: true}}
	_, err := db.Exec("INSERT INTO user (id, type, identifier, created) VALUES (?, ?, ?, ?)", user.Id, user.Type, user.Identifier, time.Now())
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func CreateAndGetCLIUser(db *sql.DB) (User, error) {
	row := db.QueryRow("SELECT id, type, identifier FROM user WHERE type = ?", "cli")

	user := User{}

	err := row.Scan(&user.Id, &user.Type, &user.Identifier)
	switch {
	case err == sql.ErrNoRows:
		return NewCLIUser(db)
	case err != nil:
		return User{}, err
	default:
		return user, nil
	}
}

func CreateAndGetTelegramUser(telegramId int, db *sql.DB) (User, error) {
	row := db.QueryRow("SELECT id, type, identifier FROM user WHERE type = ? AND identifier = ?", "telegram", telegramId)

	user := User{}

	err := row.Scan(&user.Id, &user.Type, &user.Identifier)
	switch {
	case err == sql.ErrNoRows:
		return NewTelegramUser(telegramId, db)
	case err != nil:
		return User{}, err
	default:
		return user, nil
	}
}
