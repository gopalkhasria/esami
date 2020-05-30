package models

import (
	"bitcointransaction/connection"

	"golang.org/x/crypto/bcrypt"
)

//User user model for db
type User struct {
	Name     string
	Email    string
	Password string
}

//InsertUser insericsi
func InsertUser(data User) int {
	sqlStatement := `
	INSERT INTO users (name, email, pass)
	VALUES ($1, $2, $3)
	RETURNING id`
	id := 0
	hasPass, err := bcrypt.GenerateFromPassword([]byte(data.Password), 15)

	err = connection.Db.QueryRow(sqlStatement, data.Name, data.Email, string(hasPass)).Scan(&id)
	if err != nil {
		panic(err)
	}
	return id
}
