package models

import "time"

// Ticker represents a pair of two currencies
type Ticker string

const (
	// BTCUSDTicker is a ticker for BTC-USD pair
	BTCUSDTicker Ticker = "BTC_USD"
)

// TickerPrice is an incomming ticker price from an exchange
type TickerPrice struct {
	Ticker Ticker
	Time   time.Time
	Price  string // decimal value. example: "0", "10", "12.2", "13.2345122"
}
