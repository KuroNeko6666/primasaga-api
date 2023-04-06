package router

import (
	"github.com/KuroNeko6666/prima-api/handler"
	"github.com/gofiber/fiber/v2"
)

func Follow(app *fiber.App) {
	app.Post("/api/client/follow", handler.Follow)
	app.Get("/api/client/follow", handler.GetFollowers)
}
