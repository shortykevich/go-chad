package main

import (
	"fmt"
)

type clientsController struct {
	addClient chan *Client
	delClient chan *Client
	clients   map[*Client]string
	broadcast chan *toSendMessage
}

func (cc *clientsController) linkClient(client *Client) {
	cc.clients[client] = client.name
}

func (cc *clientsController) unlinkClient(client *Client) {
	delete(cc.clients, client)
}

func (cc *clientsController) clientExists(client *Client) bool {
	_, ok := cc.clients[client]
	return ok
}

func (cc *clientsController) getClients() *map[*Client]string {
	return &cc.clients
}

func (cc *clientsController) initClientsController() {
	for {
		select {
		case cl := <-cc.addClient:
			cc.linkClient(cl)
			logger.Info(fmt.Sprintf("Added '%s' to list of clients", cl.name))
		case cl := <-cc.delClient:
			if cc.clientExists(cl) {
				cc.unlinkClient(cl)
				logger.Info(fmt.Sprintf("Deleted '%s' from list of clients", cl.name))
			}
		case msg := <-cc.broadcast:
			for cl := range *cc.getClients() {
				if msg.client == cl {
					continue
				}
				cl.toSend <- *msg
			}
		}
	}
}
