package models

import (
	"bitcointransaction/connection"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

//CheckTransactions is used for getting the id of non confirmed transaction
func CheckTransactions() {
	sqlStatement := `UPDATE transactions
					SET block = '-2' 
					WHERE  id in (
							SELECT id FROM transactions
							WHERE  block = '-1' LIMIT  5)
					RETURNING hash;`
	rows, err := connection.Db.Query(sqlStatement)
	if err != nil {
		fmt.Println(err)
	}
	var hashes []string
	for rows.Next() {
		var temp string
		_ = rows.Scan(&temp)
		hashes = append(hashes, temp)
	}
	if len(hashes) > 0 {
		confirmTransactions(hashes)
	} else {
		fmt.Println("Nothing to do")
	}
}

func confirmTransactions(data []string) {
	fmt.Println("Start mining")
	nonce := -1
	var hash string
	var inputFmt string
	for inputFmt != "000000" {
		nonce = nonce + 1
		hash = calculateHash(strings.Join(data, " "), nonce)
		inputFmt = hash[0:6]
	}
	fmt.Println(hash)
	var previousHash string
	var id string
	sqlStatement := `SELECT hash FROM block ORDER BY id DESC `
	connection.Db.QueryRow(sqlStatement).Scan(&previousHash)
	sqlStatement = `INSERT INTO block (hash, nounce, previousHash) VALUES ($1, $2, $3) returning Id`
	err := connection.Db.QueryRow(sqlStatement, hash, nonce, previousHash).Scan(&id)
	if err != nil {
		fmt.Println(err)
	}
	sqlStatement = `UPDATE transactions
					SET block = $1 
					WHERE  id in (
							SELECT id FROM transactions
							WHERE  block = '-2')`
	connection.Db.QueryRow(sqlStatement, id)
	connection.SendMsg()
	CheckTransactions()
}

func calculateHash(data string, nonce int) string {
	h := sha256.New()
	strToHash := data + strconv.Itoa(nonce)
	h.Write([]byte(strToHash))
	hashed := hex.EncodeToString(h.Sum(nil))
	return hashed
}
