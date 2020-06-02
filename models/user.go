package models

import (
	"bitcointransaction/connection"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"log"

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
	/*privateKey, _ := rsa.GenerateKey(rand.Reader, 256)
	publicKey := privateKey.PublicKey*/
	pubkeyCurve := elliptic.P256()
	privatekey := new(ecdsa.PrivateKey)
	privatekey, err = ecdsa.GenerateKey(pubkeyCurve, rand.Reader)
	pubkey := privatekey.PublicKey
	sqlStatement = `INSERT INTO keys (private_key, public_key, user_id) VALUES($1, $2, $3)`
	connection.Db.QueryRow(sqlStatement, fmt.Sprintf("%v", privatekey), fmt.Sprintf("%v", pubkey), id)
	return id
}

//FindUser used for find and user
func FindUser(email string) (User, int, error) {
	var data User
	var id int
	err := connection.Db.QueryRow(`SELECT * FROM users WHERE email=$1`, email).Scan(&id, &data.Name, &data.Email, &data.Password)
	if err == sql.ErrNoRows {
		return data, 0, errors.New("Account doesn't exist")
	}
	if err != nil {
		log.Fatal(err)
	}
	return data, id, nil
}

//GetKeys get user key
func GetKeys(id int) string{
	var pubkey string
	connection.Db.QueryRow(`SELECT public_key FROM keys WHERE user_id=$1`, id).Scan(&pubkey)
	return pubkey
}
