package expenses

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"log"
	"strconv"
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

func (m *ExpenseManager) Save(expense *Expense) bool {
	if m.db.NewRecord(expense) {
		m.db.Create(expense)
		return true
	}

	log.Printf("models: attempting insert record with existing pk - %s", expense)
	return false
}

// Parses a descriptive time period (month, week) into a time.Time object
// time.Time object
func (m *ExpenseManager) ParseTimePeriod(period string) (time.Time, error) {
	today := m.clock.Now()
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

func (m *ExpenseManager) QueryByPeriod(period string, userId uint) ([]Expense, error) {
	var expenses []Expense

	timePeriod, err := m.ParseTimePeriod(period)
	if err != nil {
		return expenses, err
	}

	query := "user_id = ? AND date >= ?"
	if period == TODAY {
		query = "user_id = ? AND date = ?"
	}

	result := m.db.Where(query, userId, timePeriod).Find(&expenses)
	if result.Error != nil {
		return expenses, result.Error
	}

	return expenses, nil
}

func (m *ExpenseManager) TotalByPeriod(period string, userId uint, privateKey rsa.PrivateKey) (float64, error) {
	expenseTotal := 0.0
	expenses, err := m.QueryByPeriod(period, userId)
	if err != nil {
		return expenseTotal, err
	}

	// We cannot sum the DB column because expense history is stored
	// as an encrypted string in the DB, so we must first query for the relevant
	// records, decrypt, and sum it ourselves
	for _, expense := range expenses {
		expense.Decrypt(privateKey)
		amount, err := strconv.ParseFloat(expense.Historical, 64)
		if err != nil {
			return 0.0, err
		}

		expenseTotal += amount
	}
	return expenseTotal, nil
}
