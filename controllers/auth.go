package controllers

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"time"
	"net/http"
	"net/smtp"

	"github.com/dgrijalva/jwt-go"
)

//HelloServer test controller
func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

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
	jwt.StandardClaims
}

//RegEmail send email dor confirm
func RegEmail(w http.ResponseWriter, r *http.Request) {
	var jwtKey = []byte("my_secret_key")
	expirationTime := time.Now().Add(60 * time.Minute)
	claims := &Claims{
		Name:     r.FormValue("name"),
		Password: r.FormValue("password"),
		Email:    r.FormValue("email"),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
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
		Name string 
	}{
		Token: tokenString,
		Name: r.FormValue("name"),
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
