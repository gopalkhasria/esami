package models

import (
	"bitcointransaction/connection"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
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
	pubkeyCurve := elliptic.P256()
	privatekey := new(ecdsa.PrivateKey)
	privatekey, err = ecdsa.GenerateKey(pubkeyCurve, rand.Reader)
	pubkey := privatekey.PublicKey
	sqlStatement = `INSERT INTO keys (private_key,user_id) VALUES($1,$2)`
	strPriv, err := x509.MarshalECPrivateKey(privatekey)
	connection.Db.QueryRow(sqlStatement, hex.EncodeToString(strPriv), id)

	privatekey2 := new(ecdsa.PrivateKey)
	privatekey2, _ = ecdsa.GenerateKey(pubkeyCurve, rand.Reader)
	pubkey2 := privatekey2.PublicKey

	transaction := Transaction{Sender: fmt.Sprintf("%v", pubkey2), Destination: fmt.Sprintf("%v", pubkey), Amount: 5}
	var h hash.Hash
	h = sha256.New()
	io.WriteString(h, fmt.Sprintf("%v", transaction))
	signhash := h.Sum(nil)
	transaction.HashID = string(hex.EncodeToString(signhash))
	r, s, err := ecdsa.Sign(rand.Reader, privatekey2, signhash)
	if err != nil {
		fmt.Println(err)
	}
	signature := r.Bytes()
	signature = append(signature, s.Bytes()...)
	transaction.Sign = hex.EncodeToString(signature)
	InsertTransaction(transaction)

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
func GetKeys(id int) string {
	var priv string
	connection.Db.QueryRow(`SELECT private_key FROM keys WHERE user_id=$1`, id).Scan(&priv)
	privdecode, _ :=  hex.DecodeString(priv)
	privateKey, _ := x509.ParseECPrivateKey(privdecode)
	response := fmt.Sprintf("%v", privateKey.PublicKey)
	return response
}
