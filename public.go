package monobank

import (
	"context"
	"fmt"
	"net/http"
)

type PublicAPI interface {
	// Currency https://api.monobank.ua/docs/#/definitions/CurrencyInfo
	Currency(context.Context) (Currencies, error)
}

func (c Client) Currency(ctx context.Context) (Currencies, error) {
	const urlPath = "/bank/currency"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlPath, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var v Currencies
	err = c.do(req, &v, http.StatusOK)

	return v, err
}
