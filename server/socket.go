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

// SocketClient is a middleman between the websocket connection and the hub.
type SocketClient struct {
	hub *SocketHub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// writePump pumps messages from the hub to the websocket connection.
func (c *SocketClient) writePump() {
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
	client := &SocketClient{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// run writing to client in goroutine
	go client.writePump()
}

// SocketHub maintains the set of active clients and broadcasts messages to the
// clients.
type SocketHub struct {
	// Registered clients.
	clients map[*SocketClient]bool

	// Register requests from the clients.
	register chan *SocketClient

	// Unregister requests from clients.
	unregister chan *SocketClient

	// status stream
	statusStream chan *Status
}

func NewSocketHub(
	// status *Status
	) *SocketHub {
	// Attach a goroutine that will push meter status information
	// periodically
	// var statusstream = make(chan *Status)
	// go func() {
	// 	for {
	// 		time.Sleep(SECONDS_BETWEEN_STATUSUPDATE * time.Second)
	// 		status.Update()
	// 		statusstream <- status
	// 	}
	// }()

	return &SocketHub{
		register:     make(chan *SocketClient),
		unregister:   make(chan *SocketClient),
		clients:      make(map[*SocketClient]bool),
		// statusStream: statusstream,
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
