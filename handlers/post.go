package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/kamaslau/trial-fiber/models"
)

func Find(c fiber.Ctx) error {
	fmt.Println("Find: ")

	var data []models.Post
	models.DBClient.Find(&data)
	if len(data) == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No data"})
	}

	response := fiber.Map{"data": data}

	return c.JSON(response)
}

func FindOne(c fiber.Ctx) error {
	var id = c.Params("id")
	fmt.Printf("FindOne: id=%s", id)

	var data models.Post
	conditions := map[string]interface{}{"ID": id}
	models.DBClient.Where(conditions).First(&data)
	if data.ID == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No data"})
	}

	response := fiber.Map{"data": data}

	return c.JSON(response)
}

func Create(c fiber.Ctx) error {
	fmt.Println("Create: ")

	var payload = models.Post{
		UUID:    uuid.NewString(),
		Name:    "DB connected",
		Content: "This is an auto generated message on database connection succeed.",
	}
	result := models.DBClient.Create(&payload)
	if result.RowsAffected == 1 {
		response := fiber.Map{"succeed": "yes", "id": payload.ID}
		return c.JSON(response)
	} else {
		fmt.Println(result.Error)
		return c.Status(500).JSON(fiber.Map{"succeed": "no", "message": result.Error})
	}
}

// TODO
func Update(c fiber.Ctx) error {
	fmt.Println("Update: ")

	response := fiber.Map{"title": "Post No.1", "content": "This would be some articles."}
	return c.JSON(response)
}
