package users

import (
	"errors"
	"log"

	"github.com/jinzhu/gorm"
	// Register SQL driver for DB
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/fmitra/dennis-bot/pkg/utils"
)

// UserManager exposes methods to interface with a User in our
// database.
type UserManager struct {
	db *gorm.DB
}

// SettingManager exposes methods to interface with a Setting in our
// database.
type SettingManager struct {
	db *gorm.DB
}

// NewUserManager returns a UserManager.
func NewUserManager(db *gorm.DB) *UserManager {
	return &UserManager{
		db: db,
	}
}

// NewSettingManager returns a SettingManager.
func NewSettingManager(db *gorm.DB) *SettingManager {
	return &SettingManager{
		db: db,
	}
}

// Save saves a User into our DB.
func (m *UserManager) Save(user *User) error {
	var existingUser User
	if m.db.Where("telegram_id = ?", user.TelegramID).First(&existingUser).RecordNotFound() {
		m.db.Create(user)
		return nil
	}

	log.Printf("models: attempting insert record with existing telegram id - %v", user)
	return errors.New("telegram ID already exists")
}

// GetByTelegramID return's a bot User based on their Telegram ID.
func (m *UserManager) GetByTelegramID(tID uint) User {
	var user User
	m.db.Where("telegram_id = ?", tID).First(&user)
	return user
}

// UpdateCurrency creates or updates a user's settings with the
// valid currency ISO.
func (m *SettingManager) UpdateCurrency(userID uint, currency string) error {
	currencyISO, err := utils.ParseISO(currency)
	if err != nil {
		return err
	}

	setting := &Setting{
		UserID:   userID,
		Currency: currencyISO,
	}

	var existing Setting
	if m.db.Where("user_id = ?", userID).First(&existing).RecordNotFound() {
		m.db.Create(setting)
		return nil
	}

	tx := m.db.Begin()
	tx.Model(&existing).Update("currency", currencyISO)
	tx.Commit()

	return nil
}

// GetCurrency returns the User's preferred currency ISO.
func (m *SettingManager) GetCurrency(userID uint) string {
	var setting Setting
	if m.db.Where("user_id = ?", userID).First(&setting).RecordNotFound() {
		return "USD"
	}

	return setting.Currency
}
