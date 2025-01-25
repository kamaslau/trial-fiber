package routes

import (
	"log"
	"time"

	"app/src/routes/post"

	"github.com/gofiber/fiber/v3"
)

// healthCheck For monitoring
func healthCheck(c fiber.Ctx) error {
	now := time.Now()
	zone, offset := now.Zone()

	return c.JSON(fiber.Map{
		"status": "ok",
		"time": fiber.Map{
			"timezone": zone,
			"offset":   offset,
			"unix":     now.Unix(),
		},
	})
}

func Root(c fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func InitRoutes(app *fiber.App) {
	app.Get("/", Root)
	app.Get("/health", healthCheck)

	// RESTful
	restRoot := app.Group("/api/v1") // NOT Recommended: remove this line if you want to directly mount RESTful routes to the root path
	if err := post.InitRoutes(restRoot); err != nil {
		log.Panicf("Failed to initialize post routes: %v", err)
	}

	// GraphQL
	// curl -X POST -H "Content-Type: application/json" --data "{\"query\":\"{}\",\"variables\":{}}" http://localhost:3000/graphql
	app.Post("/graphql", GraphQL)
}
