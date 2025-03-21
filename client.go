package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// In cases where a lot of clients send messages
	toSendBuffer = 100

	messageSize = 256
	pingPeriod  = 90 * time.Second
	pongWait    = 100 * time.Second
	writeWait   = 20 * time.Second
)

type Client struct {
	fc     *FlowController
	conn   *websocket.Conn
	name   string
	toSend chan toSendMessage
}

type MutClients struct {
	mu sync.Mutex
	mp map[*Client]string
}

type toSendMessage struct {
	data    []byte
	msgType int
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
		fc:     fc,
		conn:   conn,
		name:   name,
		toSend: make(chan toSendMessage, toSendBuffer),
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
	defer func() {
		c.fc.delClient <- c
		c.closeConn()
	}()

	for {
		msgType, msg, err := c.getConn().ReadMessage()
		if err != nil {
			processDiscError(c, err)
			break
		}
		logger.Info(fmt.Sprintf("client: '%s' wrote message: '%s'", c.name, string(msg)))
		c.fc.broadcast <- &toSendMessage{
			data:    msg,
			msgType: msgType,
		}
	}
}

func (c *Client) sendMsgToClient() {

	pinger := time.NewTicker(pingPeriod)
	defer func() {
		pinger.Stop()
		c.getConn().Close()
	}()

	for {
		select {
		case msg, ok := <-c.toSend:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.getConn().WriteMessage(msg.msgType, msg.data)
			for {
				select {
				case msg := <-c.toSend:
					c.getConn().WriteMessage(msg.msgType, msg.data)
				default:
					break
				}
			}
		case <-pinger.C:
			c.getConn().SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.getConn().WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
