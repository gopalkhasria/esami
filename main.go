package main

import (
	connections "bitcointransaction/connection"
	"bitcointransaction/controllers"

	"fmt"
	"log"
	"net/http"
)

func main() {
	connections.Connect()
	fmt.Println("Listening on the port 5000")

	http.HandleFunc("/", controllers.HelloServer)
	http.HandleFunc("/register", controllers.Register)
	http.HandleFunc("/regemail", controllers.RegEmail)

	fileServer := http.FileServer(http.Dir("./statics/"))
	http.Handle("/statics/", http.StripPrefix("/statics", fileServer))

	log.Fatal(http.ListenAndServe(":5000", nil))
}
