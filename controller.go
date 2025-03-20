package main

type FlowController struct {
	addClient chan *Client
	delClient chan *Client
	clients   *MutClients
	broadcast chan []byte
}

func (fc *FlowController) initFlowController() {
	for {
		select {
		case cl := <-fc.addClient:
			fc.clients.addConn(cl.getConn())
			logger.Info("Added new connection!")
		case cl := <-fc.delClient:
			if fc.clients.contains(cl.getConn()) {
				fc.clients.deleteConn(cl.getConn())
				logger.Info("Deleted connection!")
			}
		}
	}
}
