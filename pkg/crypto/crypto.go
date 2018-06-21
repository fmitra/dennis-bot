// Package crypto is a collection of encryption utility functions.
package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"io"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// InitializeGob registeres rsa PublicKey/PrivateKey so we may encode
// using stdlib encoding/gob.
func InitializeGob() {
	gob.Register(rsa.PublicKey{})
	gob.Register(rsa.PrivateKey{})
}

// AsymEncrypt encryptst ext using rsa.PublicKey. Text is returned as a base64 encoded
// string for convenient storage.
func AsymEncrypt(t string, pk rsa.PublicKey) (string, error) {
	r := rand.Reader
	b := []byte(t)
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), r, &pk, b, nil)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// AsymDecrypt decrypts base64  encoded text using an rsa.PrivateKey.
func AsymDecrypt(b64 string, pk rsa.PrivateKey) (string, error) {
	text, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", err
	}

	r := rand.Reader
	plaintext, err := rsa.DecryptOAEP(sha256.New(), r, &pk, text, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// ParsePrivateKey parses a base64 encoded string and type asserts the value
// into an rsa.PrivateKey
func ParsePrivateKey(b64 string, password string) (rsa.PrivateKey, error) {
	var err error
	isEncrypted := password != ""
	encodedKey := b64

	if isEncrypted {
		encodedKey, err = Decrypt(b64, password)
	}

	if err != nil {
		return rsa.PrivateKey{}, err
	}

	privateKeyInterface, err := decodeKey(encodedKey, &rsa.PrivateKey{})
	if err != nil {
		return rsa.PrivateKey{}, err
	}

	return privateKeyInterface.(rsa.PrivateKey), nil
}

// ParsePublicKey parsers a base64 encoded string and type asserts the value into
// an rsa.PublicKey
func ParsePublicKey(b64 string) (rsa.PublicKey, error) {
	publicKeyInterface, err := decodeKey(b64, &rsa.PublicKey{})
	if err != nil {
		return rsa.PublicKey{}, err
	}

	return publicKeyInterface.(rsa.PublicKey), nil
}

// CreateProtectedKeyPair creates a string encoded public/private key pair. Private
// keys are password protected to ensure we are unable to access them on
// behalf of the user.
func CreateProtectedKeyPair(password string) (string, string, error) {
	blankKey := ""
	publicKey, privateKey, err := CreateKeyPair()
	if err != nil {
		return blankKey, blankKey, err
	}

	publicKeyStr, err := encodeKey(publicKey)
	if err != nil {
		return blankKey, blankKey, err
	}

	privateKeyStr, err := encodeKey(privateKey)
	if err != nil {
		return blankKey, blankKey, err
	}

	encryptedPrivateKeyStr, err := Encrypt(privateKeyStr, password)
	if err != nil {
		return blankKey, blankKey, err
	}

	return publicKeyStr, encryptedPrivateKeyStr, nil
}

// CreateKeyPair creates an rsa.PublicKey and rsa.PrivateKey pair.
func CreateKeyPair() (rsa.PublicKey, rsa.PrivateKey, error) {
	var publicKey rsa.PublicKey
	var privateKey rsa.PrivateKey

	reader := rand.Reader
	bitSize := 1024

	key, err := rsa.GenerateKey(reader, bitSize)
	if err != nil {
		log.Printf("users: unable to generate key pair %s", err)
		return publicKey, privateKey, err
	}

	publicKey = key.PublicKey
	privateKey = *key

	return publicKey, privateKey, nil
}

// Encrypt encrypts and base64 encodes a string with a password.
func Encrypt(text string, password string) (string, error) {
	key := sha256.New()
	key.Write([]byte(password))

	block, err := aes.NewCipher(key.Sum(nil))
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(text))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(text))
	encodedText := base64.StdEncoding.EncodeToString(ciphertext)

	return encodedText, nil
}

// Decrypt decrypts a base64 encoded string with a password.
func Decrypt(text string, password string) (string, error) {
	key := sha256.New()
	key.Write([]byte(password))

	block, err := aes.NewCipher(key.Sum(nil))
	if err != nil {
		return "", err
	}

	if len(text) < aes.BlockSize {
		return "", errors.New("cipher text too short")
	}

	decodedText, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}

	iv := decodedText[:aes.BlockSize]
	decodedText = decodedText[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(decodedText, decodedText)

	return string(decodedText), nil
}

// HashText hashes a string text using bcrypt.
func HashText(text string) (string, error) {
	t := []byte(text)
	cost := 10
	hash, err := bcrypt.GenerateFromPassword(t, cost)
	if err != nil {
		log.Printf("users: failed to hash password")
		return "", err
	}

	return string(hash), nil
}

// ValidateHash validates a a string text against a bcrypt hash.
func ValidateHash(hash, text string) error {
	bHash := []byte(hash)
	bText := []byte(text)
	return bcrypt.CompareHashAndPassword(bHash, bText)
}

// encodeKey encodes a key type to a base64 encoded string for storage.
func encodeKey(key interface{}) (string, error) {
	b := bytes.Buffer{}
	encoder := gob.NewEncoder(&b)
	err := encoder.Encode(&key)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b.Bytes()), nil
}

// decodeKey returns an interface type of the encoded key for type assertion
// from a base65 encoded string.
func decodeKey(key string, keyType interface{}) (interface{}, error) {
	by, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return keyType, err
	}

	b := bytes.Buffer{}
	b.Write(by)

	decoder := gob.NewDecoder(&b)
	err = decoder.Decode(&keyType)
	if err != nil {
		return keyType, err
	}

	return keyType, nil
}
