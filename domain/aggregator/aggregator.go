package aggregator

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"

	"github.com/fildenisov/test-task-ticker-price/models"
)

// Aggregator stores all aggregated values for different tickers.
// Aggregator is concurrency safe.
type Aggregator struct {
	log         *zerolog.Logger
	tickers     map[models.Ticker]*bars
	barInverval time.Duration
	sync.RWMutex
	capPerTicker int
}

// New is an Aggregator constuctor
func New(cfg Config) *Aggregator {
	l := zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().
		Str("cmp", "http").Logger()
	return &Aggregator{
		log:          &l,
		tickers:      make(map[models.Ticker]*bars),
		barInverval:  cfg.BarInterval,
		capPerTicker: cfg.Capacity,
	}
}

// Start starts aggregator component
func (a *Aggregator) Start(context.Context) error { return nil }

// Stop stops aggregator component
func (a *Aggregator) Stop(ctx context.Context) error {
	for _, bs := range a.tickers {
		bs.stopFiller()
	}

	return nil
}

// SubscribePriceStream subscribes aggregator to a new price source
func (a *Aggregator) SubscribePriceStream(t models.Ticker) (chan models.TickerPrice, chan error) {
	prices := make(chan models.TickerPrice)
	errs := make(chan error)
	bs := a.getBars(t)

	go bs.updater(prices, errs)

	return prices, errs
}

// GetBars return last 'max' known bars
func (a *Aggregator) GetBars(t models.Ticker, max int) ([]models.Bar, bool) {
	if max <= 0 {
		return []models.Bar{}, false
	}

	// acquire read lock
	a.RLock()
	defer a.RUnlock()

	// search if ticker is known
	bs, ok := a.tickers[t]
	if !ok {
		a.log.Debug().Stringer(models.KeyTicker, t).Msg("ticker not found")
		return []models.Bar{}, false
	}

	// max is limited by current queue length
	if max > bs.count {
		max = bs.count
	}

	// the newest price is locked before current write position
	newestIndex := bs.pos - 1

	// result capacity is equal to the current count
	res := make([]models.Bar, 0, max)

	bs.Lock()
	defer bs.Unlock()
	for i := 0; i < max; i++ {
		if newestIndex < 0 {
			// in case we are going out of range
			newestIndex = bs.count - 1
		}

		res = append(res, models.Bar{
			TS:    bs.values[newestIndex].ts,
			Price: bs.values[newestIndex].val,
		})
		newestIndex--
	}

	return res, true
}

// getBars gets or created ticker bars
func (a *Aggregator) getBars(t models.Ticker) *bars {
	a.Lock()
	defer a.Unlock()
	bs, ok := a.tickers[t]

	if !ok {
		a.log.Debug().Stringer(models.KeyTicker, t).Msg("creating ticker bars")
		bs = newBars(a.log, t, a.capPerTicker, int(a.barInverval.Seconds()))
		a.tickers[t] = bs
	}
	return bs
}
