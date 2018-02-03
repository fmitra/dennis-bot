package postgres

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var Db *gorm.DB

type Config struct {
	Host     string
	Port     int32
	User     string
	Name     string
	Password string
	SSLMode  string
}

// Open DB Connection
func (config Config) Open() {
	db, err := gorm.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
			config.Host,
			config.Port,
			config.User,
			config.Name,
			config.Password,
			config.SSLMode,
		),
	)
	if err != nil {
		log.Fatal(err)
		panic("Failed to connect to database")
	}

	Db = db
}
