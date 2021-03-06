package monobank

import "github.com/vtopc/epoch"

// ClientInfo - client/user info
type ClientInfo struct {
	Name       string   `json:"name"`
	WebHookURL string   `json:"webHookUrl"`
	Accounts   Accounts `json:"accounts"`
}

type Account struct {
	AccountID    string   `json:"id"`
	Balance      int64    `json:"balance"`
	CreditLimit  int64    `json:"creditLimit"`
	CurrencyCode int      `json:"currencyCode"`
	CashbackType string   `json:"cashbackType"` // enum: None, UAH, Miles
	CardMasks    []string `json:"maskedPan"`    // card number masks
	Type         CardType `json:"type"`
	IBAN         string   `json:"iban"`
}

type CardType string

const (
	Black    CardType = "black"
	White    CardType = "white"
	Platinum CardType = "platinum"
	FOP      CardType = "fop"
)

type Accounts []Account

// Transaction - bank account statement
type Transaction struct {
	ID              string        `json:"id"`
	Time            epoch.Seconds `json:"time"`
	Description     string        `json:"description"`
	MCC             int32         `json:"mcc"`
	OriginalMCC     int32         `json:"originalMcc"`
	Hold            bool          `json:"hold"`
	Amount          int64         `json:"amount"`
	OperationAmount int64         `json:"operationAmount"`
	CurrencyCode    int           `json:"currencyCode"`
	CommissionRate  int64         `json:"commissionRate"`
	CashbackAmount  int64         `json:"cashbackAmount"`
	Balance         int64         `json:"balance"`
	Comment         string        `json:"comment"`
	ReceiptID       string        `json:"receiptId"`
	EDRPOU          string        `json:"counterEdrpou"`
	IBAN            string        `json:"counterIban"`
}

// Transactions - transactions
type Transactions []Transaction

type Currency struct {
	CurrencyCodeA int           `json:"currencyCodeA"`
	CurrencyCodeB int           `json:"currencyCodeB"`
	Date          epoch.Seconds `json:"date"`
	RateSell      float64       `json:"rateSell"`
	RateBuy       float64       `json:"rateBuy"`
	RateCross     float64       `json:"rateCross"`
}

type Currencies []Currency

type WebHookRequest struct {
	WebHookURL string `json:"webHookUrl"`
}

type WebHookResponse struct {
	Type string      `json:"type"` // "StatementItem"
	Data WebHookData `json:"data"`
}

type WebHookData struct {
	AccountID   string      `json:"account"`
	Transaction Transaction `json:"statementItem"`
}

type TokenRequest struct {
	RequestID string `json:"tokenRequestId"` // Unique token request ID.
	AcceptURL string `json:"acceptUrl"`      // URL to redirect client or build QR on top of it.
}
