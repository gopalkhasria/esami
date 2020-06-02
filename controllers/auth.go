package controllers

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"bitcointransaction/models"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

//Register return register page
func Register(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "statics/register.html")
}

//RegisterData user for templating
type RegisterData struct {
	Email    string
	Name     string
	Password string
}

//Claims token generation
type Claims struct {
	Name     string `json:"name"`
	Password string `json:"email"`
	Email    string `json:"password"`
	ID       int    `json:"id"`
	jwt.StandardClaims
}

//JwtKey key for token
var JwtKey = os.Getenv("JWTKEY")

//RegEmail send email dor confirm
func RegEmail(w http.ResponseWriter, r *http.Request) {
	expirationTime := time.Now().Add(2 * time.Hour)
	claims := &Claims{
		Name:     r.FormValue("name"),
		Password: r.FormValue("password"),
		Email:    r.FormValue("email"),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	t, _ := template.ParseFiles("statics/emailreg.html")
	var body bytes.Buffer
	headers := "Mime-version: 1.0; \nContent-Type: text/html;"
	body.Write([]byte(fmt.Sprintf("Subject: Registration bitcoin transaction\n%s\n\n", headers)))
	t.Execute(&body, struct {
		Token string
		Name  string
	}{
		Token: tokenString,
		Name:  r.FormValue("name"),
	})
	auth := smtp.PlainAuth(
		"",
		"bitcointransaction01@gmail.com",
		"Qwertyofpc_1",
		"smtp.gmail.com",
	)
	err = smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"bitcointransaction01@gmail.com",
		[]string{r.FormValue("email")},
		body.Bytes(),
	)
	if err != nil {
		log.Fatal(err)
	}
	data := RegisterData{
		Name:     r.FormValue("name"),
		Email:    r.FormValue("email"),
		Password: tokenString,
	}
	tmpl, _ := template.ParseFiles("statics/emailSent.html")
	tmpl.Execute(w, data)
}

//ConfirmEmail ultima parte per la registrazione
func ConfirmEmail(w http.ResponseWriter, r *http.Request) {
	token, ok := r.URL.Query()["token"]
	if !ok || len(token[0]) < 1 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(token[0], claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user := models.User{
		Name:     claims.Name,
		Email:    claims.Email,
		Password: claims.Password,
	}
	_ = models.InsertUser(user)
	http.Redirect(w, r, "/", http.StatusFound)
}

//Login is used for login
func Login(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("statics/login.html")
	if len(r.FormValue("email")) == 0 {
		tmpl.Execute(w, struct {
			LogError bool
			Msg      string
		}{
			LogError: false,
			Msg:      "",
		})
	} else {
		data, id, err := models.FindUser(r.FormValue("email"))
		if err != nil {
			tmpl.Execute(w, struct {
				LogError bool
				Msg      string
			}{
				LogError: true,
				Msg:      "Account not find",
			})
		} else {
			if bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(r.FormValue("password"))) != nil {
				tmpl.Execute(w, struct {
					LogError bool
					Msg      string
				}{
					LogError: true,
					Msg:      "Password is incorrect",
				})
			} else {
				expirationTime := time.Now().Add(24 * time.Hour)
				claims := &Claims{
					Name:  data.Name,
					Email: data.Email,
					ID: id,
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: expirationTime.Unix(),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, err := token.SignedString(JwtKey)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				http.SetCookie(w, &http.Cookie{
					Name:  "session_token",
					Value: tokenString,
				})
				http.Redirect(w, r, "/", http.StatusFound)
			}
		}
	}
}
