package conversation

import (
	"encoding/json"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"

	"github.com/fmitra/dennis-bot/config"
	"github.com/fmitra/dennis-bot/pkg/expenses"
	"github.com/fmitra/dennis-bot/pkg/telegram"
	"github.com/fmitra/dennis-bot/pkg/users"
	"github.com/fmitra/dennis-bot/pkg/wit"
	mocks "github.com/fmitra/dennis-bot/test"
)

func GetTestUser(db *gorm.DB) users.User {
	var user users.User
	db.Where("telegram_id = ?", mocks.TestUserId).First(&user)
	return user
}

func DeleteTestUserExpenses(db *gorm.DB) {
	user := GetTestUser(db)
	db.Where("user_id = ?", user.ID).Unscoped().Delete(expenses.Expense{})
}

func TestGetExpenseTotal(t *testing.T) {
	t.Run("Should return a list of possible responses", func(t *testing.T) {
		expenseTotal := &GetExpenseTotal{}
		assert.Equal(t, 1, len(expenseTotal.GetResponses()))
	})

	t.Run("Should return expense total message", func(t *testing.T) {
		MessageMap = mocks.MessageMapMock
		configFile := "../../config/config.json"
		appConfig := config.LoadConfig(configFile)
		db, _ := GetDb(appConfig)
		cache := GetSession(appConfig)
		action := &Actions{
			db,
			cache,
			appConfig,
		}

		rawWitResponse := []byte(`{
			"entities": {
				"amount": [],
				"datetime": [],
				"description": [],
				"total_spent": [
					{ "value": "month", "confidence": 100.00 }
				]
			}
		}`)
		var witResponse wit.WitResponse
		json.Unmarshal(rawWitResponse, &witResponse)

		var incMessage telegram.IncomingMessage
		message := mocks.GetMockMessage("")
		json.Unmarshal(message, &incMessage)

		DeleteTestUserExpenses(db)
		expenseTotal := &GetExpenseTotal{
			Context{
				Step:        0,
				WitResponse: witResponse,
				IncMessage:  incMessage,
			},
			action,
		}
		response, step := expenseTotal.Respond()
		assert.Equal(t, BotResponse("You spent 0.00"), response)
		assert.Equal(t, -1, step)
	})

	t.Run("Should return expense total error message", func(t *testing.T) {
		MessageMap = mocks.MessageMapMock
		configFile := "../../config/config.json"
		appConfig := config.LoadConfig(configFile)
		db, _ := GetDb(appConfig)
		cache := GetSession(appConfig)
		action := &Actions{
			db,
			cache,
			appConfig,
		}

		rawWitResponse := []byte(`{
			"entities": {
				"amount": [],
				"datetime": [],
				"description": [],
				"total_spent": [
					{ "value": "foo", "confidence": 100.00 }
				]
			}
		}`)
		var witResponse wit.WitResponse
		json.Unmarshal(rawWitResponse, &witResponse)

		var incMessage telegram.IncomingMessage
		message := mocks.GetMockMessage("")
		json.Unmarshal(message, &incMessage)

		DeleteTestUserExpenses(db)
		expenseTotal := &GetExpenseTotal{
			Context{
				Step:        0,
				WitResponse: witResponse,
				IncMessage:  incMessage,
			},
			action,
		}
		response, step := expenseTotal.Respond()
		assert.Equal(t, BotResponse("Whoops!"), response)
		assert.Equal(t, -1, step)
	})
}
