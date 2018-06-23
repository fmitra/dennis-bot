// Package users represents a User account with the Bot service.
package users

import (
	"crypto/rsa"
	"errors"
	"log"

	"github.com/jinzhu/gorm"

	"github.com/fmitra/dennis-bot/pkg/crypto"
)

// User is an account associated with the bot service. A Telegram User
// is considered a user of the bot service when they create a password
// and public/private key pair.
type User struct {
	gorm.Model
	TelegramID uint   `gorm:"unique_index"`
	Password   string `gorm:"type:varchar(2000);not null"`
	PublicKey  string `gorm:"type:varchar(2500);not null"`
	PrivateKey string `gorm:"type:varchar(2500);not null"`
}

// Setting describes User specific settings for the bot. For example it
// controls whether the bot returns expense history in the default currency,
// USD or a user specified currency.
type Setting struct {
	gorm.Model
	Currency string `gorm:"type:varchar(30)"`
	User     User
	UserID   uint `gorm:"unique_index"`
}

// BeforeCreate hashes a User's plaintext password and generates
// a public/private key pair on their behalf prior to creating a DB
// record.
func (u *User) BeforeCreate(scope *gorm.Scope) error {
	hashedPassword, err := crypto.HashText(u.Password)
	if err != nil {
		log.Printf("users: failed to hash password")
		return err
	}

	publicKey, privateKey, err := crypto.CreateProtectedKeyPair(u.Password)
	if err != nil {
		log.Printf("users: failed to create protected key pair")
		return err
	}

	u.Password = hashedPassword
	u.PublicKey = publicKey
	u.PrivateKey = privateKey
	return nil
}

// ValidatePassword checks if a plaintext password string validates against
// the User's stored password hash.
func (u *User) ValidatePassword(password string) error {
	return crypto.ValidateHash(u.Password, password)
}

// GetPrivateKey returns the decrypted rsa.PrivateKey of the User.
func (u *User) GetPrivateKey(password string) (rsa.PrivateKey, error) {
	if err := u.ValidatePassword(password); err != nil {
		return rsa.PrivateKey{}, errors.New("invalid password")
	}

	key, err := crypto.ParsePrivateKey(u.PrivateKey, password)
	if err != nil {
		return rsa.PrivateKey{}, err
	}

	return key, nil
}

// GetPublicKey returns the rsa.PublicKey of the User.
func (u *User) GetPublicKey() (rsa.PublicKey, error) {
	return crypto.ParsePublicKey(u.PublicKey)
}
