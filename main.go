package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

const (
	server = "server"
	client = "client"
)

var (
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	mode   = flag.String("mode", server, "'server' to run server\n'client' to run client")
	host   = flag.String("host", "127.0.0.1", "Host name")
	port   = flag.String("port", "8554", "Port")

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true }, // Probably it would be a really bad idea for production usage
	}
)

func wsHandler(fc *FlowController, w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("Error occured during upgrading HTTP to Websocket connection:", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ws.SetReadLimit(messageSize)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(appData string) error {
		ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	client := createNewClient(fc, ws, r.RemoteAddr)
	fc.addClient <- client

	go client.readFromClient()
	go client.sendMsgToClient()
}

func main() {
	flag.Parse()
	switch *mode {
	case server:
		runServer()
	case client:
		runClient()
	}
}
