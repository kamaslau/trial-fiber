package post

import (
	"log"

	"app/src/drivers"
	"app/src/handlers"
	"app/src/models"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func Count(c fiber.Ctx) error {
	log.Println("Count: ")

	// Filter
	var filter = map[string]any{}
	if err := handlers.ComposeFilter(c, &filter); err != nil {
		log.Println(err)
		return c.Status(422).JSON(handlers.GetHTTPMsg(422))
	}

	var count int64
	drivers.DBClient.Where(filter).Model(&models.Post{}).Count(&count)

	return c.JSON(fiber.Map{"succeed": "yes", "count": count})
}

func Find(c fiber.Ctx) error {
	log.Println("Find: ")

	// Filter
	var filter = map[string]any{}
	if err := handlers.ComposeFilter(c, &filter); err != nil {
		log.Println(err)
		return c.Status(422).JSON(handlers.GetHTTPMsg(422))
	}

	// Do Count
	var count int64
	drivers.DBClient.Where(filter).Model(&models.Post{}).Count(&count)
	if count == 0 {
		return c.Status(404).JSON(handlers.GetHTTPMsg(404))
	}

	// Pager
	var pager = new(handlers.Pager)
	if err := c.Bind().Query(pager); err != nil {
		log.Println(err)
		return c.Status(422).JSON(handlers.GetHTTPMsg(422))
	} else {
		// log.Printf("pager: %#v\n", pager)
	}

	// Sorter
	var sorter = "id desc"
	if count > 1 {
		sorter = c.Query("sorter", "id desc")
		// log.Printf("sorter: %#v\n", sorter)
	}

	// Do Find
	var data []models.Post
	if count > int64(pager.Offset) {
		drivers.DBClient.Where(filter).Order(sorter).Limit(pager.Limit).Offset(pager.Offset).Find(&data)
	}

	// Output
	response := fiber.Map{
		"succeed": "yes",
		"count":   count,
		"data":    data,
		"metadata": fiber.Map{
			"filter": filter,
			"pager":  pager,
			"sorter": sorter,
		},
	}
	return c.JSON(response)
}

func FindOne(c fiber.Ctx) error {
	var id = c.Params("id")
	log.Printf("FindOne: id=%s\n", id)

	var data models.Post

	conditions := map[string]any{"ID": id}

	// Do Find
	drivers.DBClient.Where(conditions).First(&data)

	// Output
	if data.ID == 0 {
		return c.Status(404).JSON(handlers.GetHTTPMsg(404))
	} else {
		response := fiber.Map{"succeed": "yes", "data": data}
		return c.JSON(response)
	}
}

func Create(c fiber.Ctx) error {
	log.Println("Create: ")

	// Parse payload
	var payload models.Post
	if err := c.Bind().Body(&payload); err != nil {
		log.Println(err)
		return c.Status(400).JSON(handlers.GetHTTPMsg(400))
	} else {
		log.Printf("payload: %#v\n", &payload)
	}
	payload.UUID = uuid.NewString()

	// Do Create
	result := drivers.DBClient.Create(&payload)

	// Output
	if result.RowsAffected == 1 {
		response := fiber.Map{"succeed": "yes", "id": payload.ID}
		return c.JSON(response)
	} else {
		log.Println(result.Error)
		return c.Status(500).JSON(handlers.GetHTTPMsg(500))
	}
}

func UpdateOne(c fiber.Ctx) error {
	var id = c.Params("id")
	log.Printf("Update: id=%s\n", id)

	conditions := map[string]any{"ID": id}

	// Lookup Target
	var data models.Post
	drivers.DBClient.Where(conditions).First(&data)
	if data.ID == 0 {
		return c.Status(404).JSON(handlers.GetHTTPMsg(404))
	} else {
		log.Printf("target: %#v\n", &data)
	}

	// Parse payload
	var payload models.Post
	err := c.Bind().Body(&payload)
	if err != nil {
		log.Println(err)
		return c.Status(400).JSON(handlers.GetHTTPMsg(400))
	} else {
		log.Printf("payload: %#v\n", &payload)
	}

	// Merge payload to current data
	// TODO Optimize to map fields automatically
	data.Name = payload.Name
	data.Content = payload.Content
	data.Excerpt = payload.Excerpt

	// Do Update
	result := drivers.DBClient.Save(&data)
	log.Printf("result: %#v\n", &result)
	if result.RowsAffected != 1 {
		log.Println(result.Error)
		return c.Status(500).JSON(handlers.GetHTTPMsg(500))
	}

	return c.JSON(fiber.Map{"succeed": "yes"})
}

func DeleteOne(c fiber.Ctx) error {
	var id = c.Params("id")
	log.Printf("Delete: id=%s\n", id)

	conditions := map[string]any{"ID": id}

	// Lookup Target(s)
	var data []models.Post
	drivers.DBClient.Where(conditions).Find(&data) // No need to add 'deleted_at is null', GORM adds it by default with gorm.Model from type
	if len(data) == 0 {
		return c.Status(404).JSON(handlers.GetHTTPMsg(404))
	}

	// Do Delete
	result := drivers.DBClient.Where(conditions).Delete(&models.Post{})
	if result.RowsAffected != 1 {
		log.Println(result.Error)
		return c.Status(500).JSON(handlers.GetHTTPMsg(500))
	}

	return c.JSON(fiber.Map{"succeed": "yes"})
}
