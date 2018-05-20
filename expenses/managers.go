package expenses

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

type Clock interface {
	Now() time.Time
}

type ExpenseManagerClock struct{}

type ExpenseManager struct {
	clock Clock
	db    *gorm.DB
}

func (em *ExpenseManagerClock) Now() time.Time {
	return time.Now()
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
