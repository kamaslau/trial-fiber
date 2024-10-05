package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/kamaslau/trial-fiber/handlers"
)

func InitPostRoutes(router fiber.Router) {
	route := router.Group("/post")

	route.Get("/", handlers.Find)
	route.Get("/:id", handlers.FindOne)
	route.Post("/", handlers.Create)
	route.Put("/:id", handlers.Update)
	route.Delete("/:id", handlers.Delete)
}
