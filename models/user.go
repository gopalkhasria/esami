package models

import (
	"bitcointransaction/connection"
	"crypto/rand"
	"crypto/rsa"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	//"golang.org/x/crypto/bcrypt"
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
	privateKey, _ := rsa.GenerateKey(rand.Reader, 256)
	publicKey := privateKey.PublicKey
	sqlStatement = `INSERT INTO keys (private_key, public_key, user_id) VALUES($1, $2, $3)`
	connection.Db.QueryRow(sqlStatement, fmt.Sprintf("%v", privateKey), fmt.Sprintf("%v", publicKey), id)
	return id
}
