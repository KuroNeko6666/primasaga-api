package router

import (
	"github.com/KuroNeko6666/prima-api/handler"
	"github.com/gofiber/fiber/v2"
)

func Post(app *fiber.App) {
	app.Get("/api/client/post", handler.GetPost)
	app.Get("/api/client/self/post", handler.GetSelfPost)         // temporary
	app.Get("/api/client/follower/post", handler.GetFollowerPost) // temporary
	app.Post("/api/client/post", handler.CreatePost)
	app.Delete("/api/client/post", handler.DeletePost)
	app.Put("/api/client/post", handler.UpdatePost)
	app.Post("/api/client/post/like", handler.PostLike)
	app.Post("/api/client/post/comment", handler.Comment)
	app.Delete("/api/client/post/comment", handler.DeleteComment)
	app.Post("/api/client/post/sub-comment", handler.SubComment)
}
