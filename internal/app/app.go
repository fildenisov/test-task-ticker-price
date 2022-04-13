package app

import (
	"context"
	"os"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	"github.com/fildenisov/test-task-ticker-price/consts"
	"github.com/fildenisov/test-task-ticker-price/delivery/http"
	"github.com/fildenisov/test-task-ticker-price/domain/aggregator"
	"github.com/fildenisov/test-task-ticker-price/internal/rep"
	"github.com/fildenisov/test-task-ticker-price/models"
)

type Cmp struct {
	Name    string
	Service rep.Lifecycle
}

type App struct {
	log  *zerolog.Logger
	cfg  Config
	cmps []Cmp
}

func New(cfg Config) *App {
	l := zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Str("cmp", "app").Logger()
	return &App{
		log:  &l,
		cfg:  cfg,
		cmps: []Cmp{},
	}
}
func (a *App) Start(ctx context.Context) error {
	a.log.Info().Msg("starting application")

	agg := aggregator.New(a.cfg.Aggregator)

	h := http.New(a.cfg.HTTP)

	a.cmps = append(a.cmps, Cmp{"http", h}, Cmp{"aggregator", agg})

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
