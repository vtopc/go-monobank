package monobank

import (
	"fmt"
	"math"
	"strconv"

	"github.com/rmg/iso4217"
)

type ClientInfo struct {
	Name     string    `json:"name"`
	Accounts []Account `json:"accounts"`
}

type Account struct {
	ID           string `json:"id"` // {account}
	CurrencyCode int    `json:"currencyCode"`
	CashbackType string `json:"cashbackType"`
	Balance      int64  `json:"balance"`
	CreditLimit  int64  `json:"creditLimit"`
}

func (c *ClientInfo) String() string {
	if c == nil {
		return "<nil>"
	}

	resp := "І'мя: " + c.Name + "\n"
	for _, a := range c.Accounts {
		resp = resp + "🔸🔸🔸\n" + a.String()
	}

	return resp
}

func (a *Account) String() string {
	var currency string
	var minorUnits int

	currency, minorUnits = iso4217.ByCode(a.CurrencyCode)
	if len(currency) == 0 {
		currency = strconv.Itoa(a.CurrencyCode)
		minorUnits = 2
	}

	rate := math.Pow(10, float64(minorUnits))

	return "Account: " + a.ID + "\n" +
		"Валюта: " + currency + "\n" +
		// "Кешбек: " + a.CashbackType + "\n" +
		fmt.Sprintf("Власні кошти: %.2f %s\n", toBanknote(a.Balance-a.CreditLimit, rate), currency) +
		fmt.Sprintf("Баланс: %.2f %s\n", toBanknote(a.Balance, rate), currency) +
		fmt.Sprintf("Кредитний ліміт: %.2f %s\n", toBanknote(a.CreditLimit, rate), currency)
}

func toBanknote(i int64, rate float64) float64 {
	return float64(i) / rate
}
