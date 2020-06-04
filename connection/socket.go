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
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		//fmt.Println(string(p))
		t := getTransactions()
		m := msg{Azione: 1, Transaction: t}
		bytes, err := json.Marshal(m)
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

type msg struct {
	Azione      int           `json:"azione"`
	Transaction []transaction `json:"transaction"`
}

type transaction struct {
	ID     int    `json:"id"`
	HashID string `json:"hash"`
	Sender string `json:"sender"`
	Sign   string `json:"sign"`
	Status int    `json:"status"`
	Output output `json:"output"`
}

type output struct {
	ID       int    `json:"id"`
	PkScript string `json:"pkScript"`
	Amount   string `json:"amount"`
	Used     bool   `json:"used"`
}

//GetTransactions fetch all transactions
func getTransactions() []transaction {
	t := []transaction{}
	/*var id int
	p := false*/
	rows, err := Db.Query(`SELECT transactions.id,hash, sender,sign,status, outputs.id,
		pkscript, amount,used 
		FROM transactions 
		INNER JOIN outputs ON transactions.id = outputs.parent`)
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		temp := transaction{}
		_ = rows.Scan(&temp.ID, &temp.HashID, &temp.Sender, &temp.Sign, &temp.Status, &temp.Output.ID, &temp.Output.PkScript, &temp.Output.Amount, &temp.Output.Used)
		t = append(t, temp)
	}
	return t
}
