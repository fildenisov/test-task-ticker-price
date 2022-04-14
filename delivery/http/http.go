package http

import (
	"context"
	"net/http"
	"os"
	"time"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/fildenisov/test-task-ticker-price/internal/rep"
)

// Server http
type Server struct {
	log *zerolog.Logger
	srv *fiber.App
	agg rep.Aggregator
	cfg Config
}

// New HTTP Server instance constructor
func New(cfg Config, agg rep.Aggregator) *Server {
	l := zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().
		Str("cmp", "http").Logger()
	srv := fiber.New(fiber.Config{
		WriteTimeout:             cfg.WriteTimeout,
		ReadTimeout:              cfg.ReadTimeout,
		DisableHeaderNormalizing: true,
	})

	return &Server{
		log: &l,
		cfg: cfg,
		srv: srv,
		agg: agg,
	}
}

// Start starts http server
func (s *Server) Start(ctx context.Context) error {
	s.setMiddlewares()
	s.setRoutes()

	errCh := make(chan error)
	s.log.Debug().Msgf("start listening %q", s.cfg.Address)
	go func() {
		if err := s.srv.Listen(s.cfg.Address); err != nil && err != http.ErrServerClosed {
			errCh <- errors.Wrap(err, "cannot listen and serve")
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-time.After(s.cfg.StartTimeout):
		return nil
	}
}

// Stop stops http server
func (s *Server) Stop(context.Context) error {
	errCh := make(chan error)
	s.log.Debug().Msgf("stopping %q", s.cfg.Address)
	go func() {
		if err := s.srv.Shutdown(); err != nil {
			errCh <- errors.Wrap(err, "cannot shutdown")
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-time.After(s.cfg.StopTimeout):
		return nil
	}
}
