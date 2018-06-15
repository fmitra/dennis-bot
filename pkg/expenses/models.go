package expenses

import (
	"time"

	"github.com/jinzhu/gorm"

	"github.com/fmitra/dennis-bot/pkg/users"
)

// Describes a tracked expense
type Expense struct {
	gorm.Model
	Date        time.Time `gorm:"index;not null"` // Date the expense was made
	Description string    `gorm:"not null"`       // Description of the expense
	Total       float64   `gorm:"not null"`       // Total amount paid for the expense
	Historical  float64   // Historical USD value of the total
	Currency    string    `gorm:"type:varchar(5);not null"` // Currency ISO of the total
	Category    string    `gorm:"type:varchar(30)"`         // Category of the expense
	User        users.User
	UserID      uint
}
