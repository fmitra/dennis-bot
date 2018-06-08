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
)

func InitializeGob() {
	gob.Register(rsa.PublicKey{})
	gob.Register(rsa.PrivateKey{})
}

// Encrypts text using an rsa public key. Text is returned as a base64 encoded
// string for convient storage
func AsymEncrypt(text string, publicKey rsa.PublicKey) (string, error) {
	reader := rand.Reader
	forEncryption := []byte(text)
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), reader, &publicKey, forEncryption, nil)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypts base64 encoded text using an rsa private key.
func AsymDecrypt(text string, privateKey rsa.PrivateKey) (string, error) {
	decodedText, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}

	reader := rand.Reader
	plaintext, err := rsa.DecryptOAEP(sha256.New(), reader, &privateKey, decodedText, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// Parses a base64 encoded string and type asserts the value
// into an rsa private key
func ParsePrivateKey(key string, password string) (rsa.PrivateKey, error) {
	var err error
	encodedKey := key
	if password != "" {
		encodedKey, err = Decrypt(key, password)
		if err != nil {
			return rsa.PrivateKey{}, err
		}
	}

	privateKeyInterface, err := decodeKey(encodedKey, &rsa.PrivateKey{})
	if err != nil {
		return rsa.PrivateKey{}, err
	}

	return privateKeyInterface.(rsa.PrivateKey), nil
}

// Parses a base64 encoded string and type asserts the value
// into an rsa public key
func ParsePublicKey(key string) (rsa.PublicKey, error) {
	publicKeyInterface, err := decodeKey(key, &rsa.PublicKey{})
	if err != nil {
		return rsa.PublicKey{}, err
	}

	return publicKeyInterface.(rsa.PublicKey), nil
}

// Creates a string encoded public/private key pair. Private keys
// are password protected to ensure we are unable to access them on
// behalf of the user.
func CreateProtectedKeyPair(password string) (string, string, error) {
	publicKey, privateKey := CreateKeyPair()

	publicKeyStr, err := encodeKey(publicKey)
	privateKeyStr, err := encodeKey(privateKey)
	encryptedPrivateKeyStr, err := Encrypt(privateKeyStr, password)
	if err != nil {
		return "", "", err
	}

	return publicKeyStr, encryptedPrivateKeyStr, nil
}

// Creates a public/private key pair.
func CreateKeyPair() (rsa.PublicKey, rsa.PrivateKey) {
	reader := rand.Reader
	bitSize := 1024

	key, err := rsa.GenerateKey(reader, bitSize)
	if err != nil {
		log.Panicf("users: unable to generate key pair %s", err)
	}

	publicKey := key.PublicKey
	privateKey := *key

	return publicKey, privateKey
}

// Encrypts and base64 encodes a string with a password
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

// Decrypts a base64 encoded string with a password
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

// Encodes a key type to a base64 encoded string for storage
func encodeKey(key interface{}) (string, error) {
	b := bytes.Buffer{}
	encoder := gob.NewEncoder(&b)
	err := encoder.Encode(&key)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b.Bytes()), nil
}

// Returns an interface type of the encoded key for type assertion
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
