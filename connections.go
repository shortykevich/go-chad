package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

type MutClients struct {
	mx sync.Mutex
	mp map[*websocket.Conn]bool
}

type FlowController struct {
	addClient chan *Client
	delClient chan *Client
	clients   *MutClients
}

func (mc *MutClients) AddConn(con *websocket.Conn) {
	mc.mx.Lock()
	defer mc.mx.Unlock()
	mc.mp[con] = true
}

func (mc *MutClients) DeleteConn(con *websocket.Conn) {
	mc.mx.Lock()
	defer mc.mx.Unlock()
	delete(mc.mp, con)
}

func (mc *MutClients) Contains(con *websocket.Conn) bool {
	mc.mx.Lock()
	defer mc.mx.Unlock()
	_, ok := mc.mp[con]
	return ok
}

func (fc *FlowController) connControlFlow() {
	for {
		select {
		case cl := <-fc.addClient:
			fc.clients.AddConn(cl.getConn())
			logger.Info("Added new connection!")
		case cl := <-fc.delClient:
			if fc.clients.Contains(cl.getConn()) {
				fc.clients.DeleteConn(cl.getConn())
				logger.Info("Deleted connection!")
			}
		}
	}
}
