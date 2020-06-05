package controllers

import (
	"bitcointransaction/connection"
	"bitcointransaction/models"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/ripemd160"
)

type body struct {
	Amount  float32  `json:"amount"`
	Address string   `json:"address"` //receiver address
	PubKey  string   `json:"pubkey"`
	Ouputs  []output `json:"outputs"`
}

type output struct {
	Amount   float32 `json:"amount"`
	ID       int     `json:"id"`
	PkScript string  `json:"pkscript"`
	Used     bool    `json:"used"`
}

//MakeTransaction make the transaction
func MakeTransaction(w http.ResponseWriter, r *http.Request) {
	var b body
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(r.Header.Get("Authorization"), claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		verifySign, CurrentAmount := verifyOutput(b.Ouputs, claims.ID)
		if CurrentAmount == 0 {
			fmt.Println("Error")
		}
		var h hash.Hash
		h = sha256.New()
		io.WriteString(h, fmt.Sprintf("%v%v%v", b.PubKey, b.Address, b.Amount))
		signhash := h.Sum(nil)
		var privatekey *ecdsa.PrivateKey
		connection.Db.QueryRow("SELECT private_key FROM keys WHERE user_id=$1", claims.ID).Scan(&privatekey)
		r, s, serr := ecdsa.Sign(rand.Reader, privatekey, signhash)
		if serr != nil {
			fmt.Println(err)
		}
		signature := r.Bytes()
		signature = append(signature, s.Bytes()...)
		sqlStatement := `
		INSERT INTO transactions (hash, sender, sign,status)
		VALUES ($1, $2, $3, 1)
		RETURNING id`
		var id int
		connection.Db.QueryRow(sqlStatement, string(hex.EncodeToString(signhash)), b.PubKey, hex.EncodeToString(signature)).Scan(&id)

		for i, t := range b.Ouputs {
			sqlStatement := `
			INSERT INTO inputs (transaction, keyHash, sign,output)
			VALUES ($1, $2, $3, 1)`
			connection.Db.QueryRow(sqlStatement, id, t.PkScript, verifySign[i], t.ID).Scan(&id)
			sqlStatement = `UPDATE outputs
				SET used = true
				WHERE $1;`
			connection.Db.QueryRow(sqlStatement,t.ID).Scan(&id)
		}

		sqlStatement = `
			INSERT INTO outputs (parent, pkscript, amount,used)
			VALUES ($1, $2, $3, false)`
		connection.Db.QueryRow(sqlStatement, id, b.Address, b.Amount).Scan(&id)
		CurrentAmount = CurrentAmount - b.Amount
		sqlStatement = `
			INSERT INTO outputs (parent, pkscript, amount,used)
			VALUES ($1, $2, $3, false)`
		connection.Db.QueryRow(sqlStatement, id, b.PubKey, CurrentAmount).Scan(&id)
		connection.SendTransaction(id)
	}
}

func verifyOutput(data []output, id int) ([]string, float32) {
	var amount float32
	var sign []string
	var hash string
	for _, t := range data {
		var temp output
		sqlStatement := `SELECT pkscript, amount, used,hash FROM outputs
						INNER JOIN transactions ON transactions.id = outputs.parent
						WHERE id=$1`
		connection.Db.QueryRow(sqlStatement, t.ID).Scan(&temp.PkScript, &temp.Amount, &temp.Used, &hash)
		if temp.Used {
			return nil, 0
		}
		h := ripemd160.New()
		h.Write([]byte(models.GetKeys(id)))
		if t.PkScript != string(hex.EncodeToString(h.Sum(nil))) {
			return nil, 0
		}
		var priv string
		connection.Db.QueryRow(`SELECT private_key FROM keys WHERE user_id=$1`, id).Scan(&priv)
		privdecode, _ := hex.DecodeString(priv)
		privateKey, _ := x509.ParseECPrivateKey(privdecode)
		r, s, err := ecdsa.Sign(rand.Reader, privateKey, []byte(hash))
		if err != nil {
			fmt.Println(err)
		}
		signature := r.Bytes()
		signature = append(signature, s.Bytes()...)
		sign = append(sign, hex.EncodeToString(signature))
		amount += temp.Amount
	}
	return sign, amount
}
