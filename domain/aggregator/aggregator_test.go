package aggregator

import (
	"context"
	"fmt"
	"math/rand"
	"price_aggregator/models"
	"sync"
	"testing"
	"time"
)

func TestAggregator(t *testing.T) {
	// t.Skip()
	t.Parallel()

	cfg := Config{
		Capacity:    15,
		BarInterval: time.Second,
	}

	ctx := context.Background()
	a := New(ctx, cfg)
	if err := a.Start(ctx); err != nil {
		t.Error("got err on Start")
	}

	tp, errs := a.SubscribePriceStream(models.BTCUSDTicker)
	tp2, errs2 := a.SubscribePriceStream(models.BTCUSDTicker)

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
		for i := 0; i < 3000; i++ {
			time.Sleep(time.Millisecond * time.Duration(rand.Int63n(10)))
			priceFloat := float64(rand.Int31n(10000)) / 10.0
			p <- models.TickerPrice{
				Ticker: models.BTCUSDTicker,
				Time:   time.Now(),
				Price:  fmt.Sprint(priceFloat),
			}
		}
		close(p)
		// t.Log("channel closed")
		wg.Done()
	}

	go spammer(tp)
	go spammer(tp2)

	wg.Wait()

	bars, ok := a.GetBars(models.BTCUSDTicker, -1)
	if len(bars) != 0 || ok {
		t.Errorf("expect len(bars)=0, got: %v", len(bars))
	}

	bars, ok = a.GetBars(models.Ticker("BTC_UNKNOWN"), 100)
	if len(bars) != 0 || ok {
		t.Errorf("expect len(bars)=0, got: %v", len(bars))
	}

	bars, ok = a.GetBars(models.BTCUSDTicker, 2000)
	if len(bars) != 15 || !ok {
		t.Errorf("expect len(bars)=15, got: %v", len(bars))
	}

	if err := a.Stop(ctx); err != nil {
		t.Error("got err on Stop")
	}
}
