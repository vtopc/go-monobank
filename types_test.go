package monobank

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testAccount = Account{
	ID:           "-An7eZ",
	CurrencyCode: 980,
	CashbackType: "UAH",
	Balance:      1743827,
	CreditLimit:  1000000,
}

func TestAccount_String(t *testing.T) {
	tests := []struct {
		name    string
		account Account
		want    string
	}{
		{
			name:    "UAH",
			account: testAccount,
			want: `Account: -An7eZ
Валюта: UAH
Баланс(чистий): 7438.27 UAH
Баланс(з кредитним лімітом): 17438.27 UAH
Кредитний ліміт: 10000.00 UAH
`,
		},
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.account.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestClientInfo_String(t *testing.T) {
	tests := []struct {
		name string
		c    ClientInfo
		want string
	}{
		{
			c: ClientInfo{
				Name: "Will",
				Accounts: []Account{
					testAccount,
				},
			},
			want: `І'мя: Will
🔸🔸🔸
Account: -An7eZ
Валюта: UAH
Баланс(чистий): 7438.27 UAH
Баланс(з кредитним лімітом): 17438.27 UAH
Кредитний ліміт: 10000.00 UAH
`,
		},
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.c.String()
			assert.Equal(t, tt.want, got)
		})
	}
}
