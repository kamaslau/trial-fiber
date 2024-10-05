package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"

	"github.com/kamaslau/trial-fiber/models"
	"github.com/kamaslau/trial-fiber/routes"
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

	models.ConnectDB()

	app := fiber.New()

	routes.InitRoutes(app)

	startUp(app)
}
