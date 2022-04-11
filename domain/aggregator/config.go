package aggregator

import "time"

type Config struct {
	Capacity    int           `default:"120"`
	BarInterval time.Duration `default:"1m"`
}
