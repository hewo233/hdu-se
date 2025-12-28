package models

import "os"

type Conversation struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	UserID         uint   `gorm:"not null" json:"user_id"`
	ConversationID string   `gorm:"not null" json:"conversation_id"`
	Name          string `gorm:"not null" json:"title"`
}

func NewConversation() *Conversation {
	return &Conversation{}
}

var CozeToken string

func SetCozeToken(path string) {
	// read token from file
	data, err := os.ReadFile(path)
	if err != nil {
		CozeToken = ""
		return
	}
	CozeToken = string(data)
}
