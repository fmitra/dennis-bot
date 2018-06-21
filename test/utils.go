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

// We use global variables here to implement testEnv as a singleton
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
func CreateTestUser(db *gorm.DB) {
	password, _ := bcrypt.GenerateFromPassword([]byte("my-password"), 10)
	user := &user{
		TelegramID: TestUserID,
		Password:   string(password),
	}
	db.Create(user)
}

// CleanUpEnv cleans common DB and cached objects. Intended to be run after any test
// suite with a DB dependency.
func CleanUpEnv(testEnv *TestEnv) {
	// TODO There is a race condition that needs to be investigated where test tear down
	// methods do not finish clearing out the DB before the next suite starts. Setting
	// a lock does not seem to resolve it. For now, we will miitigate the issue with
	// with custom SQL and a sleeper until we find a more elegant solution.
	testEnv.Db.Exec("DELETE FROM expenses e USING users u WHERE e.user_id = u.id and u.telegram_id != 100;")
	testEnv.Db.Exec("DELETE FROM users;")

	defaultUserCache := fmt.Sprintf("%s_conversation", strconv.Itoa(int(TestUserID)))
	defaultPassCache := fmt.Sprintf("%s_password", strconv.Itoa(int(TestUserID)))
	defaultCurrencyCache := "SGD_USD"
	testEnv.Cache.Delete(defaultUserCache)
	testEnv.Cache.Delete(defaultPassCache)
	testEnv.Cache.Delete(defaultCurrencyCache)

	// TODO Clean this up later
	time.Sleep(500 * time.Millisecond)
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
	return sessions.NewClient(sessions.Config{
		cacheConfig.Redis.Host,
		cacheConfig.Redis.Port,
		cacheConfig.Redis.Password,
		cacheConfig.Redis.Db,
	})
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

	db.AutoMigrate(&user{}, &expense{})
	return db
}
