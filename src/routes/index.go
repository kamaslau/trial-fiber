package routes

import (
	"app/src/routes/post"

	"github.com/gofiber/fiber/v3"
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
