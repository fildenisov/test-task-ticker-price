package http

import (
	"net/http"
	"strconv"

	"github.com/fildenisov/test-task-ticker-price/models"
	"github.com/gofiber/fiber/v2"
)

// barResp is a model for bar response
type barResp struct {
	Timestamp  int64   `json:"Timestamp"`
	IndexPrice float64 `json:"IndexPrice"`
}

// bars is a handler for /v1/tickers/{ticker}/bars/100
func (s *Server) bars(ctx *fiber.Ctx) error {
	t := ctx.Params(tickerParam)
	limitStr := ctx.Params(limitParam)

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		return sendResponse(ctx, "invalid limit", http.StatusBadRequest)
	}

	bs, ok := s.agg.GetBars(models.Ticker(t), int(limit))

	if !ok {
		return sendResponse(ctx, []barResp{}, http.StatusOK)
	}

	resp := make([]barResp, 0, len(bs))
	for _, b := range bs {
		resp = append(resp, barResp{
			Timestamp:  b.TS,
			IndexPrice: b.Price,
		})
	}

	return sendResponse(ctx, resp, http.StatusOK)
}
