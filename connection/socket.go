package connection

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

//Clients connected clients
var Clients = make(map[*websocket.Conn]bool)

//Ws web socket pointer
var Ws *websocket.Conn

//SocketStart start the scoket
func SocketStart(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("Client Connected")
	Clients[ws] = true
	reader(ws)
	Ws = ws
}
func reader(conn *websocket.Conn) {
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		fmt.Println(string(p))
		t := getTransactions()
		bytes, err := json.Marshal(t)
		for client := range Clients {
			err := client.WriteMessage(websocket.TextMessage, bytes)
			if err != nil {
				log.Printf("Websocket error: %s", err)
				client.Close()
				delete(Clients, client)
			}
		}

	}
}

type transaction struct {
	HashID      string `json:"hash"`
	Destination string `json:"destination"`
	Sender      string `json:"sender"`
	Amount      int    `json:"amount"`
	Sign        string `json:"sign"`
	Status      int    `json:"status"`
}

//GetTransactions fetch all transactions
func getTransactions() []transaction {
	t := []transaction{}
	var id int
	p := false
	rows, err := Db.Query("SELECT * FROM transactions")
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		temp := transaction{}
		_ = rows.Scan(&id, &temp.HashID, &temp.Sender, &temp.Sign, &temp.Amount, &p, &temp.Status)
		t = append(t, temp)
	}
	return t
}
