package main

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	fc   *FlowController
	conn *websocket.Conn
	name string
}

type MutClients struct {
	mu sync.Mutex
	mp map[*Client]string
}

func (mc *MutClients) addConn(client *Client) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.mp[client] = client.name
}

func (mc *MutClients) deleteConn(client *Client) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	delete(mc.mp, client)
}

func (mc *MutClients) contains(client *Client) bool {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	_, ok := mc.mp[client]
	return ok
}

func (mc *MutClients) getMap() *map[*Client]string {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	return &mc.mp
}

func createNewClient(fc *FlowController, conn *websocket.Conn, name string) *Client {
	return &Client{
		fc:   fc,
		conn: conn,
		name: name,
	}
}

func (c *Client) closeConn() {
	c.conn.Close()
}

func (c *Client) getConn() *websocket.Conn {
	return c.conn
}

func processDiscError(c *Client, err error) {
	if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
		logger.Warn(fmt.Sprintf("Unexpected close error for client '%s': %v", c.name, err))
	} else if ce, ok := err.(*websocket.CloseError); ok {
		logger.Info(fmt.Sprintf("Client '%s' disconnected with code %d: %s", c.name, ce.Code, ce.Text))
	} else {
		logger.Warn(fmt.Sprintf("Client '%s' disconnected with error: %v", c.name, err))
	}
}

func (c *Client) readFromClient() {
	for {
		msgType, msg, err := c.getConn().ReadMessage()
		if err != nil {
			processDiscError(c, err)
			break
		}
		logger.Info(fmt.Sprintf("client: '%s' wrote message: '%s'", c.name, string(msg)))
		c.sendMsgToClients(msgType, msg)
	}
}

func (c *Client) sendMsgToClients(msgType int, msg []byte) {
	for cl := range *c.fc.clients.getMap() {
		if cl == c {
			continue
		}
		if err := cl.getConn().WriteMessage(msgType, msg); err != nil {
			logger.Error(err.Error())
		}
	}
}
