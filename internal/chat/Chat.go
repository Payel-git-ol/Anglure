package chat

import (
	"Angular/internal/DataBase/postgres"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func init() {
	go hub.run()
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
