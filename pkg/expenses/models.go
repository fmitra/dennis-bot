// Package expenses represents an Expense a User logged with the Bot service.
package expenses

import (
	"crypto/rsa"
	"log"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/fmitra/dennis-bot/pkg/crypto"
	"github.com/fmitra/dennis-bot/pkg/users"
)

// Expense is a tracked expnese entry.
type Expense struct {
	gorm.Model
	Date        time.Time `gorm:"index;not null"` // Date the expense was made
	Description string    `gorm:"not null"`       // Description of the expense
	Total       string    `gorm:"not null"`       // Total amount paid for the expense
	Historical  string    // Historical USD value of the total
	Currency    string    `gorm:"not null"`         // Currency ISO of the total
	Category    string    `gorm:"type:varchar(30)"` // Category of the expense
	User        users.User
	UserID      uint `gorm:"index;not null"`
}

// Encrypt encrypts sensitive fields in an Expense record.
func (e *Expense) Encrypt(publicKey rsa.PublicKey) error {
	total, err := crypto.AsymEncrypt(e.Total, publicKey)
	if err != nil {
		log.Printf("expenses: failed to encrypt total - %s", err)
		return err
	}

	historical, err := crypto.AsymEncrypt(e.Historical, publicKey)
	if err != nil {
		log.Printf("expenses: failed to encrypt historical - %s", err)
		return err
	}

	description, err := crypto.AsymEncrypt(e.Description, publicKey)
	if err != nil {
		log.Printf("expenses: failed to encrypt description - %s", err)
		return err
	}

	currency, err := crypto.AsymEncrypt(e.Currency, publicKey)
	if err != nil {
		log.Printf("expenses: failed to encrypt currency - %s", err)
		return err
	}

	e.Total = total
	e.Historical = historical
	e.Description = description
	e.Currency = currency

	return nil
}

// Decrypt decrypts an Expense record.
func (e *Expense) Decrypt(privateKey rsa.PrivateKey) error {
	total, err := crypto.AsymDecrypt(e.Total, privateKey)
	if err != nil {
		log.Printf("expenses: failed to decrypt total - %s", err)
		return err
	}

	historical, err := crypto.AsymDecrypt(e.Historical, privateKey)
	if err != nil {
		log.Printf("expenses: failed to decrypt historical - %s", err)
		return err
	}

	description, err := crypto.AsymDecrypt(e.Description, privateKey)
	if err != nil {
		log.Printf("expenses: failed to decrypt description - %s", err)
		return err
	}

	currency, err := crypto.AsymDecrypt(e.Currency, privateKey)
	if err != nil {
		log.Printf("expenses: failed to decrypt currency - %s", err)
		return err
	}

	e.Total = total
	e.Historical = historical
	e.Description = description
	e.Currency = currency

	return nil
}
