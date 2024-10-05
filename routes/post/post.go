package post

import (
	"github.com/gofiber/fiber/v3"
	"github.com/kamaslau/trial-fiber/handlers/post"
)

func InitRoutes(router fiber.Router) {
	route := router.Group("/post")

	route.Get("/", post.Find)
	route.Get("/:id", post.FindOne)
	route.Post("/", post.Create)
	route.Put("/:id", post.Update)
	route.Delete("/:id", post.Delete)
}
