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

var port = "3000"

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Fail loading .env: %s", err)
	}
}

func startUp(app *fiber.App) {
	if strings.Count(os.Getenv("PORT"), "") > 0 {
		log.Printf("👂 env.PORT: \033[33m%s\033[0m", os.Getenv("PORT"))
		port = os.Getenv("PORT")
	}

	err := app.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal("Error trying to launching fiber: ", err)
	}
}

func main() {
	loadEnv() // Load env variable(s)

	// Initiate Components first, fail fast so we can debug faster
	drivers.ConnectDB()    // Exit on failure
	drivers.ConnectCache() // Warning on failure
	drivers.ConnectMQ()    // Warning on failure

	app := fiber.New()

	// Routes
	routes.InitRoutes(app)
	app.Use("/", static.New("./public"))

	startUp(app)
}
