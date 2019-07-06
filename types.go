package monobank

import (
	"fmt"
	"strconv"
	"time"

	"github.com/rmg/iso4217"
)

type Time time.Time

// ClientInfo - client/user info
type ClientInfo struct {
	Name     string    `json:"name"`
	Accounts []Account `json:"accounts"`
}

type Account struct {
	ID           string `json:"id"` // account ID
	Balance      int64  `json:"balance"`
	CreditLimit  int64  `json:"creditLimit"`
	CurrencyCode int    `json:"currencyCode"`
	CashbackType string `json:"cashbackType"` // enum: None, UAH, Miles
}

// Statement - bank account statement
type Statement struct {
	ID              string `json:"id"`
	Time            Time   `json:"time"`
	Description     string `json:"description"`
	MCC             int32  `json:"mcc"`
	Hold            bool   `json:"hold"`
	Amount          int64  `json:"amount"`
	OperationAmount int64  `json:"operationAmount"`
	CurrencyCode    int    `json:"currencyCode"`
	CommissionRate  int64  `json:"commissionRate"`
	CashbackAmount  int64  `json:"cashbackAmount"`
	Balance         int64  `json:"balance"`
}

// Statements - transactions
type Statements []Statement

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
	currency, minorUnits := iso4217.ByCode(a.CurrencyCode)
	if len(currency) == 0 {
		// setting defaults:
		currency = strconv.Itoa(a.CurrencyCode)
		minorUnits = 2
	}

	return "Account: /ST_" + a.ID + "\n" +
		"Валюта: " + currency + "\n" +
		// "Кешбек: " + a.CashbackType + "\n" +
		// TODO: rewrite to string concat:
		fmt.Sprintf("Власні кошти: %s %s\n", ToBanknote(a.Balance-a.CreditLimit, minorUnits), currency) +
		fmt.Sprintf("Баланс: %s %s\n", ToBanknote(a.Balance, minorUnits), currency) +
		fmt.Sprintf("Кредитний ліміт: %s %s\n", ToBanknote(a.CreditLimit, minorUnits), currency)
}

func (t *Time) UnmarshalJSON(data []byte) error {
	ts, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}

	*t = Time(time.Unix(ts, 0))

	return nil
}
