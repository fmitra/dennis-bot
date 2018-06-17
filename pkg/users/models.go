package users

import (
	"crypto/rsa"
	"errors"

	"github.com/jinzhu/gorm"

	"github.com/fmitra/dennis-bot/pkg/crypto"
)

// Represents a User account associated with the bot.
// We keep the user model separate from what is provided by the
// Telegram API to budget for the possibility of moving to, or supporting
// additional messenger platforms.
type User struct {
	gorm.Model
	TelegramID uint   `gorm:"unique_index"`
	Password   string `gorm:"type:varchar(2000);not null"`
	PublicKey  string `gorm:"type:varchar(2500);not null"`
	PrivateKey string `gorm:"type:varchar(2500);not null"`
}

func (u *User) BeforeCreate(scope *gorm.Scope) error {
	hashedPassword, err := crypto.HashText(u.Password)
	publicKey, privateKey, err := crypto.CreateProtectedKeyPair(u.Password)
	if err != nil {
		return err
	}

	u.Password = hashedPassword
	u.PublicKey = publicKey
	u.PrivateKey = privateKey
	return nil
}

func (u *User) IsPasswordValid(password string) bool {
	return crypto.ValidateHash(u.Password, password)
}

func (u *User) GetPrivateKey(password string) (rsa.PrivateKey, error) {
	if !u.IsPasswordValid(password) {
		return rsa.PrivateKey{}, errors.New("Invalid password")
	}

	key, err := crypto.ParsePrivateKey(u.PrivateKey, password)
	if err != nil {
		return rsa.PrivateKey{}, err
	}

	return key, nil
}

func (u *User) GetPublicKey() (rsa.PublicKey, error) {
	return crypto.ParsePublicKey(u.PublicKey)
}
