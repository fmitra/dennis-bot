package users

import (
	"errors"
	"log"

	"github.com/jinzhu/gorm"
	// Register SQL driver for DB
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// UserManager exposes methods to interface with a User in our
// database.
type UserManager struct {
	db *gorm.DB
}

// NewUserManager returns a UserManager.
func NewUserManager(db *gorm.DB) *UserManager {
	return &UserManager{
		db: db,
	}
}

// Save saves a User into our DB.
func (m *UserManager) Save(user *User) error {
	var existingUser User
	err := m.db.Where("telegram_id = ?", user.TelegramID).First(&existingUser).Error
	if err == nil {
		log.Printf("models: attempting insert record with existing telegram id - %v", user)
		return errors.New("telegram ID already exists")
	}

	m.db.Create(user)
	return nil
}

// GetByTelegramID return's a bot User based on their Telegram ID.
func (m *UserManager) GetByTelegramID(tID uint) User {
	var user User
	m.db.Where("telegram_id = ?", tID).First(&user)
	return user
}
