package handlers

import "github.com/gofiber/fiber/v3"

func Posts(c fiber.Ctx) error {
	response := fiber.Map{"title": "Post No.1", "content": "This would be an article."}
	return c.JSON(response)
}
