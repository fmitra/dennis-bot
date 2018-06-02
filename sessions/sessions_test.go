package sessions

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type LocalConfig struct {
	Redis struct {
		Host     string `json:"host"`
		Port     int32  `json:"port"`
		Password string `json:"password"`
		Db       int    `json:"db"`
	} `json:"redis"`
}

func GetSession() *Client {
	var config LocalConfig
	file := "../config/config.json"
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		panic(err)
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)

	return NewClient(Config{
		config.Redis.Host,
		config.Redis.Port,
		config.Redis.Password,
		config.Redis.Db,
	})
}

func TestSessions(t *testing.T) {
	t.Run("Sets and gets and deletes from session", func(t *testing.T) {
		type UserMock struct {
			UserId    string
			UserEmail string
		}

		session := GetSession()
		userMock := UserMock{
			"userId",
			"userEmail",
		}

		var cachedUser UserMock
		session.Set("userId", userMock)
		session.Get("userId", &cachedUser)

		assert.Equal(t, userMock, cachedUser)

		session.Delete("userId")
		err := session.Get("userId", &cachedUser)
		assert.EqualError(t, err, "No session found")
	})

	t.Run("Returns error if not found", func(t *testing.T) {
		type UserMock struct {
			UserId    string
			UserEmail string
		}

		session := GetSession()
		userMock := UserMock{
			"userId",
			"userEmail",
		}

		var wanted UserMock
		session.Set("userId", userMock)
		err := session.Get("nonExistentUser", &wanted)

		assert.EqualError(t, err, "No session found")
	})
}
