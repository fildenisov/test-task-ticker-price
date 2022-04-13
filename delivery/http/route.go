package http

import "fmt"

const (
	tickerParam = "ticker"
	limitParam  = "limit"
)

func (s *Server) setRoutes() {
	v1 := s.srv.Group("/v1")
	v1.Get("/health", s.health)
	v1.Get(fmt.Sprintf("/tickers/:%s/bars/:%s", tickerParam, limitParam), s.bars)
}
