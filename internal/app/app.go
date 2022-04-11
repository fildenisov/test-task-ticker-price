package app

import (
	"context"
	"os"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	"price_aggregator/consts"
	"price_aggregator/delivery/http"
	"price_aggregator/internal/rep"
)

var (
	ErrStartTimeout    = errors.New("start timeout")
	ErrShutdownTimeout = errors.New("shutdown timeout")
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

	h, err := http.New(a.cfg.HTTP)
	if err != nil {
		a.log.Fatal().Err(err).Msg("cannot create http")
	}
	a.cmps = append(a.cmps, Cmp{"http", h})

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
		return ErrStartTimeout
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
		return ErrShutdownTimeout
	case err := <-errCh:
		if err != nil {
			return err
		}
		return nil
	}
}
