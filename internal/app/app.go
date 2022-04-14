package app

import (
	"context"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	"github.com/fildenisov/test-task-ticker-price/consts"
	"github.com/fildenisov/test-task-ticker-price/delivery/http"
	"github.com/fildenisov/test-task-ticker-price/domain/aggregator"
	"github.com/fildenisov/test-task-ticker-price/internal/rep"
	"github.com/fildenisov/test-task-ticker-price/mocks/stream"
	"github.com/fildenisov/test-task-ticker-price/models"
)

type cmp struct {
	Service rep.Lifecycle
	Name    string
}

// App respesents the application.
// Import App only in cmd derectory.
type App struct {
	log  *zerolog.Logger
	cmps []cmp
	cfg  Config
}

// New is a constructor for App
func New(cfg Config) *App {
	l := zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().
		Str("cmp", "app").Logger()
	return &App{
		log:  &l,
		cfg:  cfg,
		cmps: []cmp{},
	}
}

// Start starts application
func (a *App) Start(ctx context.Context) error {
	a.log.Info().Msg("starting application")

	agg := aggregator.New(a.cfg.Aggregator)
	h := http.New(a.cfg.HTTP, agg)

	a.cmps = append(a.cmps, cmp{h, "http"}, cmp{agg, "aggregator"})

	// adding fake streams
	for t, count := range a.cfg.FakeStreams {

		for i := 0; i < count; i++ {
			fsCfg := stream.Config{
				Ticker:    t,
				PriceFrom: a.cfg.FakeStreamMinPrice,
				PriceTo:   a.cfg.FakeStreamMaxPrice,
				Period:    a.cfg.FakeStreamPeriod,
			}
			fs := stream.New(fsCfg, agg)
			a.cmps = append(a.cmps, cmp{fs, fmt.Sprintf("fake_steam_%v", i)})
		}
	}

	okCh, errCh := make(chan struct{}), make(chan error)
	go func() {
		for _, cmp := range a.cmps {
			a.log.Info().Msgf("%v is starting", cmp.Name)
			if err := cmp.Service.Start(ctx); err != nil {
				a.log.Error().Err(err).Msgf(consts.FmtCannotStart, cmp.Name)
				errCh <- errors.Wrapf(err, consts.FmtCannotStart, cmp.Name)
			}
		}

		okCh <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return models.ErrStartTimeout
	case err := <-errCh:
		return err
	case <-okCh:
		return nil
	}
}

// Stop stops application
func (a *App) Stop(ctx context.Context) error {
	a.log.Info().Msg("shutting down service...")

	errCh := make(chan error)
	go func() {
		gr, ctx := errgroup.WithContext(ctx)
		for _, cmp := range a.cmps {
			a.log.Info().Msgf("stopping %q...", cmp.Name)
			if err := cmp.Service.Stop(ctx); err != nil {
				a.log.Error().Err(err).Msgf("cannot stop %q", cmp.Name)
			}
		}
		errCh <- gr.Wait()
	}()

	select {
	case <-ctx.Done():
		return models.ErrShutdownTimeout
	case err := <-errCh:
		if err != nil {
			return err
		}
		return nil
	}
}
