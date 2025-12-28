package models

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"not null" json:"username"`
	Email    string `gorm:"unique;not null" json:"email"`
	Password string `gorm:"not null" json:"-"`
	Conversations []Conversation `gorm:"foreignKey:UserID" json:"conversations"`
}

func UserNew() *User {
	return &User{}
}
