package main

import (
	"fmt"
	"log"
	"os"

	"app/src/drivers"
	"app/src/routes"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/joho/godotenv"
)

// loadEnv Load env variable(s) from .env file
func loadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Fail loading .env: %s", err)
	}
}

// Initiate Service Components first, fail fast so we can break early
func initServices() {
	drivers.ConnectDB()    // Exit on failure
	drivers.ConnectCache() // Warning on failure
	drivers.ConnectMQ()    // Warning on failure
}

func startUp(app *fiber.App) {
	var port = "3000"

	if port = os.Getenv("PORT"); port != "" {
		log.Printf("👂 env.PORT: \033[33m%s\033[0m", os.Getenv("PORT"))
	}

	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal("Error trying to launching fiber: ", err)
	}
}

func main() {
	loadEnv()
	initServices()

	app := fiber.New()

	// Routes
	routes.InitRoutes(app)
	app.Use("/", static.New("./public")) // Serve static files from ./public

	startUp(app)
}
