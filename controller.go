package main

import "fmt"

type clientsController struct {
	addClient chan *Client
	delClient chan *Client
	clients   map[*Client]string
	broadcast chan *toSendMessage
}

func (cc *clientsController) addConn(client *Client) {
	cc.clients[client] = client.name
}

func (cc *clientsController) deleteConn(client *Client) {
	delete(cc.clients, client)
}

func (cc *clientsController) contains(client *Client) bool {
	_, ok := cc.clients[client]
	return ok
}

func (cc *clientsController) getMap() *map[*Client]string {
	return &cc.clients
}

func (cc *clientsController) initFlowController() {
	for {
		select {
		case cl := <-cc.addClient:
			cc.addConn(cl)
			logger.Info(fmt.Sprintf("Added '%s' to list of clients", cl.name))
		case cl := <-cc.delClient:
			if cc.contains(cl) {
				cc.deleteConn(cl)
				logger.Info(fmt.Sprintf("Deleted '%s' from list of clients", cl.name))
			}
		case msg := <-cc.broadcast:
			for cl := range *cc.getMap() {
				cl.toSend <- *msg
			}
		}
	}
}
