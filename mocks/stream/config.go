package stream

import "time"

type Config struct {
	Ticker    string
	PriceFrom int
	PriceTo   int
	Period    time.Duration
}
