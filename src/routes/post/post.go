package post

import (
	"app/src/handlers/post"

	"github.com/gofiber/fiber/v3"
)

func InitRoutes(router fiber.Router) {
	route := router.Group("/post")

	route.Get("/", post.Find)
	route.Get("/:id", post.FindOne)
	route.Post("/", post.Create)
	route.Put("/:id", post.Update)
	route.Delete("/:id", post.Delete)
}
