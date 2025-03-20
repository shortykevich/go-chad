package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

	host = flag.String("host", "127.0.0.1", "Host name")
	port = flag.String("port", "8554", "Port")
)

func main() {
	flag.Parse()
	mux := http.NewServeMux()
	flowController := &FlowController{
		addClient: make(chan *Client),
		delClient: make(chan *Client),
		clients:   &MutClients{mp: make(map[*websocket.Conn]bool)},
		broadcast: make(chan []byte),
	}

	go flowController.initFlowController()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		wsHandler(flowController, w, r)
	})

	url := fmt.Sprintf("%v:%v", *host, *port)
	logger.Info(fmt.Sprintf("Server running on: '%s'", url))
	err := http.ListenAndServe(url, mux)
	if err != nil {
		log.Fatal(err)
	}
}

func wsHandler(fc *FlowController, w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	client := createNewClient(ws)
	defer func() {
		fc.delClient <- client
		ws.Close()
	}()
	fc.addClient <- client

	client.readFromClient(fc)
}
