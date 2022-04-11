package http

import (
	"github.com/gofiber/fiber/v2"
)

// sendResponse is a shortcut to write a response
func sendResponse(ctx *fiber.Ctx, resp interface{}, statusCode int) error {
	ctx.Status(statusCode)

	if ctx == nil {
		return nil
	}

	if err := ctx.JSON(resp); err != nil {
		return err
	}

	return nil
}
