package Complite

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

type ChatMessage struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint
	Message   string
	CreatedAt time.Time
}

func (ChatMessage) TableName() string {
	return "chat_messages"
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
		var messages []ChatMessage
		db.Order("created_at desc").Find(&messages)

		tmpl, err := template.ParseFiles("template/Chat.html")
		if err != nil {
			http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
			return
		}

		// Передаем сообщения в шаблон
		tmpl.Execute(w, messages)
		return
	}

	if r.Method == "POST" {
		msg := r.FormValue("message")

		// Получаем текущего пользователя, если есть
		var userID uint = 0
		if email != "" {
			var user UserRegister
			result := db.Where("email = ?", email).First(&user)
			if result.Error == nil {
				userID = user.ID
			}
		}

		newMessage := ChatMessage{
			UserID:  userID,
			Message: msg,
		}

		result := db.Create(&newMessage)
		if result.Error != nil {
			log.Printf("Ошибка сохранения: %v", result.Error)
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/chat", http.StatusSeeOther)
		return
	}
}