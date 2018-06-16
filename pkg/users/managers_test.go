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
		TelegramID: mocks.TestUserId,
		Password:   password,
	}
	suite.Env.Db.Create(user)

	queriedUser := manager.GetByTelegramId(mocks.TestUserId)
	assert.Equal(suite.T(), mocks.TestUserId, queriedUser.TelegramID)
}

func (suite *Suite) TestCreatesNewUser() {
	manager := NewUserManager(suite.Env.Db)
	user := &User{
		TelegramID: mocks.TestUserId,
	}
	isCreated := manager.Save(user)
	assert.True(suite.T(), isCreated)
}

func (suite *Suite) TestHashesUserPasswordOnSave() {
	password := "my-password"
	manager := NewUserManager(suite.Env.Db)
	user := &User{
		TelegramID: mocks.TestUserId,
		Password:   password,
	}
	manager.Save(user)
	assert.NotEqual(suite.T(), password, user.Password)
	assert.NotEqual(suite.T(), "", user.Password)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
