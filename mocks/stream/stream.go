package stream

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/fildenisov/test-task-ticker-price/consts"
	"github.com/fildenisov/test-task-ticker-price/domain/aggregator"
	"github.com/fildenisov/test-task-ticker-price/models"
	"github.com/rs/zerolog"
)

// FakeStream emulates ticker index stream
type FakeStream struct {
	log  *zerolog.Logger
	agg  *aggregator.Aggregator
	done chan bool
	cfg  Config
}

// New is an FakeStream constuctor
func New(cfg Config, agg *aggregator.Aggregator) *FakeStream {
	l := zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Str("cmp", "fake_stream").Logger()
	return &FakeStream{
		log:  &l,
		cfg:  cfg,
		agg:  agg,
		done: make(chan bool),
	}
}

// Start starts aggregator component
func (f *FakeStream) Start(context.Context) error {

	tp, errs := f.agg.SubscribePriceStream(models.Ticker(f.cfg.Ticker))
	go f.worker(tp, errs)

	return nil
}

// Stop stops aggregator component
func (f *FakeStream) Stop(ctx context.Context) error {
	return nil
}

func (f *FakeStream) worker(tp chan models.TickerPrice, errs chan error) {

	f.log.Debug().Str(consts.KeyTicker, f.cfg.Ticker).Msg("fake stream worker started")
	errDone := make(chan bool)

	go func() {
		for {
			select {
			case <-errDone:
				close(errs)
				return
			case err := <-errs:
				f.log.Error().Err(err).Msg("got err in fake worker")
			}
		}
	}()

	t := time.NewTicker(f.cfg.Period)
	go func() {
		for {
			select {
			case <-f.done:
				errDone <- true
				close(tp)
				return
			case t := <-t.C:
				rPrice := float64(f.cfg.PriceFrom) + rand.Float64()*(float64(f.cfg.PriceTo)-float64(f.cfg.PriceFrom))

				tp <- models.TickerPrice{
					Ticker: models.Ticker(f.cfg.Ticker),
					Time:   t,
					Price:  fmt.Sprintf("%.2f", rPrice),
				}
			}
		}
	}()
}
