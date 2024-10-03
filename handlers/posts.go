package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
)

func Post(c fiber.Ctx) error {
	var id = c.Params("id")
	title := fmt.Sprintf("Post with id: %s", id)
	response := fiber.Map{"title": title, "content": "This should be one article."}
	return c.JSON(response)
}

func Posts(c fiber.Ctx) error {
	response := fiber.Map{"title": "Post No.1", "content": "This would be some articles."}
	return c.JSON(response)
}
