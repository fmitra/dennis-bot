package expenses

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/fmitra/dennis/alphapoint"
	"github.com/fmitra/dennis/postgres"
	"github.com/fmitra/dennis/wit"
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

// Creates an expense item from a Wit.ai response
func NewExpense(w wit.WitResponse, userId int) {
	date := w.GetDate()
	amount, currency, _ := w.GetAmount()
	description, _ := w.GetDescription()
	historical := alphapoint.Client.Convert(
		currency,
		"USD",
		amount,
	)

	postgres.Db.Create(&Expense{
		Date:        date,
		Description: description,
		Total:       amount,
		Historical:  historical,
		Currency:    currency,
		UserId:      userId,
	})
}
