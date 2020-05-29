package controllers

import (
	"fmt"
	"net/http"
)

//HelloServer test controller
func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}
