package main

import "github.com/gorilla/websocket"

type Client struct {
	conn *websocket.Conn
}

func createNewClient(conn *websocket.Conn) *Client {
	return &Client{conn: conn}
}

func (c *Client) ReadFromConns() {
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			logger.Error(err.Error())
			break
		}
		c.SendMsgToConns(msg)
	}
}

func (c *Client) SendMsgToConns(msg []byte) {
	for k := range connections.mp {
		if k == c.conn {
			continue
		}
		go func(c *websocket.Conn) {
			if err := c.WriteMessage(websocket.TextMessage, msg); err != nil {
				logger.Error(err.Error())
			}
		}(k)
	}
}
