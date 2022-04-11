package aggregator

import (
	"fmt"
	"math/rand"
	"price_aggregator/models"
	"sync"
	"testing"
	"time"
)

func TestBarUpdate(t *testing.T) {

	b := bar{}

	if err := b.update("WRONG_VALUE", 10); err == nil {
		t.Error("expected to have an error, got nil")
	}

	b = bar{val: 10, count: 10}
	if err := b.update("50.0", 10); err != nil {
		t.Errorf("got err: %v \n", err)
	}

	expect := "13.64"
	if fmt.Sprintf("%.2f", b.val) != expect {
		t.Errorf("expect: %v, got: %v", expect, b.val)
	}

	if b.count != 11 {
		t.Errorf("expect: 11, got: %v", b.count)

	}

}

func TestNewBars(t *testing.T) {
	t.Parallel()

	bs := newBars(100, 1)
	time.Sleep(6 * time.Second)
	bs.stopFiller()

	count := 0
	for i := 0; i < len(bs.values); i++ {
		if bs.values[i].ts != 0 {
			count++
		}
	}

	if count < 5 {
		t.Errorf("expected to have more than 10 values after 10 second of idle")
	}
}

func TestNewBarsOverflow(t *testing.T) {
	t.Parallel()
	bs := newBars(5, 1)
	time.Sleep(6 * time.Second)
	bs.stopFiller()

	count := 0
	gotGap := false
	for i := 0; i < len(bs.values); i++ {
		if bs.values[i].ts != 0 {
			count++
		}

		if i > 0 && bs.values[i-1].ts > bs.values[i].ts {
			gotGap = true
		}
	}

	if count != 5 {
		t.Errorf("expected to have 5 filled values, got: %v", bs.values)
	}

	if !gotGap {
		t.Error("expected to have 'gap' between values is queue")
	}
}

func TestBarsUpdater(t *testing.T) {
	t.Parallel()
	bs := newBars(20, 1)

	prices := make(chan models.TickerPrice)
	errs := make(chan error)
	go bs.updater(prices, errs)

	prices2 := make(chan models.TickerPrice)
	errs2 := make(chan error)
	go bs.updater(prices2, errs2)

	go func() {
		select {
		case err := <-errs:
			if err != nil {
				t.Errorf("got err: %v", err)
			}
		case err := <-errs2:
			if err != nil {
				t.Errorf("got err2: %v", err)
			}
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)
	spammer := func(p chan models.TickerPrice) {
		for i := 0; i < 1000; i++ {
			time.Sleep(time.Millisecond * time.Duration(rand.Int63n(10)))
			priceFloat := float64(rand.Int31n(10000)) / 10.0
			t := time.Now()
			if rand.Int31n(2) == 1 {
				t = t.Add(-1 * time.Second)
			}
			p <- models.TickerPrice{
				Ticker: models.BTCUSDTicker,
				Time:   t,
				Price:  fmt.Sprint(priceFloat),
			}
		}
		close(p)
		// t.Log("channel closed")
		wg.Done()
	}

	go spammer(prices)
	go spammer(prices2)

	wg.Wait()
	time.Sleep(3 * time.Second)
	bs.stopFiller()
	// t.Logf("pos: %v", bs.pos)
	// t.Logf("count: %v", bs.count)
	// t.Log(bs.values)
}
