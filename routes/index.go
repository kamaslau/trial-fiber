package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/kamaslau/trial-fiber/routes/post"
)

func Root(c fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func InitRoutes(app *fiber.App) {
	app.Get("/", Root)

	// RESTful
	post.InitRoutes(app)

	// TODO GraphQL
	app.Post("/graphql", GraphQL)
}
