package main

import "fmt"

type FlowController struct {
	addClient chan *Client
	delClient chan *Client
	clients   *MutClients
}

func (fc *FlowController) initFlowController() {
	for {
		select {
		case cl := <-fc.addClient:
			fc.clients.addConn(cl)
			logger.Info(fmt.Sprintf("Added '%s' to list of clients", cl.name))
		case cl := <-fc.delClient:
			if fc.clients.contains(cl) {
				fc.clients.deleteConn(cl)
				logger.Info(fmt.Sprintf("Deleted '%s' from list of clients", cl.name))
			}
		}
	}
}
