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
	Status      int    `json:"status"`
}

//InsertTransaction inserisco una transazione
func InsertTransaction(t Transaction) {
	sqlStatement := `
	INSERT INTO transactions (hash, sender, sign, amount, isUsed, status)
	VALUES ($1, $2, $3, $4, false, 1)
	RETURNING id`
	var id int
	connection.Db.QueryRow(sqlStatement, t.HashID, t.Sender, t.Sign, t.Amount).Scan(&id)
	sqlStatement = `
	INSERT INTO outputs (parent, out_transaction, condition)
	VALUES ($1, $2, $3)
	RETURNING id`
	h := ripemd160.New()
	h.Write([]byte(t.Destination))
	err := connection.Db.QueryRow(sqlStatement, id, id, string(hex.EncodeToString(h.Sum(nil)))).Scan(&id)
	if err != nil {
		fmt.Println(err)
	}
}
