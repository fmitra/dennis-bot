package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	MONTH = "month"
	WEEK  = "week"
	TODAY = "today"
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

type Clock interface {
	Now() time.Time
}

type ExpenseManagerClock struct{}

func (em *ExpenseManagerClock) Now() time.Time {
	return time.Now()
}

type ExpenseManager struct {
	clock Clock
	db    *gorm.DB
}

func NewExpenseManager(db *gorm.DB) *ExpenseManager {
	return &ExpenseManager{
		db:    db,
		clock: &ExpenseManagerClock{},
	}
}

func (e *ExpenseManager) Save(expense *Expense) bool {
	if e.db.NewRecord(expense) {
		e.db.Create(expense)
		return true
	}

	log.Printf("models: attempting insert record with existing pk - %s", e)
	return false
}

func (e *ExpenseManager) QueryByPeriod(period string, userId int) ([]Expense, error) {
	var err error
	var timePeriod time.Time
	var expenses []Expense
	var query string

	today := e.clock.Now()
	weekday := int(today.Weekday())
	year, month, day := today.Date()

	switch period {
	case MONTH:
		timePeriod = time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
		query = "user_id = ? AND created_at >= ?"
	case WEEK:
		delta := today.AddDate(0, 0, -weekday)
		timePeriod = time.Date(year, delta.Month(), delta.Day(), 0, 0, 0, 0, time.UTC)
		query = "user_id = ? AND created_at >= ?"
	case TODAY:
		timePeriod = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		query = "user_id = ? AND created_at = ?"
	default:
		errorMessage := fmt.Sprintf("%s is an invalid period", period)
		err = errors.New(errorMessage)
	}

	if err != nil {
		return expenses, err
	}

	result := e.db.Where(query, userId, timePeriod).Find(&expenses)

	if result.Error != nil {
		return expenses, result.Error
	}

	return expenses, nil
}
