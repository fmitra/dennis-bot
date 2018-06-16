package conversation

import (
	"fmt"
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

func (a *Actions) CreateNewExpense(witResponse wit.WitResponse, userId uint) bool {
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
		Total:       amount,
		Historical:  historicalAmount,
		Currency:    fromCurrency,
		UserID:      userId,
	}
	manager := expenses.NewExpenseManager(a.Db)
	return manager.Save(expense)
}

func (a *Actions) GetExpenseTotal(witResponse wit.WitResponse, userId uint) (string, error) {
	manager := expenses.NewExpenseManager(a.Db)
	period, err := witResponse.GetSpendPeriod()
	total, err := manager.TotalByPeriod(period, userId)
	messageVar := strconv.FormatFloat(total, 'f', 2, 64)

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
