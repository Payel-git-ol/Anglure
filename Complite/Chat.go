package Complite

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // В продакшене нужно реализовать правильную проверку origin
	},
}

type ChatMessage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
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

func (ChatMessage) TableName() string {
	return "chat_messages"
}

var hub = ChatHub{
	broadcast:  make(chan ChatMessage),
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

func init() {
	go hub.run()
}

func HandleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
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
		return
	}

	if r.Method == "GET" {
		// Проверяем, является ли запрос веб-сокетом
		if websocket.IsWebSocketUpgrade(r) {
			handleWebSocket(w, r)
			return
		}

		// Обычный HTTP GET запрос
		var messages []ChatMessage
		db.Order("created_at desc").Find(&messages)

		tmpl, err := template.ParseFiles("template/Chat.html")
		if err != nil {
			http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, messages)
		return
	}

	if r.Method == "POST" {
		msg := r.FormValue("message")

		var userID uint = 0
		if email != "" {
			var user UserRegister
			result := db.Where("email = ?", email).First(&user)
			if result.Error == nil {
				userID = user.ID
			}
		}

		newMessage := ChatMessage{
			UserID:    userID,
			Message:   msg,
			CreatedAt: time.Now(),
		}

		result := db.Create(&newMessage)
		if result.Error != nil {
			log.Printf("Ошибка сохранения: %v", result.Error)
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		// Отправляем сообщение всем подключенным клиентам
		hub.broadcast <- newMessage

		http.Redirect(w, r, "/chat", http.StatusSeeOther)
		return
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Ошибка при обновлении до веб-сокета:", err)
		return
	}

	client := &Client{
		conn: conn,
		send: make(chan ChatMessage, 256),
	}

	hub.register <- client

	// Запускаем горутины для чтения и записи
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
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Ошибка: %v", err)
			}
			break
		}

		// Обработка входящих сообщений (если нужно)
		var msg ChatMessage
		if err := json.Unmarshal(message, &msg); err == nil {
			// Можно обработать сообщение от клиента
		}
	}
}

func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				// Канал закрыт
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			jsonMsg, err := json.Marshal(message)
			if err != nil {
				log.Println("Ошибка маршалинга сообщения:", err)
				continue
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, jsonMsg); err != nil {
				log.Println("Ошибка отправки сообщения:", err)
				return
			}
		}
	}
}
