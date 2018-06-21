package conversation

import (
	"crypto/rsa"
	"fmt"
	"log"
	"strconv"

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
	Db     *gorm.DB
	Cache  sessions.Session
	Config config.AppConfig
}

// CreateNewExpense creates and saves a new Expense entry to the DB.
func (a *Actions) CreateNewExpense(wr wit.Response, userID uint, pk rsa.PublicKey) error {
	date := wr.GetDate()
	amount, fromCurrency, _ := wr.GetAmount()
	targetCurrency := "USD"
	description, _ := wr.GetDescription()

	var conversion alphapoint.Conversion
	var newConversion *alphapoint.Conversion
	var historicalAmount float64
	cacheKey := fmt.Sprintf("%s_%s", fromCurrency, targetCurrency)
	err := a.Cache.Get(cacheKey, &conversion)
	historicalAmount = conversion.Rate * amount

	// We ping alphapoint for an updated conversation rate if is not
	// already stored in our cache.
	if err != nil {
		ap := alphapoint.NewClient(a.Config.AlphaPoint.Token)
		historicalAmount, newConversion = ap.Convert(
			fromCurrency,
			"USD",
			amount,
		)
		oneWeek := 604800
		a.Cache.Set(cacheKey, newConversion, oneWeek)
	}

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

// GetExpenseTotal returns the sum of historical expense history over a period of time.
func (a *Actions) GetExpenseTotal(period string, userID uint, pk rsa.PrivateKey) (string, error) {
	manager := expenses.NewExpenseManager(a.Db)
	total, err := manager.TotalByPeriod(period, userID, pk)
	if err != nil {
		log.Printf("actions: failed to query expenses %s", err)
	}

	messageVar := strconv.FormatFloat(total, 'f', 2, 64)
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
