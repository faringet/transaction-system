package http

type Request struct {
	CurrencyCode int     `json:"currency_code"`
	Amount       float64 `json:"amount"`
	WalletNumber int     `json:"wallet_number"`
	CardNumber   int     `json:"card_number"`
}
