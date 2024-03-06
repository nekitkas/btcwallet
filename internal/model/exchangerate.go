package model

import "time"

type ExchangeRate struct {
	Symbol    string    `json:"symbol"`
	Value     string    `json:"value"`
	Sources   int       `json:"sources"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ExchangeRates struct {
	Data []ExchangeRate `json:"data"`
}
