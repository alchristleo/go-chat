package main

// "ConnectionHub" as Hub maintains the set of active clients and broadcasts messages to the clients.
type ConnectionHub struct {
	// Registered clients.
	clients map[*Client]string

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

// create new connection
func createConnectionHub() *ConnectionHub {
	return &ConnectionHub{
		clients:    make(map[*Client]string),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// handle channels
func (hub *ConnectionHub) run() {
	for {
		select {
		// register client to hub
		case client := <-hub.register:
			hub.clients[client] = client.name
		// unregister client to hub
		case client := <-hub.unregister:
			if _, ok := hub.clients[client]; ok {
				delete(hub.clients, client)
				close(client.send)
			}
		// loop through registered clients and send message to their send channel
		case message := <-hub.broadcast:
			for client := range hub.clients {
				select {
				case client.send <- message:
				// if send buffer is full, assume client is dead or stuck and unregister
				default:
					close(client.send)
					delete(hub.clients, client)
				}
			}
		}
	}
}
