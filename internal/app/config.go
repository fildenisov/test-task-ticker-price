package app

import (
	"time"

	"price_aggregator/delivery/http"
)

type Config struct {
	HTTP         http.Config
	StartTimeout time.Duration
	StopTimeout  time.Duration
}
