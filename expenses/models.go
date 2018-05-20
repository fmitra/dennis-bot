package expenses

import (
	"time"
)

// Describes a tracked expense
type Expense struct {
	ID          int        `gorm:"primary_key"` // Auto assigned ID
	CreatedAt   time.Time  // Timestamp of DB entry
	UpdatedAt   time.Time  // Timestamp of last save date
	DeletedAt   *time.Time // Timestamp for soft deletion
	Date        time.Time  // Date the expense was made
	Description string     // Description of the expense
	Total       float64    // Total amount paid for the expense
	Historical  float64    // Historical USD value of the total
	Currency    string     // Currency denomination of the total
	Category    string     // Category of the expense
	UserId      int        // Telegram UserId of the expense owner
}
