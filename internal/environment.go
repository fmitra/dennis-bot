package internal

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	// Register SQL driver for DB
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/fmitra/dennis-bot/config"
	"github.com/fmitra/dennis-bot/pkg/crypto"
	"github.com/fmitra/dennis-bot/pkg/expenses"
	"github.com/fmitra/dennis-bot/pkg/sessions"
	"github.com/fmitra/dennis-bot/pkg/telegram"
	"github.com/fmitra/dennis-bot/pkg/users"
)

// Env is the working environment. It exposes HTTP handlers
// to communicate with the bot and the DB/Cache layer as well
// as application configuration.
type Env struct {
	db       *gorm.DB
	cache    sessions.Session
	config   config.AppConfig
	telegram telegram.Telegram
}

// HealthCheck ensures application is running.
func (env *Env) HealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}
}

// Webhook accepts payload from Telegram API for incoming messages.
func (env *Env) Webhook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		bot := &Bot{env}
		go bot.Converse(body)

		w.Write([]byte("received"))
	}
}

// Start sets the Telegram webhook and starts the HTTP server.
func (env *Env) Start() {
	go env.telegram.SetWebhook()

	// Telegram does not send authentication headers on each request
	// and instead recommends we use their token as the webhook path
	webhookPath := fmt.Sprintf("/%s", env.config.Telegram.Token)

	http.HandleFunc("/healthcheck", env.HealthCheck())
	http.HandleFunc(webhookPath, env.Webhook())

	// Run server
	log.Printf("main: starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// LoadEnv initializes dependencies and attaches them to the environment.
func LoadEnv(config config.AppConfig) *Env {
	db, err := gorm.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
			config.Database.Host,
			config.Database.Port,
			config.Database.User,
			config.Database.Name,
			config.Database.Password,
			config.Database.SSLMode,
		),
	)
	if err != nil {
		log.Panicf("environment: database connection failed - %s", err)
	}

	db.AutoMigrate(
		&users.User{},
		&users.Setting{},
		&expenses.Expense{},
	)

	cache, err := sessions.NewClient(sessions.Config{
		Host:     config.Redis.Host,
		Port:     config.Redis.Port,
		Password: config.Redis.Password,
		Db:       config.Redis.Db,
	})
	if err != nil {
		log.Panicf("environment: redis connection failed - %s", err)
	}

	telegram := telegram.NewClient(
		config.Telegram.Token,
		config.BotDomain,
	)

	crypto.InitializeGob()

	return &Env{
		db:       db,
		cache:    cache,
		config:   config,
		telegram: telegram,
	}
}
