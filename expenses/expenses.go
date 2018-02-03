package expenses

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/fmitra/dennis/postgres"
	"github.com/fmitra/dennis/wit"
)

// Describes a tracked expense
type Expense struct {
	gorm.Model
	Date time.Time      // Date the expense was made
	Description string  // Description of the expense
	Total float64       // Total amount paid for the expense
	Historical float64  // Historical USD value of the total
	Currency string     // Currency denomination of the total
	Category string     // Category of the expense
}

// Creates an expense item from a Wit.ai response
func NewExpense(w wit.WitResponse) {
	date := w.GetDate()
	amount, currency, _ := w.GetAmount()
	description, _ := w.GetDescription()

	postgres.Db.Create(&Expense{
		Date: date,
		Description: description,
		Total: amount,
		Currency: currency,
	})
}
