package models

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"not null" json:"username"`
	Email    string `gorm:"unique;not null" json:"email"`
	Password string `gorm:"not null" json:"-"`
}

func UserNew() *User {
	return &User{}
}
