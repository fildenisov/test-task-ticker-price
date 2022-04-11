package http

import (
	"github.com/gofiber/fiber/v2"

	"price_aggregator/consts"
)

func (s *Server) setMiddlewares() {
	s.srv.Use(s.newLoggingMiddleware())
}

func (s *Server) newLoggingMiddleware() func(*fiber.Ctx) (err error) {
	return func(c *fiber.Ctx) (err error) {
		chainErr := c.Next()

		event := s.log.Info().
			Str(consts.FieldMethod, c.Method()).
			Str(consts.FieldURL, c.OriginalURL())

		if chainErr != nil {
			event = event.Err(chainErr)
		}

		event.Int("status", c.Response().StatusCode()).Msg("http request")
		return chainErr
	}
}
