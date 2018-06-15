package expenses

import (
	"fmt"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stretchr/testify/assert"

	"github.com/fmitra/dennis-bot/config"
	"github.com/fmitra/dennis-bot/pkg/users"
	mocks "github.com/fmitra/dennis-bot/test"
)

func DeleteTestUserExpenses(db *gorm.DB) {
	user := GetTestUser(db)
	db.Where("user_id = ?", user.ID).Unscoped().Delete(Expense{})
}

func GetTestUser(db *gorm.DB) users.User {
	var user users.User
	db.Where("telegram_id = ?", mocks.TestUserId).First(&user)
	return user
}

func GetDb() (*gorm.DB, error) {
	dbConfig := config.LoadConfig("../../config/config.json")
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
	// TODO Set up proper teardown/setup handling for DB related tests
	db.AutoMigrate(&users.User{}, &Expense{})

	password := "my-password"
	user := &users.User{
		TelegramID: mocks.TestUserId,
		Password: password,
	}
	db.Create(user)

	return db, err
}

func BatchCreateExpenses(db *gorm.DB, firstEntryDate time.Time, totalEntries int) {
	var user users.User
	db.Where("telegram_id = ?", mocks.TestUserId).First(&user)

	for days := 1; days <= 10; days++ {
		createdAt := firstEntryDate.AddDate(0, 0, days)
		expense := &Expense{
			Date:        createdAt,
			Description: "Food",
			Total:       26.31,
			Historical:  20.25,
			Currency:    "SGD",
			Category:    "",
			User:        user,
		}
		db.Create(expense)

		// We cannot define the creation date on the initial
		// insert so we need to update the record immediately after
		expense.CreatedAt = createdAt
		db.Save(&expense)
	}
}

func TestManager(t *testing.T) {
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

		user := GetTestUser(db)
		expense := &Expense{
			Date:        time.Now(),
			Description: "Food",
			Total:       26.30,
			Historical:  20.25,
			Currency:    "SGD",
			Category:    "",
			User:        user,
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
		mockTime := &mocks.MockTime{currentTime}
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

		user := GetTestUser(db)
		for _, test := range testCases {
			expenses, err := expenseManager.QueryByPeriod(test.input, user.ID)
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
		mockTime := &mocks.MockTime{currentTime}
		expenseManager := &ExpenseManager{
			db:    db,
			clock: mockTime,
		}

		firstEntryDate := time.Date(2018, 3, 8, 0, 0, 0, 0, time.UTC)
		BatchCreateExpenses(db, firstEntryDate, 5)

		user := GetTestUser(db)
		_, errorMessage := expenseManager.QueryByPeriod("some-date", user.ID)
		assert.EqualError(t, errorMessage, "some-date is an invalid period")
	})

	t.Run("It should parse a string option to time", func(t *testing.T) {
		db, err := GetDb()
		defer db.Close()
		assert.NoError(t, err)

		// Clear out previous test data before we run any queries
		DeleteTestUserExpenses(db)

		currentTime := time.Date(2018, 3, 12, 0, 0, 0, 0, time.UTC)
		mockTime := &mocks.MockTime{currentTime}
		expenseManager := &ExpenseManager{
			db:    db,
			clock: mockTime,
		}

		var testCases = []struct {
			input    string
			expected time.Time
		}{
			{"month", time.Date(2018, 3, 1, 0, 0, 0, 0, time.UTC)},
			{"week", time.Date(2018, 3, 11, 0, 0, 0, 0, time.UTC)},
			{"today", time.Date(2018, 3, 12, 0, 0, 0, 0, time.UTC)},
		}

		for _, test := range testCases {
			timePeriod, err := expenseManager.ParseTimePeriod(test.input)
			assert.NoError(t, err)
			assert.Equal(t, test.expected, timePeriod)
		}
	})

	t.Run("It should sum historical totals by period", func(t *testing.T) {
		db, err := GetDb()
		defer db.Close()
		assert.NoError(t, err)

		// Clear out previous test data before we run any queries
		DeleteTestUserExpenses(db)

		currentTime := time.Date(2018, 3, 12, 0, 0, 0, 0, time.UTC)
		mockTime := &mocks.MockTime{currentTime}
		expenseManager := &ExpenseManager{
			db:    db,
			clock: mockTime,
		}

		firstEntryDate := time.Date(2018, 3, 8, 0, 0, 0, 0, time.UTC)
		BatchCreateExpenses(db, firstEntryDate, 10)

		var testCases = []struct {
			input    string
			expected float64
		}{
			{"month", 202.5},
			{"week", 162.0},
			{"today", 20.25},
		}

		user := GetTestUser(db)
		for _, test := range testCases {
			total, err := expenseManager.TotalByPeriod(test.input, user.ID)
			assert.NoError(t, err)
			assert.Equal(t, test.expected, total)
		}
	})
}
