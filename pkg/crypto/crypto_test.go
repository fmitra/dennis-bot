package crypto

import (
	"crypto/rsa"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeys(t *testing.T) {
	t.Run("Should encode and decode a private/public key pair", func(t *testing.T) {
		InitializeGob()

		publicKey, privateKey := CreateKeyPair()

		publicKeyStr, _ := encodeKey(publicKey)
		privateKeyStr, _ := encodeKey(privateKey)

		privateKeyInterface, _ := decodeKey(privateKeyStr, &rsa.PrivateKey{})
		publicKeyInterface, _ := decodeKey(publicKeyStr, &rsa.PublicKey{})

		publicKeyDecoded := publicKeyInterface.(rsa.PublicKey)
		privateKeyDecoded := privateKeyInterface.(rsa.PrivateKey)

		assert.Equal(t, privateKey.D, privateKeyDecoded.D)
		assert.Equal(t, privateKey.Primes, privateKeyDecoded.Primes)
		assert.Equal(t, publicKey.N, publicKeyDecoded.N)
	})

	t.Run("Should encrypt and decrypt text with a password", func(t *testing.T) {
		password := "moscow-in-june"
		text := "This is pivate content"

		encrypted, _ := Encrypt(text, password)
		assert.NotEqual(t, encrypted, text)

		decrypted, _ := Decrypt(encrypted, password)
		assert.Equal(t, text, decrypted)

		decrypted, _ = Decrypt(encrypted, "invalid-password")
		assert.NotEqual(t, decrypted, text)
	})

	t.Run("Should create an encrypted encoded key pair", func(t *testing.T) {
		password := "moscow-in-june"
		publicKeyStr, encryptedPrivateKey, _ := CreateProtectedKeyPair(password)
		privateKeyStr, _ := Decrypt(encryptedPrivateKey, password)

		privateKeyInterface, _ := decodeKey(privateKeyStr, &rsa.PrivateKey{})
		publicKeyInterface, _ := decodeKey(publicKeyStr, &rsa.PublicKey{})

		assert.IsType(t, rsa.PrivateKey{}, privateKeyInterface)
		assert.IsType(t, rsa.PublicKey{}, publicKeyInterface)
	})

	t.Run("Should encrypt and decrypt with public/private keys", func(t *testing.T) {
		msg := "hello world i made an expense today and earlier"
		password := "moscow-in-june"
		publicKeyStr, privateKeyStr, _ := CreateProtectedKeyPair(password)

		privateKey, _ := ParsePrivateKey(privateKeyStr, password)
		publicKey, _ := ParsePublicKey(publicKeyStr)

		encryptedText, _ := AsymEncrypt(msg, publicKey)
		decryptedText, _ := AsymDecrypt(encryptedText, privateKey)

		assert.Equal(t, msg, decryptedText)
	})

	t.Run("Should parse protected and encoded text to private key", func(t *testing.T) {
		password := "moscow-in-june"
		_, encryptedPrivateKey, _ := CreateProtectedKeyPair(password)

		privateKey, _ := ParsePrivateKey(encryptedPrivateKey, password)
		assert.IsType(t, rsa.PrivateKey{}, privateKey)
		assert.NotNil(t, privateKey.D)
		assert.NotNil(t, privateKey.Primes)
	})

	t.Run("Should parse encoded text to private key", func(t *testing.T) {
		_, key := CreateKeyPair()
		privateKeyStr, _ := encodeKey(key)
		password := ""

		privateKey, _ := ParsePrivateKey(privateKeyStr, password)
		assert.IsType(t, rsa.PrivateKey{}, privateKey)
		assert.NotNil(t, privateKey.D)
		assert.NotNil(t, privateKey.Primes)
	})

	t.Run("Should parse encoded text to public key", func(t *testing.T) {
		publicKeyStr, _, _ := CreateProtectedKeyPair("")

		publicKey, _ := ParsePublicKey(publicKeyStr)
		assert.IsType(t, rsa.PublicKey{}, publicKey)
		assert.NotNil(t, publicKey.N)
	})
}
