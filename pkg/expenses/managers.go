package expenses

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	// Register SQL driver for DB
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	// MONTH is a one month period of expenses
	MONTH = "month"

	// WEEK is a one week period of expenses
	WEEK = "week"

	// TODAY is a one day period of expenses
	TODAY = "today"
)

// Clock is an interface that provides a Now method.
type Clock interface {
	Now() time.Time
}

// ExpenseManagerClock uses stdlib time to satisfy the Clock interface
type ExpenseManagerClock struct{}

// ExpenseManager exposes methods to interface with an Expense in our database.
type ExpenseManager struct {
	clock Clock
	db    *gorm.DB
}

// Now returns the current time.
func (em *ExpenseManagerClock) Now() time.Time {
	return time.Now()
}

// NewExpenseManager returns an ExpenseManager with a default clock.
func NewExpenseManager(db *gorm.DB) *ExpenseManager {
	return &ExpenseManager{
		db:    db,
		clock: &ExpenseManagerClock{},
	}
}

// Save saves an Expense into our DB.
func (m *ExpenseManager) Save(expense *Expense) error {
	if m.db.NewRecord(expense) {
		m.db.Create(expense)
		return nil
	}

	log.Printf("models: attempting insert record with existing pk - %v", expense)
	return errors.New("expense ID already exists")
}

// ParseTimePeriod parses a string period (ex. month) into a time.Time object.
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

// QueryByPeriod finds all expenses within a specific period.
func (m *ExpenseManager) QueryByPeriod(period string, userID uint) ([]Expense, error) {
	var expenses []Expense

	timePeriod, err := m.ParseTimePeriod(period)
	if err != nil {
		return expenses, err
	}

	query := "user_id = ? AND date >= ?"
	if period == TODAY {
		query = "user_id = ? AND date = ?"
	}

	if err = m.db.Where(query, userID, timePeriod).Find(&expenses).Error; err != nil {
		return expenses, err
	}

	return expenses, nil
}

// TotalByPeriod sums the total historical value of a list of Expenses.
func (m *ExpenseManager) TotalByPeriod(period string, userID uint, pk rsa.PrivateKey) (float64, error) {
	expenseTotal := 0.0
	expenses, err := m.QueryByPeriod(period, userID)
	if err != nil {
		return expenseTotal, err
	}

	// We cannot sum the DB column because expense history is stored
	// as an encrypted string in the DB, so we must first query for the relevant
	// records, decrypt, and sum it ourselves
	for _, expense := range expenses {
		expense.Decrypt(pk)
		amount, err := strconv.ParseFloat(expense.Historical, 64)
		if err != nil {
			return 0.0, err
		}

		expenseTotal += amount
	}
	return expenseTotal, nil
}
