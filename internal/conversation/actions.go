package conversation

import (
	"crypto/rsa"
	"fmt"
	"log"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

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

func (a *Actions) CreateNewExpense(witResponse wit.WitResponse, userId uint, publicKey rsa.PublicKey) bool {
	date := witResponse.GetDate()
	amount, fromCurrency, _ := witResponse.GetAmount()
	targetCurrency := "USD"
	description, _ := witResponse.GetDescription()

	var conversion alphapoint.Conversion
	var newConversion *alphapoint.Conversion
	cacheKey := fmt.Sprintf("%s_%s", fromCurrency, targetCurrency)
	a.Cache.Get(cacheKey, &conversion)

	historicalAmount := conversion.Rate * amount
	if conversion.Rate == 0 {
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
		UserID:      userId,
	}
	expense.Encrypt(publicKey)
	manager := expenses.NewExpenseManager(a.Db)
	return manager.Save(expense)
}

func (a *Actions) GetExpenseTotal(expensePeriod string, userId uint, privateKey rsa.PrivateKey) (string, error) {
	manager := expenses.NewExpenseManager(a.Db)
	total, err := manager.TotalByPeriod(expensePeriod, userId, privateKey)
	messageVar := strconv.FormatFloat(total, 'f', 2, 64)
	if err != nil {
		log.Printf("actions: failed to query expenses %s", err)
	}

	return messageVar, err
}

func (a *Actions) CreateNewUser(userId uint, password string) bool {
	user := &users.User{
		TelegramID: userId,
		Password:   password,
	}
	manager := users.NewUserManager(a.Db)
	return manager.Save(user)
}
