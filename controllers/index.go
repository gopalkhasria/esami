package controllers

import (
	"bitcointransaction/models"
	"html/template"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

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
			tmpl, _ := template.ParseFiles("statics/index.html")
			tmpl.Execute(w, struct{
				Name string
				Email string
				PubKey string
			}{
				Name: claims.Name,
				Email: claims.Email,
				PubKey: models.GetKeys(claims.ID) ,
			})
		}
	}
}
