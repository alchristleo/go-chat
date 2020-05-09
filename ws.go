package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// needed to allow connections from any origin for :3000 -> :8081
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Client struct {
	name string
	hub  *ConnectionHub
	conn *websocket.Conn
	send chan []byte
}

type JSONData struct {
	Name      string `json:"name"`
	Text      string `json:"text"`
	Timestamp string `json:"timestamp"`
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (c *Client) readPump() {
	// schedule client to be disconnected
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	// init client connection
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		fmt.Println("read from client")
		// read JSON data from connection
		message := JSONData{}
		if err := c.conn.ReadJSON(&message); err != nil {
			fmt.Println("Error reading JSON", err)
		}
		fmt.Printf("Get response: %#v\n", message)

		messageJSON, _ := json.Marshal(message)
		// queue message for writing
		c.hub.broadcast <- messageJSON
	}
}

func wsHandler(hub *ConnectionHub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		http.Error(w, "Fail to open websocket connection", http.StatusBadRequest)
	}

	// register to connection hub
	name := r.URL.Query().Get("name")
	client := &Client{
		name: name,
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}
	client.hub.register <- client

	// construct JSON list of connected client names and send to new client for display
	names := make([]string, len(client.hub.clients)+1)
	i := 0
	for k := range client.hub.clients {
		names[i] = client.hub.clients[k]
		i++
	}

	names[i] = name
	namesJSON, _ := json.Marshal(names)

	client.hub.broadcast <- namesJSON

	// read and write concurrent func for websocket
	go client.writePump()
	go client.readPump()
}
