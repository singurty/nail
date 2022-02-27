package user

import (
	"context"
	
	"golang.org/x/crypto/bcrypt"

	"github.com/singurty/nail/db"
)

func Register(username, password string) error {
	passCrypt, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err!= nil {
		return err
	}
	_, err = db.DBpool.Exec(context.Background(), "INSERT INTO users(username, password) VALUES ($1, $2);", username, string(passCrypt))
	if err != nil {
		return err
	}
	return nil
}

func Login(username, password string) (int, error) {
	var id int
	userRow := db.DBpool.QueryRow(context.Background(), "SELECT id, password FROM USERS WHERE username = $1;", username)
	var passCrypt string
	err := userRow.Scan(&id, &passCrypt)
	if err != nil {
		return 0, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(passCrypt), []byte(password))
	if err != nil {
		return 0, err
	}
	return id, nil
}
