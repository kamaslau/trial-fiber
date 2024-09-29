package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Fail loading .env: %s", err)
	}
}

func main() {
	// Load env variable(s)
	loadEnv()
	port := "3000"
	fmt.Printf("env.PORT: %s", os.Getenv("PORT"))
	if strings.Count(os.Getenv("PORT"), "") > 0 {
		port = os.Getenv("PORT")
	}

	app := fiber.New()

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("posts", func(c fiber.Ctx) error {
		response := fiber.Map{"title": "Post No.1", "content": "This would be an article."}
		return c.JSON(response)
	})

	err := app.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		fmt.Println("Error trying to launching fiber: ", err)
	}
}
