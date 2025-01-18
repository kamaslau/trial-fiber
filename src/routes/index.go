package routes

import (
	"app/src/routes/post"

	"github.com/gofiber/fiber/v3"
)

// healthCheck For monitoring
func healthCheck(c fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "ok",
	})
}

func Root(c fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func InitRoutes(app *fiber.App) {
	app.Get("/", Root)
	app.Get("/health", healthCheck)

	// RESTful
	post.InitRoutes(app)

	// GraphQL
	// curl -X POST -H "Content-Type: application/json" --data "{\"query\":\"{}\",\"variables\":{}}" http://localhost:3000/graphql
	app.Post("/graphql", GraphQL)
}
