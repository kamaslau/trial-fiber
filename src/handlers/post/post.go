package post

import (
	"log"
	"net/http"

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
		log.Printf("Count: filter composition failed: %v", err)
		return c.Status(http.StatusUnprocessableEntity).JSON(handlers.GetHTTPMsg(http.StatusUnprocessableEntity))
	}
	// log.Printf("filter: %#v\n", filter)

	var count int64
	if err := drivers.DBClient.Where(filter).Model(&models.Post{}).Count(&count).Error; err != nil {
		log.Printf("Count: database query failed: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(handlers.GetHTTPMsg(http.StatusInternalServerError))
	}

	return c.JSON(fiber.Map{"succeed": true, "count": count})
}

func Find(c fiber.Ctx) error {
	log.Println("Find: ")

	// Filter
	var filter = map[string]any{}
	if err := handlers.ComposeFilter(c, &filter); err != nil {
		log.Printf("Count: filter composition failed: %v", err)
		return c.Status(http.StatusUnprocessableEntity).JSON(handlers.GetHTTPMsg(http.StatusUnprocessableEntity))
	}
	// log.Printf("filter: %#v\n", filter)

	// Query Instance
	var query = drivers.DBClient.Where(filter)

	// Do Count
	var count int64
	if err := query.Model(&models.Post{}).Count(&count).Error; err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(handlers.GetHTTPMsg(http.StatusInternalServerError))
	}
	if count == 0 {
		return c.Status(http.StatusNotFound).JSON(handlers.GetHTTPMsg(http.StatusNotFound))
	}

	// Sorter
	var sorter = "id desc"
	if count > 1 {
		sorter = c.Query("sorter", "id desc")
		// log.Printf("sorter: %#v\n", sorter)
	}

	// Pager
	var pager = new(handlers.Pager)
	if err := c.Bind().Query(pager); err != nil {
		log.Println(err)
		return c.Status(http.StatusUnprocessableEntity).JSON(handlers.GetHTTPMsg(http.StatusUnprocessableEntity))
	}
	// log.Printf("pager: %#v\n", pager)

	// Do Find
	var data []models.Post
	if count > int64(pager.Offset) {
		query.Order(sorter).Limit(pager.Limit).Offset(pager.Offset).Find(&data)

		if err := query.Order(sorter).Limit(pager.Limit).Offset(pager.Offset).Find(&data).Error; err != nil {
			log.Println(err)
			return c.Status(http.StatusInternalServerError).JSON(handlers.GetHTTPMsg(http.StatusInternalServerError))
		}
	}

	return c.JSON(fiber.Map{
		"succeed": true,
		"count":   count,
		"data":    data,
		"metadata": fiber.Map{
			"filter": filter,
			"pager":  pager,
			"sorter": sorter,
		},
	})
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
		return c.Status(http.StatusNotFound).JSON(handlers.GetHTTPMsg(http.StatusNotFound))
	} else {
		response := fiber.Map{"succeed": true, "data": data}
		return c.JSON(response)
	}
}

func Create(c fiber.Ctx) error {
	log.Println("Create: ")

	// Parse payload
	var payload models.Post
	if err := c.Bind().Body(&payload); err != nil {
		log.Println(err)
		return c.Status(http.StatusBadRequest).JSON(handlers.GetHTTPMsg(http.StatusBadRequest))
	} else {
		log.Printf("payload: %#v\n", &payload)
	}
	payload.UUID = uuid.NewString()

	// Do Create
	result := drivers.DBClient.Create(&payload)

	// Output
	if result.RowsAffected == 1 {
		response := fiber.Map{"succeed": true, "id": payload.ID}
		return c.JSON(response)
	} else {
		log.Println(result.Error)
		return c.Status(http.StatusInternalServerError).JSON(handlers.GetHTTPMsg(http.StatusInternalServerError))
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
		return c.Status(http.StatusNotFound).JSON(handlers.GetHTTPMsg(http.StatusNotFound))
	} else {
		log.Printf("target: %#v\n", &data)
	}

	// Parse payload
	var payload models.Post
	err := c.Bind().Body(&payload)
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusBadRequest).JSON(handlers.GetHTTPMsg(http.StatusBadRequest))
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
		return c.Status(http.StatusInternalServerError).JSON(handlers.GetHTTPMsg(http.StatusInternalServerError))
	}

	return c.JSON(fiber.Map{"succeed": true})
}

func DeleteOne(c fiber.Ctx) error {
	var id = c.Params("id")
	log.Printf("Delete: id=%s\n", id)

	conditions := map[string]any{"ID": id}

	// Lookup Target(s)
	var data []models.Post
	drivers.DBClient.Where(conditions).Find(&data) // No need to add 'deleted_at is null', GORM adds it by default with gorm.Model from type
	if len(data) == 0 {
		return c.Status(http.StatusNotFound).JSON(handlers.GetHTTPMsg(http.StatusNotFound))
	}

	// Do Delete
	result := drivers.DBClient.Where(conditions).Delete(&models.Post{})
	if result.RowsAffected != 1 {
		log.Println(result.Error)
		return c.Status(http.StatusInternalServerError).JSON(handlers.GetHTTPMsg(http.StatusInternalServerError))
	}

	return c.JSON(fiber.Map{"succeed": true})
}
