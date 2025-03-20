package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

type ConcSafeConns struct {
	mx sync.Mutex
	mp map[*websocket.Conn]bool
}

func (cfc *ConcSafeConns) AddConn(con *websocket.Conn) {
	cfc.mx.Lock()
	defer cfc.mx.Unlock()
	cfc.mp[con] = true
}

func (cfc *ConcSafeConns) DeleteConn(con *websocket.Conn) {
	cfc.mx.Lock()
	defer cfc.mx.Unlock()
	delete(cfc.mp, con)
}
