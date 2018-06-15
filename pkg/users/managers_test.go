package users

import (
	"fmt"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stretchr/testify/assert"

	"github.com/fmitra/dennis-bot/config"
)

func GetDb() (*gorm.DB, error) {
	dbConfig := config.LoadConfig("../../config/config.json")
	db, err := gorm.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
			dbConfig.Database.Host,
			dbConfig.Database.Port,
			dbConfig.Database.User,
			dbConfig.Database.Name,
			dbConfig.Database.Password,
			dbConfig.Database.SSLMode,
		),
	)
	// TODO Set up proper teardown/setup handling for DB related tests
	db.AutoMigrate(User{})

	return db, err
}

func DeleteTestUser(db *gorm.DB, testUserId uint) {
	db.Unscoped().Delete(User{}, "telegram_id = ?", testUserId)
}

func TestManagers(t *testing.T) {
	t.Run("It should return a user by their telegram ID", func(t *testing.T) {
		db, _ := GetDb()
		manager := NewUserManager(db)
		testUserId := uint(4567)
		DeleteTestUser(db, testUserId)
		password := "my-password"
		user := &User{
			TelegramID: testUserId,
			Password: password,
		}
		db.Create(user)

		queriedUser := manager.GetByTelegramId(testUserId)
		assert.Equal(t, testUserId, queriedUser.TelegramID)
	})

	t.Run("It should create a new user", func(t *testing.T) {
		db, _ := GetDb()
		testUserId := uint(4567)
		DeleteTestUser(db, testUserId)
		manager := NewUserManager(db)
		user := &User{
			TelegramID: testUserId,
		}
		isCreated := manager.Save(user)
		DeleteTestUser(db, testUserId)
		assert.True(t, isCreated)
	})

	t.Run("It should store a hash of the user password on save", func(t *testing.T) {
		db, _ := GetDb()
		testUserId := uint(4567)
		DeleteTestUser(db, testUserId)
		password := "my-password"
		manager := NewUserManager(db)
		user := &User{
			TelegramID: testUserId,
			Password: password,
		}
		manager.Save(user)
		assert.NotEqual(t, password, user.Password)
		assert.NotEqual(t, "", user.Password)
	})
}
