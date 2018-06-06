package wit

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/kierdavis/dateparser"
	"github.com/stretchr/testify/assert"
)

func getWitResponse(raw []byte) WitResponse {
	var witResponse WitResponse
	json.Unmarshal(raw, &witResponse)
	return witResponse
}

func TestWitParser(t *testing.T) {
	t.Run("Returns unknown intent for empty response", func(t *testing.T) {
		witResponse := &WitResponse{}
		assert.Equal(t, "default", witResponse.GetIntent())
	})

	t.Run("Returns amount", func(t *testing.T) {
		witResponse := getWitResponse([]byte(`
			{
				"entities": {
					"amount": [
						{ "value": "20 USD", "confidence": 100.00 }
					],
					"datetime": [
						{ "value": "", "confidence": 100.00 }
					],
					"description": [
						{ "value": "Food", "confidence": 100.00 }
					]
				}
			}
		`))

		amount, currency, _ := witResponse.GetAmount()
		assert.Equal(t, 20.0, amount)
		assert.Equal(t, "USD", currency)
	})

	t.Run("Returns error if amount is empty", func(t *testing.T) {
		witResponse := getWitResponse([]byte(`
			{
				"entities": {
					"amount": [],
					"datetime": [
						{ "value": "", "confidence": 100.00 }
					],
					"description": [
						{ "value": "Food", "confidence": 100.00 }
					]
				}
			}
		`))

		amount, currency, err := witResponse.GetAmount()
		assert.Equal(t, 0.0, amount)
		assert.Equal(t, "", currency)
		assert.EqualError(t, err, "No amount")
	})

	t.Run("Returns error if amount is blank", func(t *testing.T) {
		witResponse := getWitResponse([]byte(`
			{
				"entities": {
					"amount": [
						{ "value": "", "confidence": 100.00 }
					],
					"datetime": [
						{ "value": "", "confidence": 100.00 }
					],
					"description": [
						{ "value": "Food", "confidence": 100.00 }
					]
				}
			}
		`))

		amount, currency, err := witResponse.GetAmount()
		assert.Equal(t, 0.0, amount)
		assert.Equal(t, "", currency)
		assert.EqualError(t, err, "Invalid amount")
	})

	t.Run("Returns error if currency is blank", func(t *testing.T) {
		witResponse := getWitResponse([]byte(`
			{
				"entities": {
					"amount": [
						{ "value": "20", "confidence": 100.00 }
					],
					"datetime": [
						{ "value": "", "confidence": 100.00 }
					],
					"description": [
						{ "value": "Food", "confidence": 100.00 }
					]
				}
			}
		`))

		amount, currency, err := witResponse.GetAmount()
		assert.Equal(t, 0.0, amount)
		assert.Equal(t, "", currency)
		assert.EqualError(t, err, "Invalid amount")
	})

	t.Run("Returns description", func(t *testing.T) {
		witResponse := getWitResponse([]byte(`
			{
				"entities": {
					"amount": [],
					"datetime": [],
					"description": [
						{ "value": "Food", "confidence": 100.00 }
					]
				}
			}
		`))

		description, _ := witResponse.GetDescription()
		assert.Equal(t, "Food", description)
	})

	t.Run("Returns error if description is blank", func(t *testing.T) {
		witResponse := getWitResponse([]byte(`
			{
				"entities": {
					"amount": [],
					"datetime": [],
					"description": []
				}
			}
		`))

		description, err := witResponse.GetDescription()
		assert.Equal(t, "", description)
		assert.EqualError(t, err, "No description")
	})

	t.Run("Returns date", func(t *testing.T) {
		parser := &dateparser.Parser{}
		expectedDate, _ := parser.Parse("2018/10/10")
		witResponse := getWitResponse([]byte(`
			{
				"entities": {
					"amount": [],
					"datetime": [
						{ "value": "2018/10/10", "confidence": 100.00 }
					],
					"description": []
				}
			}
		`))

		date := witResponse.GetDate()
		assert.Equal(t, expectedDate, date)
	})

	t.Run("Returns default date if blank", func(t *testing.T) {
		witResponse := getWitResponse([]byte(`
			{
				"entities": {
					"amount": [],
					"datetime": [],
					"description": []
				}
			}
		`))

		date := witResponse.GetDate()
		assert.IsType(t, time.Time{}, date)
	})

	t.Run("Infers tracking", func(t *testing.T) {
		witResponse := getWitResponse([]byte(`
			{
				"entities": {
					"amount": [
						{ "value": "20 USD", "confidence": 100.00 }
					],
					"datetime": [
						{ "value": "", "confidence": 100.00 }
					],
					"description": [
						{ "value": "Food", "confidence": 100.00 }
					]
				}
			}
		`))

		isTracking, err := witResponse.IsTracking()
		assert.True(t, isTracking)
		assert.NoError(t, err)
	})

	t.Run("Fails to infer tracking if missing amount", func(t *testing.T) {
		witResponse := getWitResponse([]byte(`
			{
				"entities": {
					"amount": [],
					"datetime": [
						{ "value": "", "confidence": 100.00 }
					],
					"description": [
						{ "value": "Food", "confidence": 100.00 }
					]
				}
			}
		`))

		isTracking, err := witResponse.IsTracking()
		assert.False(t, isTracking)
		assert.EqualError(t, err, "No amount")
	})

	t.Run("Fails to infer tracking if missing description", func(t *testing.T) {
		witResponse := getWitResponse([]byte(`
			{
				"entities": {
					"amount": [
						{ "value": "20 USD", "confidence": 100.00 }
					],
					"datetime": [
						{ "value": "", "confidence": 100.00 }
					],
					"description": []
				}
			}
		`))

		isTracking, err := witResponse.IsTracking()
		assert.True(t, isTracking)
		assert.EqualError(t, err, "No description")
	})

	t.Run("Returns tracking error", func(t *testing.T) {
		witResponse := getWitResponse([]byte(`
			{
				"entities": {
					"amount": [
						{ "value": "20 USD", "confidence": 100.00 }
					],
					"datetime": [
						{ "value": "", "confidence": 100.00 }
					],
					"description": []
				}
			}
		`))

		intent := witResponse.GetIntent()
		assert.Equal(t, intent, INTENT_TRACKING_ERROR)
	})

	t.Run("Returns tracking success", func(t *testing.T) {
		witResponse := getWitResponse([]byte(`
			{
				"entities": {
					"amount": [
						{ "value": "20 USD", "confidence": 100.00 }
					],
					"datetime": [
						{ "value": "", "confidence": 100.00 }
					],
					"description": [
						{ "value": "Food", "confidence": 100.00 }
					]
				}
			}
		`))

		intent := witResponse.GetIntent()
		assert.Equal(t, INTENT_TRACKING_SUCCESS, intent)
	})

	t.Run("Returns period success", func(t *testing.T) {
		witResponse := getWitResponse([]byte(`
			{
				"entities": {
					"amount": [],
					"datetime": [],
					"description": [],
					"total_spent": [
						{ "value": "monthly", "confidence": 100.00 }
					]
				}
			}
		`))

		intent := witResponse.GetIntent()
		assert.Equal(t, INTENT_PERIOD_TOTAL_SUCCESS, intent)
	})

	t.Run("Returns unknown intent", func(t *testing.T) {
		witResponse := getWitResponse([]byte(`
			{
				"entities": {
					"amount": [],
					"datetime": [],
					"description": [],
					"total_spent": []
				}
			}
		`))

		intent := witResponse.GetIntent()
		assert.Equal(t, INTENT_UNKNOWN, intent)
	})
}
