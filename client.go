package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
}

type MutClients struct {
	mx sync.Mutex
	mp map[*websocket.Conn]bool
}

func (mc *MutClients) addConn(con *websocket.Conn) {
	mc.mx.Lock()
	defer mc.mx.Unlock()
	mc.mp[con] = true
}

func (mc *MutClients) deleteConn(con *websocket.Conn) {
	mc.mx.Lock()
	defer mc.mx.Unlock()
	delete(mc.mp, con)
}

func (mc *MutClients) contains(con *websocket.Conn) bool {
	mc.mx.Lock()
	defer mc.mx.Unlock()
	_, ok := mc.mp[con]
	return ok
}

func createNewClient(conn *websocket.Conn) *Client {
	return &Client{conn: conn}
}

func (c *Client) getConn() *websocket.Conn {
	return c.conn
}

func (c *Client) readFromClient(fc *FlowController) {
	for {
		msgType, msg, err := c.conn.ReadMessage()
		if err != nil {
			logger.Error(err.Error())
			break
		}
		c.sendMsgToClients(fc, msgType, msg)
	}
}

func (c *Client) sendMsgToClients(fc *FlowController, msgType int, msg []byte) {
	for k := range fc.clients.mp {
		if k == c.conn {
			continue
		}
		if err := k.WriteMessage(msgType, msg); err != nil {
			logger.Error(err.Error())
		}
	}
}
