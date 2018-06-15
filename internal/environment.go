package internal

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/fmitra/dennis-bot/config"
	"github.com/fmitra/dennis-bot/pkg/expenses"
	"github.com/fmitra/dennis-bot/pkg/sessions"
	"github.com/fmitra/dennis-bot/pkg/telegram"
	"github.com/fmitra/dennis-bot/pkg/users"
)

// Working environment for the application
type Env struct {
	db       *gorm.DB
	cache    sessions.Session
	config   config.AppConfig
	telegram telegram.Telegram
}

func (env *Env) HealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}
}

func (env *Env) Webhook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		bot := &Bot{env}
		go bot.Converse(body)

		w.Write([]byte("received"))
	}
}

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

	cache := sessions.NewClient(sessions.Config{
		config.Redis.Host,
		config.Redis.Port,
		config.Redis.Password,
		config.Redis.Db,
	})

	db.AutoMigrate(&users.User{}, &expenses.Expense{})

	telegram := telegram.NewClient(
		config.Telegram.Token,
		config.BotDomain,
	)

	return &Env{
		db:       db,
		cache:    cache,
		config:   config,
		telegram: telegram,
	}
}
