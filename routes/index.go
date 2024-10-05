package routes

import "github.com/gofiber/fiber/v3"

func Root(c fiber.Ctx) error {
	return c.SendString("Hello, World!")
}
