package main

import (
	"bitcointransaction/controllers"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	http.HandleFunc("/", controllers.HelloServer)
	http.ListenAndServe(":"+port, nil)
}
