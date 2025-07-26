package postgres

import "time"

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
