package http

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

var start = time.Now()

// healthResp is a model for health response
type healthResp struct {
	Uptime    string
	Time      string
	Timestamp int32
}

// health is a handler for /v1/health path
func (s *Server) health(ctx *fiber.Ctx) error {
	resp := healthResp{
		Uptime:    time.Since(start).String(),
		Timestamp: int32(time.Now().Unix()),
		Time:      time.Now().Format(time.RFC3339Nano),
	}
	return sendResponse(ctx, resp, 200)
}
