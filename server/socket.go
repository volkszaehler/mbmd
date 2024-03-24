package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer
	socketWriteWait = 10 * time.Second

	// Frequency at which status updates are sent
	statusFrequency = 1 * time.Second
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
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

	// status channel
	status *Status
}

// NewSocketHub creates a web socket hub that distributes meter status and
// query results for the ui or other clients
func NewSocketHub(status *Status) *SocketHub {
	return &SocketHub{
		register:   make(chan *SocketClient),
		unregister: make(chan *SocketClient),
		clients:    make(map[*SocketClient]bool),
		status:     status,
	}
}

func (h *SocketHub) broadcast(i interface{}) {
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

// Run starts data and status distribution
func (h *SocketHub) Run(in <-chan QuerySnip) {
	// Periodically push meter status information
	statusChannel := make(chan *Status)
	go func() {
		for {
			time.Sleep(statusFrequency)
			statusChannel <- h.status
		}
	}()

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
			h.broadcast(&obj)
		case obj := <-statusChannel:
			h.broadcast(obj)
		}
	}
}
