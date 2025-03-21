package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

func runServer() {
	flag.Parse()
	mux := http.NewServeMux()
	flowController := &FlowController{
		addClient: make(chan *Client),
		delClient: make(chan *Client),
		clients:   &MutClients{mp: make(map[*Client]string)},
		broadcast: make(chan *toSendMessage),
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

func runClient() {
	url := fmt.Sprintf("ws://%s:%s/", *host, *port)

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		logger.Error("Failed to connect to server:", "error", err)
		return
	}
	defer conn.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					logger.Error("Connection error:", "error", err)
				}
				return
			}
			fmt.Printf("%s\n", message)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Connected to chat server. Type messages and press enter to send:")

	for {
		select {
		case <-done:
			return
		default:
			if scanner.Scan() {
				text := scanner.Text()
				err := conn.WriteMessage(websocket.TextMessage, []byte(text))
				if err != nil {
					logger.Error("Failed to send message:", "error", err)
					return
				}
			} else {
				if err := scanner.Err(); err != nil {
					logger.Error("Input error:", "error", err)
				}

				conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(
					websocket.CloseNormalClosure, ""))
				return
			}
		}
	}
}
