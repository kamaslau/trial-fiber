package post

import (
	"fmt"
	"log"

	"app/src/drivers"
	"app/src/handlers"
	"app/src/models"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func Find(c fiber.Ctx) error {
	log.Println("Find: ")

	// TODO Extract inputs from JSON body
	type Pager struct {
		limit  int
		offset int
	}
	type Sorter map[string]interface{}
	type Filter map[string]interface{}
	type FindInput struct {
		Pager
		Sorter
		Filter
	}
	var payload FindInput
	c.Bind().Body(&payload)

	var count int64
	var data []models.Post
	var pager = map[string]int{"limit": 10, "offset": 0}
	var sorter Sorter = nil
	var filter Filter = nil // Empty filter

	// Do
	drivers.DBClient.Where(filter).Model(&models.Post{}).Count(&count)
	drivers.DBClient.Where(filter).Order(sorter).Limit(pager["limit"]).Offset(pager["offset"]).Find(&data)

	response := fiber.Map{
		"count":  count,
		"data":   data,
		"pager":  pager,
		"sorter": sorter,
		"filter": filter,
	}

	return c.JSON(response)
}

func FindOne(c fiber.Ctx) error {
	var id = c.Params("id")
	log.Printf("FindOne: id=%s\n", id)

	var data models.Post
	conditions := map[string]interface{}{"ID": id}
	drivers.DBClient.Where(conditions).First(&data)
	if data.ID == 0 {
		return c.Status(404).JSON(handlers.ResNotFound)
	}

	response := fiber.Map{"data": data}

	return c.JSON(response)
}

func Create(c fiber.Ctx) error {
	log.Println("Create: ")

	// Parse payload
	var payload models.Post
	err := c.Bind().Body(&payload)
	if err != nil {
		fmt.Println(err)
		return c.Status(400).JSON(fiber.Map{"succeed": "no", "message": "input error"})
	} else {
		fmt.Printf("payload: %#v\n", &payload)
	}
	payload.UUID = uuid.NewString()

	// Do
	result := drivers.DBClient.Create(&payload)
	if result.RowsAffected == 1 {
		response := fiber.Map{"succeed": "yes", "id": payload.ID}
		return c.JSON(response)
	} else {
		fmt.Println(result.Error)
		return c.Status(500).JSON(fiber.Map{"succeed": "no", "message": result.Error})
	}
}

func Update(c fiber.Ctx) error {
	var id = c.Params("id")
	log.Printf("Update: id=%s\n", id)

	conditions := map[string]interface{}{"ID": id}

	// Lookup Target
	var data models.Post
	drivers.DBClient.Where(conditions).First(&data)
	if data.ID == 0 {
		return c.Status(404).JSON(handlers.ResNotFound)
	} else {
		fmt.Printf("target: %#v\n", &data)
	}

	// Parse payload
	var payload models.Post
	err := c.Bind().Body(&payload)
	if err != nil {
		fmt.Println(err)
		return c.Status(400).JSON(fiber.Map{"succeed": "no", "message": "input error"})
	} else {
		fmt.Printf("payload: %#v\n", &payload)
	}

	// Merge payload to current data
	// TODO Optimize to map fields automatically
	data.Name = payload.Name
	data.Content = payload.Content
	data.Excerpt = payload.Excerpt

	// Do
	result := drivers.DBClient.Save(&data)
	if result.RowsAffected != 1 {
		fmt.Println(result.Error)
		return c.Status(500).JSON(fiber.Map{"succeed": "no", "message": "Failed to update"})
	}

	return c.JSON(fiber.Map{"succeed": "yes"})
}

func Delete(c fiber.Ctx) error {
	var id = c.Params("id")
	log.Printf("Delete: id=%s\n", id)

	conditions := map[string]interface{}{"ID": id}

	// Lookup Target(s)
	var data []models.Post
	drivers.DBClient.Where(conditions).Find(&data) // No need to add 'deleted_at is null', GORM adds it by default with gorm.Model from type
	if len(data) == 0 {
		return c.Status(404).JSON(handlers.ResNotFound)
	}

	// Do
	result := drivers.DBClient.Where(conditions).Delete(&models.Post{})
	if result.RowsAffected != 1 {
		fmt.Println(result.Error)
		return c.Status(500).JSON(fiber.Map{"succeed": "no", "message": "Failed to delete"})
	}

	return c.JSON(fiber.Map{"succeed": "yes"})
}
