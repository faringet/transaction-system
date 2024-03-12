package domain

import "time"

type Transactions struct {
	ID         int
	ClientID   int
	Client     *Clients
	CurrencyID int
	Currency   *Currencies
	Amount     float64
	Status     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
