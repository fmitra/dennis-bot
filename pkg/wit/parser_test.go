package wit

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/kierdavis/dateparser"
	"github.com/stretchr/testify/assert"
)

func getResponse(b []byte) Response {
	var response Response
	json.Unmarshal(b, &response)
	return response
}

func TestWitParser(t *testing.T) {
	t.Run("Returns unknown intent for empty response", func(t *testing.T) {
		response := &Response{}
		assert.Equal(t, "unknown_request", response.GetMessageOverview())
	})

	t.Run("Returns amount", func(t *testing.T) {
		response := getResponse([]byte(`
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

		amount, currency, _ := response.GetAmount()
		assert.Equal(t, 20.0, amount)
		assert.Equal(t, "USD", currency)
	})

	t.Run("Returns error if amount is empty", func(t *testing.T) {
		response := getResponse([]byte(`
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

		amount, currency, err := response.GetAmount()
		assert.Equal(t, 0.0, amount)
		assert.Equal(t, "", currency)
		assert.EqualError(t, err, "no amount")
	})

	t.Run("Returns error if amount is blank", func(t *testing.T) {
		response := getResponse([]byte(`
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

		amount, currency, err := response.GetAmount()
		assert.Equal(t, 0.0, amount)
		assert.Equal(t, "", currency)
		assert.EqualError(t, err, "invalid amount")
	})

	t.Run("Returns error if currency is blank", func(t *testing.T) {
		response := getResponse([]byte(`
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

		amount, currency, err := response.GetAmount()
		assert.Equal(t, 0.0, amount)
		assert.Equal(t, "", currency)
		assert.EqualError(t, err, "invalid amount")
	})

	t.Run("Returns description", func(t *testing.T) {
		response := getResponse([]byte(`
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

		description, _ := response.GetDescription()
		assert.Equal(t, "Food", description)
	})

	t.Run("Returns error if description is blank", func(t *testing.T) {
		response := getResponse([]byte(`
			{
				"entities": {
					"amount": [],
					"datetime": [],
					"description": []
				}
			}
		`))

		description, err := response.GetDescription()
		assert.Equal(t, "", description)
		assert.EqualError(t, err, "no description")
	})

	t.Run("Returns date", func(t *testing.T) {
		parser := &dateparser.Parser{}
		expectedDate, _ := parser.Parse("2018/10/10")
		response := getResponse([]byte(`
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

		date := response.GetDate()
		assert.Equal(t, expectedDate, date)
	})

	t.Run("Returns default date if blank", func(t *testing.T) {
		response := getResponse([]byte(`
			{
				"entities": {
					"amount": [],
					"datetime": [],
					"description": []
				}
			}
		`))

		date := response.GetDate()
		assert.IsType(t, time.Time{}, date)
	})

	t.Run("Infers tracking", func(t *testing.T) {
		response := getResponse([]byte(`
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

		isTracking, err := response.IsTracking()
		assert.True(t, isTracking)
		assert.NoError(t, err)
	})

	t.Run("Fails to infer tracking if missing amount", func(t *testing.T) {
		response := getResponse([]byte(`
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

		isTracking, err := response.IsTracking()
		assert.False(t, isTracking)
		assert.EqualError(t, err, "no amount")
	})

	t.Run("Fails to infer tracking if missing description", func(t *testing.T) {
		response := getResponse([]byte(`
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

		isTracking, err := response.IsTracking()
		assert.True(t, isTracking)
		assert.EqualError(t, err, "no description")
	})

	t.Run("Returns tracking error", func(t *testing.T) {
		response := getResponse([]byte(`
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

		overview := response.GetMessageOverview()
		assert.Equal(t, TrackingRequestedError, overview)
	})

	t.Run("Returns tracking success", func(t *testing.T) {
		response := getResponse([]byte(`
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

		overview := response.GetMessageOverview()
		assert.Equal(t, TrackingRequestedSuccess, overview)
	})

	t.Run("Returns period success", func(t *testing.T) {
		response := getResponse([]byte(`
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

		overview := response.GetMessageOverview()
		assert.Equal(t, ExpenseTotalRequestedSuccess, overview)
	})

	t.Run("Returns unknown intent", func(t *testing.T) {
		response := getResponse([]byte(`
			{
				"entities": {
					"amount": [],
					"datetime": [],
					"description": [],
					"total_spent": []
				}
			}
		`))

		overview := response.GetMessageOverview()
		assert.Equal(t, UnknownRequest, overview)
	})
}
