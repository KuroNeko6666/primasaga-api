package router

import (
	"github.com/KuroNeko6666/prima-api/handler"
	"github.com/gofiber/fiber/v2"
)

func User(app *fiber.App) {
	app.Get("/api/client/user", handler.GetUsers)
	app.Get("/api/client/user/:id", handler.GetUser)
}
