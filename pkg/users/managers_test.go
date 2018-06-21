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

func (suite *Suite) ValidatesUserPassword() {
	hashedPassword, _ := crypto.HashText("my-password")
	user := &User{
		Password:   hashedPassword,
		TelegramID: mocks.TestUserID,
	}

	assert.NoError(suite.T(), user.ValidatePassword("my-password"))
	assert.EqualError(suite.T(), user.ValidatePassword("not-my-password"), "")
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
