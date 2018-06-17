package expenses

import (
	"crypto/rsa"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/fmitra/dennis-bot/pkg/crypto"
	"github.com/fmitra/dennis-bot/pkg/users"
)

// Describes a tracked expense
type Expense struct {
	gorm.Model
	Date        time.Time `gorm:"index;not null"` // Date the expense was made
	Description string    `gorm:"not null"`       // Description of the expense
	Total       string    `gorm:"not null"`       // Total amount paid for the expense
	Historical  string    // Historical USD value of the total
	Currency    string    `gorm:"not null"`         // Currency ISO of the total
	Category    string    `gorm:"type:varchar(30)"` // Category of the expense
	User        users.User
	UserID      uint
}

func (e *Expense) Encrypt(publicKey rsa.PublicKey) error {
	total, err := crypto.AsymEncrypt(e.Total, publicKey)
	historical, err := crypto.AsymEncrypt(e.Historical, publicKey)
	description, err := crypto.AsymEncrypt(e.Description, publicKey)
	currency, err := crypto.AsymEncrypt(e.Currency, publicKey)
	if err != nil {
		return err
	}

	e.Total = total
	e.Historical = historical
	e.Description = description
	e.Currency = currency

	return nil
}

func (e *Expense) Decrypt(privateKey rsa.PrivateKey) error {
	total, err := crypto.AsymDecrypt(e.Total, privateKey)
	historical, err := crypto.AsymDecrypt(e.Historical, privateKey)
	description, err := crypto.AsymDecrypt(e.Description, privateKey)
	currency, err := crypto.AsymDecrypt(e.Currency, privateKey)
	if err != nil {
		return err
	}

	e.Total = total
	e.Historical = historical
	e.Description = description
	e.Currency = currency

	return nil
}
