package expenses

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fmitra/dennis-bot/pkg/crypto"
)

func TestExpenseModel(t *testing.T) {
	t.Run("Should encrypt and decrypt model fields", func(t *testing.T) {
		expense := Expense{
			Description: "Food",
			Total:       "100.00",
			Historical:  "1.58",
			Currency:    "RUB",
		}
		publicKey, privateKey, _ := crypto.CreateKeyPair()

		expense.Encrypt(publicKey)
		assert.NotEqual(t, expense.Total, "100.00")
		assert.NotEqual(t, expense.Historical, "1.58")
		assert.NotEqual(t, expense.Description, "Food")
		assert.NotEqual(t, expense.Currency, "RUB")

		expense.Decrypt(privateKey)
		assert.Equal(t, expense.Total, "100.00")
		assert.Equal(t, expense.Historical, "1.58")
		assert.Equal(t, expense.Description, "Food")
		assert.Equal(t, expense.Currency, "RUB")
	})
}
