package router

import (
	"github.com/KuroNeko6666/prima-api/handler"
	"github.com/gofiber/fiber/v2"
)

func Auth(app *fiber.App) {
	app.Post("/api/auth/login", handler.Login)
	app.Post("/api/auth/register", handler.Register)
	app.Delete("/api/auth/logout", handler.Logout)
	app.Get("/api/auth/verified", handler.EmailVerififation)
	app.Get("/api/auth/check", handler.Check)
}
