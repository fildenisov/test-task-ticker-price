package app

import (
	"time"

	"github.com/fildenisov/test-task-ticker-price/delivery/http"
	"github.com/fildenisov/test-task-ticker-price/domain/aggregator"
)

type Config struct {
	HTTP         http.Config
	Aggregator   aggregator.Config
	StartTimeout time.Duration
	StopTimeout  time.Duration
}
