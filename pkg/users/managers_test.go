package users

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

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

func (suite *Suite) AfterTest(suiteName, testName string) {
	mocks.CleanUpEnv(suite.Env)
}

func (suite *Suite) TestReturnsUserByTelegramID() {
	manager := NewUserManager(suite.Env.Db)
	password := "my-password"
	user := &User{
		TelegramID: mocks.TestUserID,
		Password:   password,
	}
	suite.Env.Db.Create(user)

	queriedUser := manager.GetByTelegramID(mocks.TestUserID)
	assert.Equal(suite.T(), mocks.TestUserID, queriedUser.TelegramID)
}

func (suite *Suite) TestCreatesNewUser() {
	manager := NewUserManager(suite.Env.Db)
	user := &User{
		TelegramID: mocks.TestUserID,
	}
	err := manager.Save(user)
	assert.NoError(suite.T(), err)
}

func (suite *Suite) TestUpdateCurrency() {
	password := "my-password"
	user := User{
		TelegramID: uint(400),
		Password:   password,
	}
	suite.Env.Db.Create(&user)

	manager := NewSettingManager(suite.Env.Db)
	err := manager.UpdateCurrency(user.ID, "ABC")
	assert.EqualError(suite.T(), err, "invalid currency")

	err = manager.UpdateCurrency(user.ID, "php")
	assert.NoError(suite.T(), err)

	err = manager.UpdateCurrency(user.ID, "SGD")
	assert.NoError(suite.T(), err)

	var setting Setting
	suite.Env.Db.Where("user_id = ?", user.ID).First(&setting)
	assert.Equal(suite.T(), setting.UserID, user.ID)
	assert.Equal(suite.T(), setting.Currency, "SGD")
}

func (suite *Suite) TestGetCurrency() {
	password := "my-password"
	user := User{
		TelegramID: uint(400),
		Password:   password,
	}
	suite.Env.Db.Create(&user)

	manager := NewSettingManager(suite.Env.Db)
	currency := manager.GetCurrency(user.ID)

	assert.Equal(suite.T(), "USD", currency)

	manager.UpdateCurrency(user.ID, "JPY")
	currency = manager.GetCurrency(user.ID)
	assert.Equal(suite.T(), "JPY", currency)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
