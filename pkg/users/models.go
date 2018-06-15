package users

import (
	"log"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// Represents a User account associated with the bot.
// We keep the user model separate from what is provided by the
// Telegram API to budget for the possibility of moving to, or supporting
// additional messenger platforms.
type User struct {
	gorm.Model
	Password   string `gorm:"type:varchar(2000);not null"`
	TelegramID uint   `gorm:"unique_index"`
}

func (u *User) BeforeCreate() error {
	p := []byte(u.Password)
	cost := 10
	hash, err := bcrypt.GenerateFromPassword(p, cost)
	if err != nil {
		log.Printf("users: failed to hash password")
		return err
	}

	u.Password = string(hash)
	return nil
}
