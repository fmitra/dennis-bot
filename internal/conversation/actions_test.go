package conversation

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stretchr/testify/assert"

	"github.com/fmitra/dennis-bot/config"
	"github.com/fmitra/dennis-bot/pkg/alphapoint"
	"github.com/fmitra/dennis-bot/pkg/expenses"
	"github.com/fmitra/dennis-bot/pkg/sessions"
	"github.com/fmitra/dennis-bot/pkg/wit"
	mocks "github.com/fmitra/dennis-bot/test"
)

func GetSession(config config.AppConfig) *sessions.Client {
	return sessions.NewClient(sessions.Config{
		config.Redis.Host,
		config.Redis.Port,
		config.Redis.Password,
		config.Redis.Db,
	})
}

func GetDb(config config.AppConfig) (*gorm.DB, error) {
	db, err := gorm.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
			config.Database.Host,
			config.Database.Port,
			config.Database.User,
			config.Database.Name,
			config.Database.Password,
			config.Database.SSLMode,
		),
	)
	db.AutoMigrate(&expenses.Expense{})

	return db, err
}

func TestExpenseCreation(t *testing.T) {
	t.Run("Creates a new expense", func(t *testing.T) {
		configFile := "../../config/config.json"
		appConfig := config.LoadConfig(configFile)
		db, _ := GetDb(appConfig)
		cache := GetSession(appConfig)
		cache.Delete("SGD_USD")

		rawWitResponse := []byte(`{
			"entities": {
				"amount": [
					{ "value": "20 SGD", "confidence": 100.00 }
				],
				"datetime": [
					{ "value": "", "confidence": 100.00 }
				],
				"description": [
					{ "value": "Food", "confidence": 100.00 }
				]
			}
		}`)
		alphapointResponse := `{
			"Realtime Currency Exchange Rate": {
				"5. Exchange Rate": ".7"
			}
		}`
		alphapointServer := mocks.MakeTestServer(alphapointResponse)
		defer alphapointServer.Close()

		var witResponse wit.WitResponse
		json.Unmarshal(rawWitResponse, &witResponse)

		alphapoint.BaseUrl = alphapointServer.URL

		action := &Actions{
			db,
			cache,
			appConfig,
		}

		isCreated := action.CreateNewExpense(witResponse, mocks.TestUserId)
		assert.True(t, isCreated)
	})

	t.Run("Creates a new expense from cache", func(t *testing.T) {
		configFile := "../../config/config.json"
		appConfig := config.LoadConfig(configFile)
		db, _ := GetDb(appConfig)
		cache := GetSession(appConfig)

		cache.Delete("SGD_USD")
		rawWitResponse := []byte(`{
			"entities": {
				"amount": [
					{ "value": "20 SGD", "confidence": 100.00 }
				],
				"datetime": [
					{ "value": "", "confidence": 100.00 }
				],
				"description": [
					{ "value": "Food", "confidence": 100.00 }
				]
			}
		}`)
		alphapointResponse := `{
			"Realtime Currency Exchange Rate": {
				"5. Exchange Rate": ".7"
			}
		}`
		alphapointServer := mocks.MakeTestServer(alphapointResponse)

		var witResponse wit.WitResponse
		json.Unmarshal(rawWitResponse, &witResponse)

		alphapoint.BaseUrl = alphapointServer.URL

		action := &Actions{
			db,
			cache,
			appConfig,
		}

		// Initial call without cache
		action.CreateNewExpense(witResponse, mocks.TestUserId)

		// Second call should not hit server
		alphapointServer.Close()
		isCreated := action.CreateNewExpense(witResponse, mocks.TestUserId)
		assert.True(t, isCreated)
	})
}

func TestExpenseTotal(t *testing.T) {
	t.Run("Gets expense total", func(t *testing.T) {
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

		db.Where("user_id = ?", mocks.TestUserId).
			Unscoped().
			Delete(expenses.Expense{})

		total, err := action.GetExpenseTotal(witResponse, mocks.TestUserId)
		assert.Equal(t, "0.00", total)
		assert.NoError(t, err)
	})

	t.Run("Returns error for invalid period", func(t *testing.T) {
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
		_, err := action.GetExpenseTotal(witResponse, mocks.TestUserId)
		assert.EqualError(t, err, "foo is an invalid period")
	})
}
