package main

import "github.com/gorilla/websocket"

type Client struct {
	conn *websocket.Conn
}

func createNewClient(conn *websocket.Conn) *Client {
	return &Client{conn: conn}
}

func (c *Client) getConn() *websocket.Conn {
	return c.conn
}

func (c *Client) readFromClient() {
	for {
		msgType, msg, err := c.conn.ReadMessage()
		if err != nil {
			logger.Error(err.Error())
			break
		}
		c.sendMsgToClients(msgType, msg)
	}
}

// TODO: find a better way broadcasting to users
// cause control flow instance should be moved into main func
func (c *Client) sendMsgToClients(msgType int, msg []byte) {
	for k := range flowCont.clients.mp {
		if k == c.conn {
			continue
		}
		// go func(c *websocket.Conn) {
		if err := k.WriteMessage(msgType, msg); err != nil {
			logger.Error(err.Error())
		}
		// }(k)
	}
}
