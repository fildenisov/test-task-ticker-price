package http

func (s *Server) setRoutes() {
	v1 := s.srv.Group("/v1")
	v1.Get("/health", s.health)
}
