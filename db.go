package main

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var Db *gorm.DB

type DbConfig struct {
	Host     string
	Port     int32
	User     string
	Name     string
	Password string
	SSLMode  string
}

// Open DB Connection
func (dbConfig DbConfig) Open() {
	db, err := gorm.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
			dbConfig.Host,
			dbConfig.Port,
			dbConfig.User,
			dbConfig.Name,
			dbConfig.Password,
			dbConfig.SSLMode,
		),
	)
	if err != nil {
		log.Fatal(err)
		panic("Failed to connect to database")
	}

	db.AutoMigrate(&Expense{})
	Db = db
}
