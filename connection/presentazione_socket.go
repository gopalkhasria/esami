package connection

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var presentazioneClients = make(map[*websocket.Conn]bool)

//PresentazioneSocketStart starto la presentazione socket
func PresentazioneSocketStart(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	presentazioneClients[ws] = true
	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		for client := range presentazioneClients {
			err := client.WriteMessage(websocket.TextMessage, p)
			if err != nil {
				log.Printf("Websocket error: %s", err)
				client.Close()
				delete(Clients, client)
			}
		}
	}
}
