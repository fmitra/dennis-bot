package main

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Describes a tracked expense
type Expense struct {
	gorm.Model
	Date        time.Time // Date the expense was made
	Description string    // Description of the expense
	Total       float64   // Total amount paid for the expense
	Historical  float64   // Historical USD value of the total
	Currency    string    // Currency denomination of the total
	Category    string    // Category of the expense
	UserId      int       // Telegram UserId of the expense owner
}
