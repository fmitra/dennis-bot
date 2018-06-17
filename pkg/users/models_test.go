package users

import (
	"crypto/rsa"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fmitra/dennis-bot/pkg/crypto"
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
}
