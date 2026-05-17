package mpesa

import "github.com/shopspring/decimal"

type TransactionPayload struct {
	TransactionID string `json:"transaction_id"`
	Amount decimal.Decimal
	CurrencyCode string `json:"currency_code"`
	Reference string `json:"reference"`
	Timestamp string `json:"timestamp"`
}