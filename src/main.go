package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"app/src/drivers"
	"app/src/routes"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/joho/godotenv"
)

func loadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Fail loading .env: %s", err)
	}
}

func startUp(app *fiber.App) {
	var port = "3000"

	if strings.Count(os.Getenv("PORT"), "") > 0 {
		log.Printf("👂 env.PORT: \033[33m%s\033[0m", os.Getenv("PORT"))
		port = os.Getenv("PORT")
	}

	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal("Error trying to launching fiber: ", err)
	}
}

func initServices() {
	drivers.ConnectDB()    // Exit on failure
	drivers.ConnectCache() // Warning on failure
	drivers.ConnectMQ()    // Warning on failure
}

func main() {
	loadEnv() // Load env variable(s)

	// Initiate Service Components first, fail fast so we can break ea
	initServices()

	app := fiber.New()

	// Routes
	routes.InitRoutes(app)
	app.Use("/", static.New("./public"))

	startUp(app)
}
