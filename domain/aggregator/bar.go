package aggregator

import (
	"strconv"
	"sync"
	"time"

	"github.com/rs/zerolog"

	"github.com/fildenisov/test-task-ticker-price/models"
)

// bar stores agregated index value per time interval.
type bar struct {
	ts    int64   // unix timestamp
	val   float64 // aggregated index value
	count int     // incidates how many indexes were agregated in val
}

func (b *bar) update(tp string) error {
	val, err := strconv.ParseFloat(tp, 64)
	if err != nil {
		return err
	}

	// calculation average val
	b.val = (b.val*float64(b.count) + val) / float64(b.count+1)
	b.count++

	return nil
}

// bars stores circular queue of []bar.
// It provides concurrently safe add() method that cicles the queue if necessory.
type bars struct {
	log         *zerolog.Logger
	stopFill    chan struct{}
	ticker      models.Ticker
	values      []bar
	pos         int
	count       int
	intervalSec int
	sync.Mutex
}

func newBars(log *zerolog.Logger, t models.Ticker, cap, intervalSec int) *bars {
	bs := &bars{
		log:         log,
		ticker:      t,
		values:      make([]bar, cap),
		intervalSec: intervalSec,
		stopFill:    make(chan struct{}),
	}
	bs.startFiller()
	return bs

}

// add adds new empty bar into the queue.
func (bs *bars) add(ts int64) *bar {
	bs.values[bs.pos] = bar{
		ts: ts,
	}

	if bs.pos == len(bs.values)-1 {
		bs.pos = 0
	} else {
		bs.pos++
	}

	if bs.count != len(bs.values) {
		bs.count++
	}

	return &bs.values[bs.pos]
}

// get finds and return the bar, ts MUST be correct interval
func (bs *bars) get(ts int64) *bar {
	// get index of the newest bar
	pos := bs.pos - 1
	for i := 0; i < bs.count; i++ {
		if pos < 0 {
			pos = bs.count - 1
		}

		// if ts matched = we've found the bar
		if bs.values[pos].ts == ts {
			return &bs.values[pos]
		}

		// we will not be able to find the bar in the future
		if bs.values[pos].ts < ts {
			return nil
		}
		pos--
	}

	return nil
}

func (bs *bars) updater(prices <-chan models.TickerPrice, errs chan<- error) {
	for tp := range prices {
		ts := tp.Time.Unix()
		ts -= ts % int64(bs.intervalSec)
		bs.Lock()
		b := bs.get(ts)
		if b == nil {
			b = bs.add(ts)
		}

		err := b.update(tp.Price)
		bs.Unlock()

		if err != nil {
			bs.log.Error().Err(err).
				Stringer(models.KeyTicker, tp.Ticker).
				Int64("ts", ts).
				Str("price", tp.Price).
				Msg("bar update failed")
			errs <- err
		}
	}
}

// startFiller starts filling empty bars in case there are no new tickers in the channel
func (bs *bars) startFiller() {
	ticker := time.NewTicker(time.Duration(bs.intervalSec) * time.Second)
	go func() {
		bs.log.Info().Stringer(models.KeyTicker, bs.ticker).Msg("filler is started")
		for {
			select {
			case <-bs.stopFill: // can be useful later if will decide to stop filler
				return
			case t := <-ticker.C:
				ts := t.Unix()
				ts -= ts % int64(bs.intervalSec)
				bs.Lock()
				if b := bs.get(ts); b == nil {
					bs.add(ts)
				}
				bs.Unlock()
			}
		}
	}()
}

func (bs *bars) stopFiller() {
	bs.log.Info().Stringer(models.KeyTicker, bs.ticker).Msg("filler is stopped")
	bs.stopFill <- struct{}{}
}
