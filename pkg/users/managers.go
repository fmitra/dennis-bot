package users

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type UserManager struct {
	db *gorm.DB
}

func NewUserManager(db *gorm.DB) *UserManager {
	return &UserManager{
		db: db,
	}
}

func (m *UserManager) Save(user *User) bool {
	var existingUser User
	m.db.Where("telegram_id = ?", user.TelegramID).First(&existingUser)
	noId := int(0)
	if int(existingUser.ID) != noId {
		log.Printf("models: attempting insert record with existing pk - %s", user)
		return false
	}

	if m.db.NewRecord(user) {
		m.db.Create(user)
		return true
	}

	return false
}

func (m *UserManager) GetByTelegramId(telegramUserId uint) User {
	var user User
	m.db.Where("telegram_id = ?", telegramUserId).First(&user)
	return user
}
