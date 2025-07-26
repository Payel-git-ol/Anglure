package chat

import (
	"Angular/internal/DataBase/postgres"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
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

func init() {
	go hub.run()
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

func HandleChat(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "DELETE":
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Нет id", http.StatusBadRequest)
			return
		}
		result := postgres.Db.Delete(&postgres.ChatMessage{}, id)
		if result.Error != nil {
			http.Error(w, "Ошибка удаления", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

	case "GET":
		if websocket.IsWebSocketUpgrade(r) {
			handleWebSocket(w, r)
			return
		}

		var messages []postgres.ChatMessage
		postgres.Db.Order("created_at desc").Find(&messages)

		// Получаем профиль пользователя из куков (или по-другому)
		nameCookie, err := r.Cookie("user_name")
		if err != nil {
			nameCookie = &http.Cookie{Value: "Гость"}
		}
		emailCookie, err := r.Cookie("user_email")
		if err != nil {
			emailCookie = &http.Cookie{Value: "guest@example.com"}
		}

		data := postgres.ChatPageData{
			Profile: postgres.UserProfileData{
				Name:  nameCookie.Value,
				Email: emailCookie.Value,
			},
			Messages: messages,
		}

		tmpl, err := template.ParseFiles("web/templates/ChatTemplates/Chat.html")
		if err != nil {
			http.Error(w, "Ошибка загрузки шаблона Chat", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, data)
		tmpl.Execute(w, messages)

	case "POST":
		emailCookie, err := r.Cookie("user_email")
		if err != nil {
			emailCookie = &http.Cookie{Value: "guest@example.com"}
		}
		nameCookie, err := r.Cookie("user_name")
		if err != nil {
			nameCookie = &http.Cookie{Value: "Guest"}
		}

		var userID uint
		var user postgres.UserRegister
		result := postgres.Db.Where("email = ?", emailCookie.Value).First(&user)
		if result.Error == nil {
			userID = user.ID
		}

		newMessage := postgres.ChatMessage{
			UserID:    userID,
			Message:   r.FormValue("message"),
			Email:     emailCookie.Value,
			Name:      nameCookie.Value,
			CreatedAt: time.Now(),
		}

		result = postgres.Db.Create(&newMessage)
		if result.Error != nil {
			log.Printf("Ошибка сохранения: %v", result.Error)
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		hub.broadcast <- newMessage
		w.WriteHeader(http.StatusOK) // или можно вернуть JSON
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
