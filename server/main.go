package main

import (
	"net/http"

	socketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	logger "github.com/sirupsen/logrus"
)

type Message struct {
	Id   int    `json:"id"`
	Text string `json:"text"`
}

func main() {
	bc := BlockChain{Chain: make([]Block, 0), BlockSize: 5, CurrentDifficulty: 1, MiningReward: 100, PendingTransactions: make([]Transaction, 0)}
	bc.generateGenesisBlock()

	server := socketio.NewServer(transport.GetDefaultWebsocketTransport())
	server.On(socketio.OnConnection, func(c *socketio.Channel) {
		logger.Println("Connected")
	})
	server.On(socketio.OnDisconnection, func(c *socketio.Channel) {
		logger.Println("Disconnected")
	})
	server.On("get-balance", func(c *socketio.Channel, s string) string {
		logger.Println(s)
		return s + " OK"
	})

	serveMux := http.NewServeMux()
	serveMux.Handle("/socket.io/", server)

	logger.Println("Starting server...")
	// listen on port 8000
	http.ListenAndServe(":8000", serveMux)
}
