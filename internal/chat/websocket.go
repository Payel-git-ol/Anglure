package chat

import (
	"Angular/internal/DataBase/postgres"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Client struct {
	conn *websocket.Conn
	send chan postgres.ChatMessage
}

type ChatHub struct {
	clients    map[*Client]bool
	broadcast  chan postgres.ChatMessage
	register   chan *Client
	unregister chan *Client
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var hub = ChatHub{
	broadcast:  make(chan postgres.ChatMessage),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	clients:    make(map[*Client]bool),
}

func (h *ChatHub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
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
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Ошибка WebSocket:", err)
		return
	}

	client := &Client{
		conn: conn,
		send: make(chan postgres.ChatMessage, 256),
	}

	hub.register <- client

	go client.writePump()
	go client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		var msg postgres.ChatMessage
		if err := json.Unmarshal(message, &msg); err == nil {
			// обработка входящих сообщений (не используется сейчас)
		}
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			jsonMsg, err := json.Marshal(message)
			if err != nil {
				continue
			}

			c.conn.WriteMessage(websocket.TextMessage, jsonMsg)
		}
	}
}
