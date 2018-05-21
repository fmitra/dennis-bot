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

type ExpenseTotal struct {
	Total float64
}

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

// Parses a descriptive time period (month, week) into a time.Time object
// time.Time object
func (e *ExpenseManager) ParseTimePeriod(period string) (time.Time, error) {
	today := e.clock.Now()
	weekday := int(today.Weekday())
	year, month, day := today.Date()

	switch period {
	case MONTH:
		return time.Date(year, month, 1, 0, 0, 0, 0, time.UTC), nil
	case WEEK:
		year, month, day := today.AddDate(0, 0, -weekday).Date()
		return time.Date(year, month, day, 0, 0, 0, 0, time.UTC), nil
	case TODAY:
		return time.Date(year, month, day, 0, 0, 0, 0, time.UTC), nil
	default:
		errorMessage := fmt.Sprintf("%s is an invalid period", period)
		return time.Time{}, errors.New(errorMessage)
	}
}

func (e *ExpenseManager) QueryByPeriod(period string, userId int) ([]Expense, error) {
	var expenses []Expense
	var query string

	timePeriod, err := e.ParseTimePeriod(period)
	if err != nil {
		return expenses, err
	}

	if period == TODAY {
		query = "user_id = ? AND created_at = ?"
	} else {
		query = "user_id = ? AND created_at >= ?"
	}

	result := e.db.Where(query, userId, timePeriod).Find(&expenses)
	if result.Error != nil {
		return expenses, result.Error
	}

	return expenses, nil
}

func (e *ExpenseManager) TotalByPeriod(period string, userId int) (float64, error) {
	var query string
	var expenseTotal ExpenseTotal

	timePeriod, err := e.ParseTimePeriod(period)
	if err != nil {
		return expenseTotal.Total, err
	}

	if period == TODAY {
		query = "user_id = ? AND created_at = ?"
	} else {
		query = "user_id = ? AND created_at >= ?"
	}

	result := e.db.Table("expenses").
		Select("sum(historical) as total").
		Where(query, userId, timePeriod).
		Scan(&expenseTotal)

	if result.Error != nil {
		return expenseTotal.Total, err
	}

	return expenseTotal.Total, nil
}
