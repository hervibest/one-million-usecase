package config

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func NewFiber() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      "upload",
		Prefork:      false,
		ErrorHandler: newErrorHandler(),
		BodyLimit:    100 * 1024 * 1024, // 50MB
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origins, Content-Type, Accept",
	}))

	return app
}

func newErrorHandler() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		code := http.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}

		response := &ErrorResponse{
			Success: false,
			Message: err.Error(),
		}

		return ctx.Status(code).JSON(response)

	}
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
