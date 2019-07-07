package monobank

import (
	"strconv"
	"time"
)

// ClientInfo - client/user info
type ClientInfo struct {
	Name     string   `json:"name"`
	Accounts Accounts `json:"accounts"`
}

type Account struct {
	ID           string `json:"id"` // account ID
	Balance      int64  `json:"balance"`
	CreditLimit  int64  `json:"creditLimit"`
	CurrencyCode int    `json:"currencyCode"`
	CashbackType string `json:"cashbackType"` // enum: None, UAH, Miles
}

type Accounts []Account

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

type Currency struct {
	CurrencyCodeA int     `json:"currencyCodeA"`
	CurrencyCodeB int     `json:"currencyCodeB"`
	Date          Time    `json:"date"`
	RateSell      float64 `json:"rateSell"`
	RateBuy       float64 `json:"rateBuy"`
	RateCross     float64 `json:"rateCross"`
}

type Currencies []Currency

type Time time.Time

func (t *Time) UnmarshalJSON(data []byte) error {
	ts, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}

	*t = Time(time.Unix(ts, 0))

	return nil
}
