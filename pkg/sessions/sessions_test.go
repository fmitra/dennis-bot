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
	file := "../../config/config.json"
	configFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer configFile.Close()
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
			UserID    string
			UserEmail string
		}

		session := GetSession()
		userMock := UserMock{
			"userID",
			"userEmail",
		}

		expiresIn := 60
		var cachedUser UserMock
		session.Set("userID", userMock, expiresIn)
		session.Get("userID", &cachedUser)

		assert.Equal(t, userMock, cachedUser)

		session.Delete("userID")
		err := session.Get("userID", &cachedUser)
		assert.EqualError(t, err, "no session found")
	})

	t.Run("Returns error if not found", func(t *testing.T) {
		type UserMock struct {
			UserID    string
			UserEmail string
		}

		session := GetSession()
		userMock := UserMock{
			"userID",
			"userEmail",
		}

		expiresIn := 60
		var wanted UserMock
		session.Set("userID", userMock, expiresIn)
		err := session.Get("nonExistentUser", &wanted)

		assert.EqualError(t, err, "no session found")
	})
}
