package aggregator

import (
	"context"
	"price_aggregator/models"
	"sync"
	"time"
)

// Aggregator stores all aggregated values for different tickers.
// Aggregator is concurrency safe.
type Aggregator struct {
	sync.RWMutex
	tickers      map[models.Ticker]*bars
	capPerTicker int
	barInverval  time.Duration
}

// New is an Aggregator constuctor
func New(ctx context.Context, cfg Config) *Aggregator {
	return &Aggregator{
		tickers:      make(map[models.Ticker]*bars),
		barInverval:  cfg.BarInterval,
		capPerTicker: cfg.Capacity,
	}
}

func (a *Aggregator) Start(context.Context) error { return nil }
func (a *Aggregator) Stop(ctx context.Context) error {
	for _, bs := range a.tickers {
		bs.stopFiller()
	}

	return nil
}

// GetBars return last 'max' known bars
func (a *Aggregator) GetBars(ticker models.Ticker, max int) ([]models.Bar, bool) {
	if max <= 0 {
		return []models.Bar{}, false
	}

	// acquire read lock
	a.RLock()
	defer a.RUnlock()

	// search if ticker is known
	bs, ok := a.tickers[ticker]
	if !ok {
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
	a.RLock()
	bs, ok := a.tickers[t]
	a.RUnlock()

	if !ok {
		a.Lock()

		bs = newBars(a.capPerTicker, int(a.barInverval.Seconds()))
		a.tickers[t] = bs

		a.Unlock()
	}
	return bs
}

// SubscribePriceStream subscribes aggregator to a new price source
func (a *Aggregator) SubscribePriceStream(t models.Ticker) (chan models.TickerPrice, chan error) {
	prices := make(chan models.TickerPrice)
	errs := make(chan error)
	bs := a.getBars(t)

	go bs.updater(prices, errs)

	return prices, errs
}
