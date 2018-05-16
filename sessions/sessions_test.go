package sessions

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SessionsTestSuite struct {
	suite.Suite
}

type LocalConfig struct {
	Redis struct {
		Host     string `json:"host"`
		Port     int32  `json:"port"`
		Password string `json:"password"`
		Db       int    `json:"db"`
	} `json:"redis"`
}

func (suite *SessionsTestSuite) SetupAllSuite() {
	var config LocalConfig
	file := "../config.json"
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		panic(err)
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)

	Init(Config{
		config.Redis.Host,
		config.Redis.Port,
		config.Redis.Password,
		config.Redis.Db,
	})
}

func (suite *SessionsTestSuite) SetsAndGetsFromSession() {
	type UserMock struct {
		UserId    string
		UserEmail string
	}

	userMock := UserMock{
		"userId",
		"userEmail",
	}

	Set("userId", userMock)
	cachedUser, _ := Get("userId")

	assert.Equal(suite.T(), cachedUser, userMock)
}

func (suite *SessionsTestSuite) ReturnsErrorIfNotFound() {
	type UserMock struct {
		UserId    string
		UserEmail string
	}

	userMock := UserMock{
		"userId",
		"userEmail",
	}

	Set("userId", userMock)
	_, err := Get("nonExistentUser")

	assert.EqualError(suite.T(), err, "No session found")

}

func TestSessionsTestSuite(t *testing.T) {
	suite.Run(t, new(SessionsTestSuite))
}
