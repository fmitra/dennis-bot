package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stretchr/testify/assert"

	"github.com/fmitra/dennis/config"
	"github.com/fmitra/dennis/mocks"
)

func DeleteTestUserExpenses(db *gorm.DB) {
	db.Where("user_id = ?", mocks.TestUserId).Unscoped().Delete(Expense{})
}

func GetDb() (*gorm.DB, error) {
	dbConfig := config.LoadConfig("config/config.json")
	db, err := gorm.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
			dbConfig.Database.Host,
			dbConfig.Database.Port,
			dbConfig.Database.User,
			dbConfig.Database.Name,
			dbConfig.Database.Password,
			dbConfig.Database.SSLMode,
		),
	)

	return db, err
}

func BatchCreateExpenses(db *gorm.DB, firstEntryDate time.Time, totalEntries int) {
	for days := 1; days <= 10; days++ {
		createdAt := firstEntryDate.AddDate(0, 0, days)
		expense := &Expense{
			Date:        time.Now(),
			Description: "Food",
			Total:       26.31,
			Historical:  20.25,
			Currency:    "SGD",
			Category:    "",
			UserId:      mocks.TestUserId,
		}
		db.Create(expense)

		// We cannot define the creation date on the initial
		// insert so we need to update the record immediately after
		expense.CreatedAt = createdAt
		db.Save(&expense)
	}
}

func TestModels(t *testing.T) {
	t.Run("It should return an ExpenseManager", func(t *testing.T) {
		db, err := GetDb()
		defer db.Close()
		assert.NoError(t, err)

		expenseManager := NewExpenseManager(db)
		assert.IsType(t, &ExpenseManager{}, expenseManager)
	})

	t.Run("It should create an expense", func(t *testing.T) {
		db, err := GetDb()
		defer db.Close()
		assert.NoError(t, err)

		expense := &Expense{
			Date:        time.Now(),
			Description: "Food",
			Total:       26.30,
			Historical:  20.25,
			Currency:    "SGD",
			Category:    "",
			UserId:      mocks.TestUserId,
		}
		assert.True(t, db.NewRecord(expense))

		expenseManager := NewExpenseManager(db)
		isCreated := expenseManager.Save(expense)

		assert.True(t, isCreated)
		assert.False(t, db.NewRecord(expense))
	})

	t.Run("It should query expenses by period", func(t *testing.T) {
		db, err := GetDb()
		defer db.Close()
		assert.NoError(t, err)

		// Clear out previous test data before we run any queries
		DeleteTestUserExpenses(db)

		currentTime := time.Date(2018, 3, 12, 0, 0, 0, 0, time.UTC)
		mockTime := &mocks.MockTime{
			CurrentTime: currentTime,
		}
		expenseManager := &ExpenseManager{
			db:    db,
			clock: mockTime,
		}

		firstEntryDate := time.Date(2018, 3, 8, 0, 0, 0, 0, time.UTC)
		BatchCreateExpenses(db, firstEntryDate, 10)

		var testCases = []struct {
			input    string
			expected int
		}{
			{"month", 10},
			{"week", 8},
			{"today", 1},
		}

		for _, test := range testCases {
			expenses, err := expenseManager.QueryByPeriod(test.input, mocks.TestUserId)
			assert.NoError(t, err)
			assert.Equal(t, test.expected, len(expenses))
		}
	})

	t.Run("It should error out when querying by invalid time period", func(t *testing.T) {
		db, err := GetDb()
		defer db.Close()
		assert.NoError(t, err)

		// Clear out previous test data before we run any queries
		DeleteTestUserExpenses(db)

		currentTime := time.Date(2018, 3, 12, 0, 0, 0, 0, time.UTC)
		mockTime := &mocks.MockTime{
			CurrentTime: currentTime,
		}
		expenseManager := &ExpenseManager{
			db:    db,
			clock: mockTime,
		}

		firstEntryDate := time.Date(2018, 3, 8, 0, 0, 0, 0, time.UTC)
		BatchCreateExpenses(db, firstEntryDate, 5)
		_, errorMessage := expenseManager.QueryByPeriod("some-date", mocks.TestUserId)
		assert.EqualError(t, errorMessage, "some-date is an invalid period")
	})
}
