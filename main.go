package main

import (
	connections "bitcointransaction/connection"
	"bitcointransaction/controllers"
	"bitcointransaction/models"
	"os"

	"fmt"
	"log"
	"net/http"
)

func main() {
	connections.Connect()
	fmt.Println("Listening")

	http.HandleFunc("/", controllers.Index)
	http.HandleFunc("/register", controllers.Register)
	http.HandleFunc("/regemail", controllers.RegEmail)
	http.HandleFunc("/confirmEmail", controllers.ConfirmEmail)
	http.HandleFunc("/login", controllers.Login)
	http.HandleFunc("/ws", connections.SocketStart)
	http.HandleFunc("/makeTransaction", controllers.MakeTransaction)
	http.HandleFunc("/transaction", controllers.GetTransaction)
	http.HandleFunc("/block", controllers.GetBlock)
	http.HandleFunc("/presentazione", controllers.Presentazione)
	http.HandleFunc("/presSocket", connections.PresentazioneSocketStart)
	http.HandleFunc("/actions", controllers.PresentazioneAction)
	http.HandleFunc("/actionstouch", controllers.PresentazioneActionTouch)
	go models.CheckTransactions()
	fileServer := http.FileServer(http.Dir("./statics/"))
	http.Handle("/statics/", http.StripPrefix("/statics", fileServer))

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
	//log.Fatal(http.ListenAndServe(":5000", nil))
}
