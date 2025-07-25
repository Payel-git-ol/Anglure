package Complite

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type ChatMessage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	Message   string    `json:"message"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type ChatPageData struct {
	Profile  UserProfileData
	Messages []ChatMessage
}

type UserProfileData struct {
	Name  string
	Email string
}

func (ChatMessage) TableName() string {
	return "chat_messages"
}

type Client struct {
	conn *websocket.Conn
	send chan ChatMessage
}

type ChatHub struct {
	clients    map[*Client]bool
	broadcast  chan ChatMessage
	register   chan *Client
	unregister chan *Client
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var hub = ChatHub{
	broadcast:  make(chan ChatMessage),
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
		result := db.Delete(&ChatMessage{}, id)
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

		var messages []ChatMessage
		db.Order("created_at desc").Find(&messages)

		// Получаем профиль пользователя из куков (или по-другому)
		nameCookie, err := r.Cookie("user_name")
		if err != nil {
			nameCookie = &http.Cookie{Value: "Гость"}
		}
		emailCookie, err := r.Cookie("user_email")
		if err != nil {
			emailCookie = &http.Cookie{Value: "guest@example.com"}
		}

		data := ChatPageData{
			Profile: UserProfileData{
				Name:  nameCookie.Value,
				Email: emailCookie.Value,
			},
			Messages: messages,
		}

		tmpl, err := template.ParseFiles("template/Chat.html")
		if err != nil {
			http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
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
		var user UserRegister
		result := db.Where("email = ?", emailCookie.Value).First(&user)
		if result.Error == nil {
			userID = user.ID
		}

		newMessage := ChatMessage{
			UserID:    userID,
			Message:   r.FormValue("message"),
			Email:     emailCookie.Value,
			Name:      nameCookie.Value,
			CreatedAt: time.Now(),
		}

		result = db.Create(&newMessage)
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
		send: make(chan ChatMessage, 256),
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

		var msg ChatMessage
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
