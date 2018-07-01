package mocks

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	// Register SQL driver for DB related tests
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"

	"github.com/fmitra/dennis-bot/config"
	"github.com/fmitra/dennis-bot/pkg/crypto"
	"github.com/fmitra/dennis-bot/pkg/sessions"
)

// We use package scoped variables here to implement testEnv as a singleton
// across the different test suites.
var (
	testEnv *TestEnv
	once    sync.Once
)

// TestEnv is the working environmenet for test suites.
type TestEnv struct {
	Db     *gorm.DB
	Cache  *sessions.Client
	Config config.AppConfig
}

// GetTestEnv returns TestEnv as a singleton.
func GetTestEnv(configFile string) *TestEnv {
	// We only need a single Postgres/Redis connection
	// to share across tests
	once.Do(func() {
		crypto.InitializeGob()

		appConfig := config.LoadConfig(configFile)
		db := getDb(appConfig)
		cache := getSessions(appConfig)
		testEnv = &TestEnv{
			Db:     db,
			Cache:  cache,
			Config: appConfig,
		}
	})

	return testEnv
}

// CreateTestUser creates default test user with telegram user ID specified in mocks.
func CreateTestUser(db *gorm.DB, userID uint) {
	telegramID := userID
	if userID == 0 {
		telegramID = TestUserID
	}
	password, _ := bcrypt.GenerateFromPassword([]byte("my-password"), 10)
	user := &user{
		TelegramID: telegramID,
		Password:   string(password),
	}
	db.Create(user)
}

// CleanUpEnv cleans common DB and cached objects. Intended to be run after any test
// suite with a DB dependency.
func CleanUpEnv(testEnv *TestEnv) {
	defaultUserCache := fmt.Sprintf("%s_conversation", strconv.Itoa(int(TestUserID)))
	defaultPassCache := fmt.Sprintf("%s_password", strconv.Itoa(int(TestUserID)))
	defaultCurrencyCache := "SGD_USD"
	testEnv.Cache.Delete(defaultUserCache)
	testEnv.Cache.Delete(defaultPassCache)
	testEnv.Cache.Delete(defaultCurrencyCache)

	tx := testEnv.Db.Begin()
	tx.Exec("DELETE FROM expenses;")
	tx.Exec("DELETE FROM users;")
	tx.Exec("DELETE FROM settings;")
	tx.Commit()
}

// Duplicate of the users pkg model. We define this here to prevent circular
// imports when creating the test environment.
type user struct {
	gorm.Model
	Password   string `gorm:"type:varchar(2000);not null"`
	TelegramID uint   `gorm:"unique_index"`
	PublicKey  string `gorm:"type:varchar(2500);not null"`
	PrivateKey string `gorm:"type:varchar(2500);not null"`
}

// Duplicate of the users pkg model. We define this here to prevent circular
// imports when creating the test environment.
type setting struct {
	gorm.Model
	Currency string `gorm:"type:varchar(30)"`
	User     user
	UserID   uint `gorm:"unique_index"`
}

// Duplicate of the expense pkg model. We define this here to prevent circular
// imports when creating the test environment.
type expense struct {
	gorm.Model
	Date        time.Time `gorm:"index;not null"` // Date the expense was made
	Description string    `gorm:"not null"`       // Description of the expense
	Total       string    `gorm:"not null"`       // Total amount paid for the expense
	Historical  string    // Historical USD value of the total
	Currency    string    `gorm:"not null"`         // Currency ISO of the total
	Category    string    `gorm:"type:varchar(30)"` // Category of the expense
	User        user
	UserID      uint
}

// getSessions returns sessions cache for test environment.
func getSessions(cacheConfig config.AppConfig) *sessions.Client {
	client, _ := sessions.NewClient(sessions.Config{
		Host:     cacheConfig.Redis.Host,
		Port:     cacheConfig.Redis.Port,
		Password: cacheConfig.Redis.Password,
		Db:       cacheConfig.Redis.Db,
	})
	return client
}

// getDb returns a DB for the test environment with initial migrations.
func getDb(dbConfig config.AppConfig) *gorm.DB {
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
	if err != nil {
		log.Panicf("test environment: database connection failed - %s", err)
	}

	db.AutoMigrate(&user{}, &expense{}, &setting{})
	return db
}
