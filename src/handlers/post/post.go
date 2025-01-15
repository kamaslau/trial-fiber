package post

import (
	"log"

	"app/src/drivers"
	"app/src/handlers"
	"app/src/models"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// Pager Paging segments
type Pager struct {
	limit  int `query:"limit"`
	offset int `query:"offset"`
}

func Count (c fiber.Ctx) error {
	log.Println("Count: ")

	var count int64
	drivers.DBClient.Model(&models.Post{}).Count(&count)

	return c.JSON(fiber.Map{"count": count})
}

func Find(c fiber.Ctx) error {
	log.Println("Find: ")

	var pager = new(Pager)
	if err := c.Bind().Query(pager); err != nil {
		log.Println(err)
		return c.Status(422).JSON(handlers.GetHTTPStatus(422))
	} else {
		log.Printf("pager: %#v\n", pager)
	}

	var sorter = map[string]string{"id": "desc"}
	var filter = map[string]any{}

	var count int64
	var data []models.Post

	// Do
	drivers.DBClient.Where(filter).Model(&models.Post{}).Count(&count)
	drivers.DBClient.Where(filter).Order(sorter).Limit(pager.limit).Offset(pager.offset).Find(&data)

	// Output
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

	conditions := map[string]any{"ID": id}

	// Do
	drivers.DBClient.Where(conditions).First(&data)

	// Output
	if data.ID == 0 {
		return c.Status(404).JSON(handlers.GetHTTPStatus(404))
	} else {
		response := fiber.Map{"data": data}
		return c.JSON(response)
	}
}

func Create(c fiber.Ctx) error {
	log.Println("Create: ")

	// Parse payload
	var payload models.Post
	if err := c.Bind().Body(&payload); err != nil {
		log.Println(err)
		return c.Status(400).JSON(handlers.GetHTTPStatus(400))
	} else {
		log.Printf("payload: %#v\n", &payload)
	}
	payload.UUID = uuid.NewString()

	// Do
	result := drivers.DBClient.Create(&payload)

	// Output
	if result.RowsAffected == 1 {
		response := fiber.Map{"succeed": "yes", "id": payload.ID}
		return c.JSON(response)
	} else {
		log.Println(result.Error)
		return c.Status(500).JSON(handlers.GetHTTPStatus(500))
	}
}

func Update(c fiber.Ctx) error {
	var id = c.Params("id")
	log.Printf("Update: id=%s\n", id)

	conditions := map[string]any{"ID": id}

	// Lookup Target
	var data models.Post
	drivers.DBClient.Where(conditions).First(&data)
	if data.ID == 0 {
		return c.Status(404).JSON(handlers.GetHTTPStatus(404))
	} else {
		log.Printf("target: %#v\n", &data)
	}

	// Parse payload
	var payload models.Post
	err := c.Bind().Body(&payload)
	if err != nil {
		log.Println(err)
		return c.Status(400).JSON(handlers.GetHTTPStatus(400))
	} else {
		log.Printf("payload: %#v\n", &payload)
	}

	// Merge payload to current data
	// TODO Optimize to map fields automatically
	data.Name = payload.Name
	data.Content = payload.Content
	data.Excerpt = payload.Excerpt

	// Do
	result := drivers.DBClient.Save(&data)
	if result.RowsAffected != 1 {
		log.Println(result.Error)
		return c.Status(500).JSON(handlers.GetHTTPStatus(500))
	}

	return c.JSON(fiber.Map{"succeed": "yes"})
}

func Delete(c fiber.Ctx) error {
	var id = c.Params("id")
	log.Printf("Delete: id=%s\n", id)

	conditions := map[string]any{"ID": id}

	// Lookup Target(s)
	var data []models.Post
	drivers.DBClient.Where(conditions).Find(&data) // No need to add 'deleted_at is null', GORM adds it by default with gorm.Model from type
	if len(data) == 0 {
		return c.Status(404).JSON(handlers.GetHTTPStatus(404))
	}

	// Do
	result := drivers.DBClient.Where(conditions).Delete(&models.Post{})
	if result.RowsAffected != 1 {
		log.Println(result.Error)
		return c.Status(500).JSON(handlers.GetHTTPStatus(500))
	}

	return c.JSON(fiber.Map{"succeed": "yes"})
}
