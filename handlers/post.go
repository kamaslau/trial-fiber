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
		return c.Status(404).JSON(ResNotFound)
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
		return c.Status(404).JSON(ResNotFound)
	}

	response := fiber.Map{"data": data}

	return c.JSON(response)
}

func Create(c fiber.Ctx) error {
	fmt.Println("Create: ")

	payload := new(models.Post)
	err := c.Bind().Body(payload)
	if err != nil {
		fmt.Println(err)
		return c.Status(400).JSON(fiber.Map{"succeed": "no", "message": "input error"})
	}
	payload.UUID = uuid.NewString()

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
	var id = c.Params("id")
	fmt.Printf("Update: id=%s", id)

	conditions := map[string]interface{}{"ID": id}

	// Lookup Target(s)
	var data []models.Post
	models.DBClient.Where(conditions).Find(&data) // No need to add 'deleted_at is null', GORM adds it by default with gorm.Model from type
	if len(data) == 0 {
		return c.Status(404).JSON(ResNotFound)
	}

	// TODO Do
	var payload = models.Post{}
	result := models.DBClient.Save(&payload)
	if result.RowsAffected != 1 {
		fmt.Println(result.Error)
		return c.Status(500).JSON(fiber.Map{"succeed": "no", "message": "Failed to update"})
	}

	return c.JSON(fiber.Map{"succeed": "yes"})
}

func Delete(c fiber.Ctx) error {
	var id = c.Params("id")
	fmt.Printf("Delete: id=%s", id)

	conditions := map[string]interface{}{"ID": id}

	// Lookup Target(s)
	var data []models.Post
	models.DBClient.Where(conditions).Find(&data) // No need to add 'deleted_at is null', GORM adds it by default with gorm.Model from type
	if len(data) == 0 {
		return c.Status(404).JSON(ResNotFound)
	}

	// Do
	result := models.DBClient.Where(conditions).Delete(&models.Post{})
	if result.RowsAffected != 1 {
		fmt.Println(result.Error)
		return c.Status(500).JSON(fiber.Map{"succeed": "no", "message": "Failed to delete"})
	}

	return c.JSON(fiber.Map{"succeed": "yes"})
}
