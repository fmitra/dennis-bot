package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/fmitra/dennis/config"
	"github.com/fmitra/dennis/sessions"
	"github.com/fmitra/dennis/telegram"
)

// Working environment for the application
type Env struct {
	db     *gorm.DB
	cache  *sessions.Session
	config config.AppConfig
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
		log.Printf("main: incoming message - %s", body)

		bot := &Bot{env}
		go bot.Converse(body)

		w.Write([]byte("received"))
	}
}

func (env *Env) Start() {
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

	cache := sessions.New(sessions.Config{
		config.Redis.Host,
		config.Redis.Port,
		config.Redis.Password,
		config.Redis.Db,
	})

	db.AutoMigrate(&Expense{})

	if err != nil {
		log.Fatal(err)
		panic("Failed to connect to database")
	}

	return &Env{
		db:     db,
		cache:  cache,
		config: config,
	}
}

func main() {
	// Set up the environment
	env := LoadEnv(config.LoadConfig())

	t := telegram.Client(
		env.config.Telegram.Token,
		env.config.BotDomain,
		&http.Client{},
	)
	go t.SetWebhook()

	// Start the application
	env.Start()
}
