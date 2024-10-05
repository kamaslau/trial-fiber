package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"

	// Local packages
	"github.com/kamaslau/trial-fiber/handlers"
	"github.com/kamaslau/trial-fiber/models"
)

var port = "3000"

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Fail loading .env: %s", err)
	}
}

func startUp(app *fiber.App) {
	fmt.Printf("env.PORT: %s", os.Getenv("PORT"))

	if strings.Count(os.Getenv("PORT"), "") > 0 {
		port = os.Getenv("PORT")
	}

	err := app.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		fmt.Println("Error trying to launching fiber: ", err)
	}
}

func main() {
	loadEnv() // Load env variable(s)

	app := fiber.New()

	models.ConnectDB()

	// Routers
	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// RESTful
	app.Get("/post/:id", handlers.Post)
	app.Get("/posts", handlers.Posts)

	// TODO GraphQL
	app.Get("/graphql", handlers.Posts) // TODO Placeholder

	startUp(app)
}
