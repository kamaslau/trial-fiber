package main

import (
	"fmt"
	"log"
	"os"

	"app/src/internal/middlewares"
	"app/src/internal/routes"
	"app/src/internal/utils"
	"app/src/internal/utils/drivers"

	"github.com/gofiber/fiber/v3"
)

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
	// Setup
	utils.LoadEnv()
	drivers.InitDrivers()

	// Mount Middlewares and Route Handlers
	app := fiber.New()
	app.Use(middlewares.TryToken)
	routes.InitRoutes(app)
	startUp(app)
}
