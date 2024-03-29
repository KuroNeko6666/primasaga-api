package app

import (
	"log"

	"github.com/KuroNeko6666/prima-api/config"
	"github.com/KuroNeko6666/prima-api/database"
	"github.com/KuroNeko6666/prima-api/middleware"
	"github.com/KuroNeko6666/prima-api/router"
	"github.com/gofiber/fiber/v2"
)

func StartApp() {
	//NEW APP
	app := fiber.New()

	//CONNECT DATABASE
	database.DBConnect()

	// ROUTES
	router.Auth(app)
	app.Use("/api/client", middleware.Auth)
	router.Follow(app)
	router.User(app)
	router.Post(app)

	//RUN APP
	log.Fatal(app.Listen(config.SERVER_PORT))
}
