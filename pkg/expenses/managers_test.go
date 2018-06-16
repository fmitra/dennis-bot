package expenses

import (
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/fmitra/dennis-bot/pkg/users"
	mocks "github.com/fmitra/dennis-bot/test"
)

type Suite struct {
	suite.Suite
	Env *mocks.TestEnv
}

func (suite *Suite) SetupSuite() {
	configFile := "../../config/config.json"
	suite.Env = mocks.GetTestEnv(configFile)
}

func (suite *Suite) BeforeTest(suiteName, testName string) {
	mocks.CreateTestUser(suite.Env.Db)
}

func (suite *Suite) AfterTest(suiteName, testName string) {
	mocks.CleanUpEnv(suite.Env)
}

func GetTestUser(db *gorm.DB) users.User {
	var user users.User
	db.Where("telegram_id = ?", mocks.TestUserId).First(&user)
	return user
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

func (suite *Suite) TestReturnsExpenseManager() {
	expenseManager := NewExpenseManager(suite.Env.Db)
	assert.IsType(suite.T(), &ExpenseManager{}, expenseManager)
}

func (suite *Suite) TestCreateExpense() {
	user := GetTestUser(suite.Env.Db)
	expense := &Expense{
		Date:        time.Now(),
		Description: "Food",
		Total:       26.30,
		Historical:  20.25,
		Currency:    "SGD",
		Category:    "",
		User:        user,
	}
	assert.True(suite.T(), suite.Env.Db.NewRecord(expense))

	expenseManager := NewExpenseManager(suite.Env.Db)
	isCreated := expenseManager.Save(expense)

	assert.True(suite.T(), isCreated)
	assert.False(suite.T(), suite.Env.Db.NewRecord(expense))
}

func (suite *Suite) TestQueryExpensesByPeriod() {
	currentTime := time.Date(2018, 3, 12, 0, 0, 0, 0, time.UTC)
	mockTime := &mocks.MockTime{currentTime}
	expenseManager := &ExpenseManager{
		db:    suite.Env.Db,
		clock: mockTime,
	}

	firstEntryDate := time.Date(2018, 3, 8, 0, 0, 0, 0, time.UTC)
	BatchCreateExpenses(suite.Env.Db, firstEntryDate, 10)

	var testCases = []struct {
		input    string
		expected int
	}{
		{"month", 10},
		{"week", 8},
		{"today", 1},
	}

	user := GetTestUser(suite.Env.Db)
	for _, test := range testCases {
		expenses, err := expenseManager.QueryByPeriod(test.input, user.ID)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), test.expected, len(expenses))
	}
}

func (suite *Suite) TestQueryInvalidPeriodWillError() {
	currentTime := time.Date(2018, 3, 12, 0, 0, 0, 0, time.UTC)
	mockTime := &mocks.MockTime{currentTime}
	expenseManager := &ExpenseManager{
		db:    suite.Env.Db,
		clock: mockTime,
	}

	firstEntryDate := time.Date(2018, 3, 8, 0, 0, 0, 0, time.UTC)
	BatchCreateExpenses(suite.Env.Db, firstEntryDate, 5)

	user := GetTestUser(suite.Env.Db)
	_, errorMessage := expenseManager.QueryByPeriod("some-date", user.ID)
	assert.EqualError(suite.T(), errorMessage, "some-date is an invalid period")

}

func (suite *Suite) TestParseStringOptionToTime() {
	currentTime := time.Date(2018, 3, 12, 0, 0, 0, 0, time.UTC)
	mockTime := &mocks.MockTime{currentTime}
	expenseManager := &ExpenseManager{
		db:    suite.Env.Db,
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
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), test.expected, timePeriod)
	}
}

func (suite *Suite) TestSumsHistoricalTotalsByPeriod() {
	currentTime := time.Date(2018, 3, 12, 0, 0, 0, 0, time.UTC)
	mockTime := &mocks.MockTime{currentTime}
	expenseManager := &ExpenseManager{
		db:    suite.Env.Db,
		clock: mockTime,
	}

	firstEntryDate := time.Date(2018, 3, 8, 0, 0, 0, 0, time.UTC)
	BatchCreateExpenses(suite.Env.Db, firstEntryDate, 10)

	var testCases = []struct {
		input    string
		expected float64
	}{
		{"month", 202.5},
		{"week", 162.0},
		{"today", 20.25},
	}

	user := GetTestUser(suite.Env.Db)
	for _, test := range testCases {
		total, err := expenseManager.TotalByPeriod(test.input, user.ID)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), test.expected, total)
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
