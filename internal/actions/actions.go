// Package actions are a collection of actions that the bot
// performs on behalf of the user, such as creating an account.
package actions

import (
	"crypto/rsa"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/fmitra/dennis-bot/config"
	"github.com/fmitra/dennis-bot/pkg/alphapoint"
	"github.com/fmitra/dennis-bot/pkg/expenses"
	"github.com/fmitra/dennis-bot/pkg/sessions"
	"github.com/fmitra/dennis-bot/pkg/users"
	"github.com/fmitra/dennis-bot/pkg/wit"
)

// Actions are taken by the bot in response to a user request during
// conversation. Typically a user contacts the bot to request some action
// to be performed, such as expense tracking.
type Actions struct {
	Db         *gorm.DB
	Cache      sessions.Session
	Config     config.AppConfig
	Alphapoint alphapoint.Alphapoint
}

// CreateNewExpense creates and saves a new Expense entry to the DB.
func (a *Actions) CreateNewExpense(wr wit.Response, userID uint, pk rsa.PublicKey) error {
	date := wr.GetDate()
	amount, fromCurrency, _ := wr.GetAmount()
	targetCurrency := "USD"
	description, _ := wr.GetDescription()

	historicalAmount := a.ConvertCurrency(fromCurrency, targetCurrency, amount)
	expense := &expenses.Expense{
		Date:        date,
		Description: description,
		Total:       strconv.FormatFloat(amount, 'f', -1, 64),
		Historical:  strconv.FormatFloat(historicalAmount, 'f', -1, 64),
		Currency:    fromCurrency,
		UserID:      userID,
	}
	expense.Encrypt(pk)
	manager := expenses.NewExpenseManager(a.Db)
	return manager.Save(expense)
}

// ConvertCurrency attempts to convert currency based on a cached rate.
// If the cache is empty, it queries alphapoint for an updated rate
func (a *Actions) ConvertCurrency(from, to string, amount float64) float64 {
	var conversion alphapoint.Conversion
	cacheKey := fmt.Sprintf("%s_%s", from, to)
	err := a.Cache.Get(cacheKey, &conversion)

	// We ping alphapoint for an updated conversation rate if is not
	// already stored in our cache.
	if err != nil {
		oneWeek := 604800
		convertedAmount, newConversion := a.Alphapoint.Convert(from, to, amount)
		a.Cache.Set(cacheKey, newConversion, oneWeek)
		return convertedAmount
	}

	convertedAmount := conversion.Rate * amount
	return convertedAmount
}

// GetExpenseTotal returns the sum of historical expense history over a period of time.
func (a *Actions) GetExpenseTotal(period string, userID uint, pk rsa.PrivateKey) (string, error) {
	expenseM := expenses.NewExpenseManager(a.Db)
	total, err := expenseM.TotalByPeriod(period, userID, pk)
	if err != nil {
		// Just log for now, we'll handle the error in the message
		log.Printf("actions: failed to query expenses %s", err)
	}

	fromCurrency := "USD"
	settingsM := users.NewSettingManager(a.Db)
	toCurrency := settingsM.GetCurrency(userID)

	convertedAmount := total
	if toCurrency != fromCurrency {
		convertedAmount = a.ConvertCurrency(fromCurrency, toCurrency, total)
	}
	strAmount := strconv.FormatFloat(convertedAmount, 'f', 2, 64)
	messageVar := fmt.Sprintf("%s %s", strAmount, toCurrency)
	return messageVar, err
}

// CreateNewUser saves a user to the DB.
func (a *Actions) CreateNewUser(userID uint, password string) error {
	user := &users.User{
		TelegramID: userID,
		Password:   password,
	}
	manager := users.NewUserManager(a.Db)
	return manager.Save(user)
}

// SetUserCurrency creates a settings entry with the user's requested currency.
func (a *Actions) SetUserCurrency(userID uint, currency string) error {
	manager := users.NewSettingManager(a.Db)
	return manager.UpdateCurrency(userID, currency)
}

// GetExpenseCSV creates a CSV of the user's expense history for a specific time
// period and returns the name of the file.
func (a *Actions) GetExpenseCSV(period string, userID uint, pk rsa.PrivateKey) (string, error) {
	tmpFileName := fmt.Sprintf("expenses_%s", strconv.Itoa(int(userID)))
	tmpFile, err := ioutil.TempFile(os.TempDir(), tmpFileName)
	if err != nil {
		log.Printf("actions: failed to create csv file %s", err)
	}

	fileName := fmt.Sprintf("%s.csv", tmpFile.Name())
	os.Rename(tmpFile.Name(), fileName)

	csvFile, err := os.OpenFile(fileName, os.O_RDWR, os.ModeTemporary)
	if err != nil {
		log.Printf("actions: failed to open csv file %s", err)
	}
	defer csvFile.Close()

	m := expenses.NewExpenseManager(a.Db)
	expenses, err := m.QueryByPeriod(period, userID)
	if err != nil {
		log.Printf("actions: failed to query expenses %s", err)
		return fileName, err
	}

	w := csv.NewWriter(csvFile)
	for _, expense := range expenses {
		expense.Decrypt(pk)
		row := []string{
			expense.Date.Format(time.UnixDate),
			expense.Description,
			expense.Total,
			expense.Historical,
			expense.Currency,
		}
		if err := w.Write(row); err != nil {
			log.Printf("actions: failed to write to csv %s", err)
			return fileName, err
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		log.Printf("actions: failed to write to csv %s", err)
		return fileName, err
	}

	return fileName, nil
}
