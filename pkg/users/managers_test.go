package users

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/fmitra/dennis-bot/pkg/crypto"
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

func (suite *Suite) ValidatesUserPassword() {
	hashedPassword, _ := crypto.HashText("my-password")
	user := &User{
		Password:   hashedPassword,
		TelegramID: mocks.TestUserId,
	}

	assert.True(suite.T(), user.IsPasswordValid("my-password"))
	assert.False(suite.T(), user.IsPasswordValid("not-my-password"))
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
