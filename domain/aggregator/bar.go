package aggregator

import (
	"price_aggregator/models"
	"strconv"
	"sync"
	"time"
)

// bar stores agregated index value per time interval.
type bar struct {
	// sync.Mutex         // TODO: REMOVE
	ts    int64   // unix timestamp
	val   float64 // aggregated index value
	count int     // incidates how many indexes were agregated in val
}

func (b *bar) update(tp string, intervalSec int) error {
	val, err := strconv.ParseFloat(tp, 64)
	if err != nil {
		return err
	}

	// b.Lock()
	// defer b.Unlock()
	// calculation average val
	b.val = (b.val*float64(b.count) + val) / float64(b.count+1)
	b.count++

	return nil
}

// bars stores circular queue of []bar.
// It provides concurrently safe add() method that cicles the queue if necessory.
type bars struct {
	sync.Mutex
	values      []bar         // circular queue for all aggregated bars
	stopFill    chan struct{} // send anything to stop filler
	pos         int           // stores current index for next write
	count       int           // stores total for all filled bars. max(count) = len(values)
	intervalSec int           // stores interval duration in seconds
}

func newBars(cap, intervalSec int) *bars {
	bs := &bars{
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
		ts = ts - ts%int64(bs.intervalSec)
		bs.Lock()
		b := bs.get(ts)
		if b == nil {
			b = bs.add(ts)
		}

		err := b.update(tp.Price, bs.intervalSec)
		bs.Unlock()

		if err != nil {
			errs <- err
		}
	}
}

// startFiller starts filling empty bars in case there are no new tickers in the channel
func (bs *bars) startFiller() {
	ticker := time.NewTicker(time.Duration(bs.intervalSec) * time.Second)
	go func() {
		for {
			select {
			case <-bs.stopFill: // can be useful later if will decide to stop filler
				return
			case t := <-ticker.C:
				ts := t.Unix()
				ts = ts - ts%int64(bs.intervalSec)
				bs.Lock()
				b := bs.get(ts)
				if b == nil {
					b = bs.add(ts)
				}
				bs.Unlock()
			}
		}
	}()
}

func (bs *bars) stopFiller() {
	bs.stopFill <- struct{}{}
}
