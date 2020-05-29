package main

import (
	"bitcointransaction/controllers"
	"net/http"
)

func main() {
	http.HandleFunc("/", controllers.HelloServer)
	http.ListenAndServe(":8080", nil)
}
