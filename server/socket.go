// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	socketWriteWait = 10 * time.Second
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *SocketHub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()
	for {
		msg := <-c.send
		if err := c.conn.SetWriteDeadline(time.Now().Add(socketWriteWait)); err != nil {
			return
		}
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
}

// ServeWebsocket handles websocket requests from the peer.
func ServeWebsocket(hub *SocketHub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// run writing to client in goroutine
	go client.writePump()
}

// SocketHub maintains the set of active clients and broadcasts messages to the
// clients.
type SocketHub struct {
	// Registered clients.
	clients map[*Client]bool

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// status stream
	statusStream chan *Status
}

func NewSocketHub(status *Status) *SocketHub {
	// Attach a goroutine that will push meter status information
	// periodically
	var statusstream = make(chan *Status)
	go func() {
		for {
			time.Sleep(SECONDS_BETWEEN_STATUSUPDATE * time.Second)
			status.Update()
			statusstream <- status
		}
	}()

	return &SocketHub{
		register:     make(chan *Client),
		unregister:   make(chan *Client),
		clients:      make(map[*Client]bool),
		statusStream: statusstream,
	}
}

func (h *SocketHub) Broadcast(i interface{}) {
	if len(h.clients) > 0 {
		message, err := json.Marshal(i)
		if err != nil {
			log.Fatal(err)
		}

		for client := range h.clients {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(h.clients, client)
			}
		}
	}
}

func (h *SocketHub) Run(in QuerySnipChannel) {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case obj, ok := <-in:
			if !ok {
				return // break if channel closed
			}
			// make sure to pass a pointer or MarshalJSON won't work
			h.Broadcast(&obj)
		case obj := <-h.statusStream:
			h.Broadcast(obj)
		}
	}
}
