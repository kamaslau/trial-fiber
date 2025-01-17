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

	var filter = map[string]any{"ID": id}

	var data models.Post

	// Do Find
	if err := drivers.DBClient.Where(filter).First(&data).Error; err != nil {
		if err.Error() == "record not found" {
			return c.Status(http.StatusNotFound).JSON(handlers.GetHTTPMsg(http.StatusNotFound))
		}

		log.Printf("FindOne: database query failed: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(handlers.GetHTTPMsg(http.StatusInternalServerError))
	}

	return c.JSON(fiber.Map{"succeed": true, "data": data})
}

func Create(c fiber.Ctx) error {
	log.Println("Create: ")

	// Parse payload
	var payload models.Post
	if err := c.Bind().Body(&payload); err != nil {
		log.Printf("Create: failed to parse request body: %v", err)
		return c.Status(http.StatusBadRequest).JSON(handlers.GetHTTPMsg(http.StatusBadRequest))
	}

	// Append metadata fields
	payload.UUID = uuid.NewString()
	// log.Printf("payload: %#v\n", &payload)

	// Start transaction
	tx := drivers.DBClient.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Do Create
	if err := tx.Create(&payload).Error; err != nil {
		tx.Rollback()
		log.Printf("Create: database operation failed: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(handlers.GetHTTPMsg(http.StatusInternalServerError))
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Create: failed to commit transaction: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(handlers.GetHTTPMsg(http.StatusInternalServerError))
	}

	return c.JSON(fiber.Map{
		"succeed": true,
		"id":      payload.ID,
	})
}

func UpdateOne(c fiber.Ctx) error {
	var id = c.Params("id")
	if id == "" {
		log.Print("UpdateOne: failed to parse request id")
		return c.Status(http.StatusBadRequest).JSON(handlers.GetHTTPMsg(http.StatusBadRequest))
	}
	log.Printf("UpdateOne: id=%s\n", id)

	filter := map[string]any{"ID": id}

	// Lookup Target
	var data models.Post
	if err := drivers.DBClient.Where(filter).First(&data).Error; err != nil {
		if err.Error() == "record not found" {
			return c.Status(http.StatusNotFound).JSON(handlers.GetHTTPMsg(http.StatusNotFound))
		}

		log.Printf("UpdateOne: database query failed: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(handlers.GetHTTPMsg(http.StatusInternalServerError))
	}

	// Parse payload
	var payload models.Post
	if err := c.Bind().Body(&payload); err != nil {
		log.Printf("UpdateOne: failed to parse request body: %v", err)
		return c.Status(http.StatusBadRequest).JSON(handlers.GetHTTPMsg(http.StatusBadRequest))
	}
	// log.Printf("payload: %#v\n", &payload)

	// Merge payload to current data
	// TODO Optimize to map fields automatically
	data.Name = payload.Name
	data.Content = payload.Content
	data.Excerpt = payload.Excerpt

	// Start transaction
	tx := drivers.DBClient.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Do Update
	if err := tx.Save(&data).Error; err != nil {
		tx.Rollback()
		log.Printf("UpdateOne: failed to update record: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(handlers.GetHTTPMsg(http.StatusInternalServerError))
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("UpdateOne: failed to commit transaction: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(handlers.GetHTTPMsg(http.StatusInternalServerError))
	}

	return c.JSON(fiber.Map{"succeed": true})
}

func DeleteOne(c fiber.Ctx) error {
	var id = c.Params("id")
	if id == "" {
		log.Print("DeleteOne: failed to parse request id")
		return c.Status(http.StatusBadRequest).JSON(handlers.GetHTTPMsg(http.StatusBadRequest))
	}
	log.Printf("Delete: id=%s\n", id)

	filter := map[string]any{"ID": id}

	// Lookup Target
	var data models.Post
	if err := drivers.DBClient.Where(filter).First(&data).Error; err != nil {
		if err.Error() == "record not found" {
			return c.Status(http.StatusNotFound).JSON(handlers.GetHTTPMsg(http.StatusNotFound))
		}

		log.Printf("DeleteOne: database query failed: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(handlers.GetHTTPMsg(http.StatusInternalServerError))
	}

	// Start transaction
	tx := drivers.DBClient.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Do Delete
	if err := tx.Delete(&data).Error; err != nil {
		tx.Rollback()
		log.Printf("DeleteOne: failed to delete record: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(handlers.GetHTTPMsg(http.StatusInternalServerError))
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("DeleteOne: failed to commit transaction: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(handlers.GetHTTPMsg(http.StatusInternalServerError))
	}

	return c.JSON(fiber.Map{"succeed": true})
}
