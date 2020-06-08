package models

import (
	"bitcointransaction/connection"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/ripemd160"
)

//Transaction define the type of transaction
type Transaction struct {
	HashID      string `json:"hash"`
	Destination string `json:"destination"`
	Sender      string `json:"sender"`
	Amount      int    `json:"amount"`
	Sign        string `json:"sign"`
	Block       int    `json:"block"`
}

//InsertTransaction inserisco una transazione
func InsertTransaction(t Transaction) {
	h := ripemd160.New()
	sqlStatement := `
	INSERT INTO transactions (hash, sender, sign,block)
	VALUES ($1, $2, $3, '-1')
	RETURNING id`
	var id int
	h.Write([]byte(t.Sender))
	connection.Db.QueryRow(sqlStatement, t.HashID, string(hex.EncodeToString(h.Sum(nil))), t.Sign).Scan(&id)
	sqlStatement = `
	INSERT INTO outputs (parent, pkscript, amount, used)
	VALUES ($1, $2, $3, false)
	RETURNING id`
	h2 := ripemd160.New()
	h2.Write([]byte(t.Destination))
	err := connection.Db.QueryRow(sqlStatement, id, string(hex.EncodeToString(h2.Sum(nil))), t.Amount).Scan(&id)
	if err != nil {
		fmt.Println(err)
	}
}
