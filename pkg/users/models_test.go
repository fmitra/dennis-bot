package users

import (
	"crypto/rsa"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fmitra/dennis-bot/pkg/crypto"
	mocks "github.com/fmitra/dennis-bot/test"
)

func TestUserModel(t *testing.T) {
	t.Run("It should return private key", func(t *testing.T) {
		password, _ := crypto.HashText("my-password")
		publicKey, privateKey, _ := crypto.CreateProtectedKeyPair(password)
		user := &User{
			TelegramID: uint(123),
			Password:   password,
			PublicKey:  publicKey,
			PrivateKey: privateKey,
		}

		userPrivateKey, _ := user.GetPrivateKey("my-password")
		assert.IsType(t, rsa.PrivateKey{}, userPrivateKey)
	})

	t.Run("It shoudl validate users password", func(t *testing.T) {
		hashedPassword, _ := crypto.HashText("my-password")
		user := &User{
			Password:   hashedPassword,
			TelegramID: mocks.TestUserID,
		}

		errMsg := "crypto/bcrypt: hashedPassword is not the hash of the given password"
		assert.NoError(t, user.ValidatePassword("my-password"))
		assert.EqualError(t, user.ValidatePassword("not-my-password"), errMsg)
	})
}
