package controllers

import (
	"bitcointransaction/connection"
	"bitcointransaction/models"
	"encoding/hex"
	"fmt"
	"net/http"
	"text/template"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/ripemd160"
)

type data struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	PubKey string `json:"Pubkey"`
	Token  string `json:"token"`
}

//Presentazione rispondo mandando la pagine della presentazione
func Presentazione(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "statics/presentazione.html")
}

//PresentazioneAction rispondo mandando la pagine per muovere la presentazione
func PresentazioneAction(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "statics/presentazioneAction.html")
}

//PresentazioneActionTouch rispondo mandando la pagine per muovere la presentazione
func PresentazioneActionTouch(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "statics/presentazioneActionTouch.html")
}

//Index test controller
func Index(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tkn := c.Value
	if len(tkn) == 0 {
		http.Redirect(w, r, "/login", http.StatusFound)
	} else {
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tkn, claims, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			tmp, _ := template.ParseFiles("statics/index.html")
			h := ripemd160.New()
			h.Write([]byte(models.GetKeys(claims.ID)))
			data := data{Email: claims.Email, Name: claims.Name, PubKey: string(hex.EncodeToString(h.Sum(nil)))}
			data.Token, _ = tkn.SignedString(JwtKey)
			tmp.Execute(w, data)
		}
	}
}

//TransactionResult the result of query for get transaction
type TransactionResult struct {
	Hash     string
	Sender   string
	Sign     string
	Block    string
	PkScript string
	Output   []OutputTransaction
	Input    InputTransaction
}

//OutputTransaction output of the transaction
type OutputTransaction struct {
	Amount   string
	PkScript string
}

//InputTransaction input of the transaction
type InputTransaction struct {
	Keyhash string
	Sign    string
	OutID   string
	Amount  string
}

//GetTransaction show the transaction
func GetTransaction(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tkn := c.Value
	if len(tkn) == 0 {
		http.Redirect(w, r, "/login", http.StatusFound)
	} else {
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tkn, claims, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			id, ok := r.URL.Query()["id"]
			if !ok || len(id[0]) < 1 {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			sqlStatement := `SELECT transactions.hash, transactions.sender, transactions.sign, transactions.block, outputs.amount, outputs.pkscript,
					inputs.sign, inputs.keyhash,inputs.output
					FROM transactions
					INNER JOIN outputs ON transactions.id = outputs.parent
					INNER JOIN inputs ON transactions.id = inputs.transaction
					WHERE transactions.id = $1`
			rows, err := connection.Db.Query(sqlStatement, id[0])
			if err != nil {
				http.ServeFile(w, r, "statics/CoinbaseError.html")
			}
			var response TransactionResult
			var output []OutputTransaction
			i := 0
			for rows.Next() {
				var tempOut OutputTransaction
				var tempIn InputTransaction
				err := rows.Scan(&response.Hash, &response.Sender, &response.Sign, &response.Block, &tempOut.Amount, &tempOut.PkScript, &response.Input.Sign, &response.Input.Keyhash, &response.Input.OutID)
				if err != nil {
					fmt.Println(err)
				}
				tempIn.Keyhash = hex.EncodeToString([]byte(tempIn.Keyhash))
				if i == 0 {
					output = append(output, tempOut)
					i++
				} else {
					if output[i-1].PkScript != tempOut.PkScript {
						output = append(output, tempOut)
						i++
					}
				}
			}
			sqlStatement = `SELECT outputs.pkscript, outputs.amount
					FROM inputs
					INNER JOIN outputs on inputs.output = outputs.id
					WHERE inputs.output = $1`
			err = connection.Db.QueryRow(sqlStatement, response.Input.OutID).Scan(&response.PkScript, &response.Input.Amount)
			if err != nil {
				fmt.Println(err)
			}
			response.Output = output
			if response.Hash == "" {
				http.ServeFile(w, r, "statics/CoinbaseError.html")
			}
			tmp, _ := template.ParseFiles("statics/showTransactions.html")
			tmp.Execute(w, response)
		}
	}
}

//Block block data
type Block struct {
	ID           string
	Hash         string
	Nounce       string
	PreviousHash string
	TimeStamp    string
	Transactions []Transaction
}

//Transaction transactions of block
type Transaction struct {
	ID   string
	Hash string
}

//GetBlock get a specific bloxk data
func GetBlock(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tkn := c.Value
	if len(tkn) == 0 {
		http.Redirect(w, r, "/login", http.StatusFound)
	} else {
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tkn, claims, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			id, ok := r.URL.Query()["id"]
			if !ok || len(id[0]) < 1 {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			var Response Block
			sqlStatement := `SELECT * FROM block WHERE id = $1`
			err := connection.Db.QueryRow(sqlStatement, id[0]).Scan(&Response.ID, &Response.Hash, &Response.Nounce, &Response.PreviousHash, &Response.TimeStamp)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
			}
			sqlStatement = `SELECT hash, id FROM transactions WHERE block = $1`
			rows, err := connection.Db.Query(sqlStatement, id[0])
			var transactions []Transaction
			for rows.Next() {
				var t Transaction
				err = rows.Scan(&t.Hash, &t.ID)
				transactions = append(transactions, t)
			}
			Response.Transactions = transactions
			tmp, _ := template.ParseFiles("statics/showBlock.html")
			tmp.Execute(w, Response)
		}
	}
}
