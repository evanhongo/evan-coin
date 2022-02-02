package socket

import (
	socketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	logger "github.com/sirupsen/logrus"
)

func InitSocketClient() (c *socketio.Client) {
	var err error
	c, err = socketio.Dial(
		socketio.GetUrl("localhost", 8000, false),
		transport.GetDefaultWebsocketTransport())
	if err != nil {
		logger.Fatal(err)
	}
	//defer c.Close()

	err = c.On(socketio.OnConnection, func(h *socketio.Channel) {})
	if err != nil {
		logger.Fatal(err)
	}

	err = c.On(socketio.OnDisconnection, func(h *socketio.Channel) {})
	if err != nil {
		logger.Fatal(err)
	}

	return
}
