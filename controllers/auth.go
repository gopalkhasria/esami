package controllers

import (
	"fmt"
	"net/http"
)

//HelloServer test controller
func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

//Register return register page
func Register(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "statics/register.html")
}
