package controllers

import (
	"bitcointransaction/models"
	"encoding/hex"
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
			
		}
	}
}
