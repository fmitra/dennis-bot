package users

import (
	"github.com/jinzhu/gorm"

	"github.com/fmitra/dennis-bot/pkg/crypto"
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

func (u *User) IsPasswordValid(password string) bool {
	return crypto.ValidateHash(u.Password, password)
}
