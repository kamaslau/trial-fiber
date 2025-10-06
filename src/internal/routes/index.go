package routes

import (
	"log"
	"os"
	"time"

	"app/src/internal/routes/post"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
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
	// RESTful
	// restRoot := app.Group("/aptx4869") // Recommended: mount routes after random characters, safer aginst anonymous sniffers
	// restRoot := app.Group("/") // NOT Recommended: mount routes directly to the root path, very exposed
	// restRoot := app.Group("/api/v1") // NOT Recommended: mount routes to versioned path, likely to be frequently sniffed
	restRootPath := "/"
	if path := os.Getenv("REST_ROOT"); path != "" {
		restRootPath = path
		log.Printf("👂 env.REST_ROOT: \033[33m%s\033[0m", os.Getenv("REST_ROOT"))
	} else {
		log.Print("⚠️ env.REST_ROOT not set, using '/' as RESTful root path, not safe against sniffers")
	}
	restRoot := app.Group(restRootPath)

	if err := post.InitRoutes(restRoot); err != nil {
		log.Panicf("Failed initializing post routes: %v", err)
	}
	// app.Get("/", Root)
	app.Use("/", static.New("./public"))
	app.Get("/health", healthCheck)

	// GraphQL
	// curl -X POST -H "Content-Type: application/json" --data "{\"query\":\"{}\",\"variables\":{}}" http://localhost:3000/graphql
	app.Post("/graphql", GraphQL)
}
