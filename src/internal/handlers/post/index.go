package post

import (
	"errors"
	"log"
	"net/http"
	"reflect"

	"app/src/internal/handlers"
	"app/src/internal/middlewares"
	"app/src/internal/models"
	"app/src/internal/utils/drivers"
	"app/src/internal/utils/uuid"

	"github.com/gofiber/fiber/v3"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

const ModelName = "post"

func Count(c fiber.Ctx) error {
	log.Println("Count: ")

	acRequirement := middlewares.ACRequirement{
		Roles: []string{"adminer"},
		Permissions: []string{
			ModelName + "::read::all",
		},
	}
	if err := middlewares.TryAC(c, acRequirement); err != nil {
		log.Printf("🛑 TryAC failed: %v", err)
		return err
	}

	// Filter
	var filter = map[string]any{}
	var allowedFields = []string{"name"}
	if err := handlers.ComposeFilter(c, &filter, &allowedFields); err != nil {
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

	acRequirement := middlewares.ACRequirement{
		Roles: []string{"adminer"},
		Permissions: []string{
			ModelName + "::read::all",
		},
	}
	if err := middlewares.TryAC(c, acRequirement); err != nil {
		log.Printf("🛑 TryAC failed: %v", err)
		return err
	}

	// Filter
	var filter = map[string]any{}
	var allowedFields = []string{"name"}
	if err := handlers.ComposeFilter(c, &filter, &allowedFields); err != nil {
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
	}
	// log.Printf("sorter: %#v\n", sorter)

	// Pager
	// var pager = new(handlers.Pager)
	var pager = handlers.Pager{
		Limit:  10,
		Offset: 0,
	}
	if err := c.Bind().Query(&pager); err != nil {
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

	acRequirement := middlewares.ACRequirement{
		Roles: []string{"adminer"},
		Permissions: []string{
			ModelName + "::read::all",
		},
	}
	if err := middlewares.TryAC(c, acRequirement); err != nil {
		log.Printf("🛑 TryAC failed: %v", err)
		return err
	}

	var filter = map[string]any{"id": id}

	// Do Find
	var data models.Post
	if err := drivers.DBClient.Where(filter).First(&data).Error; err != nil {
		if err.Error() == "record not found" {
			return c.Status(http.StatusNotFound).JSON(handlers.GetHTTPMsg(http.StatusNotFound))
		}

		log.Printf("FindOne: database query failed: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(handlers.GetHTTPMsg(http.StatusInternalServerError))
	}

	return c.JSON(fiber.Map{"succeed": true, "data": data})
}

func CreateOne(c fiber.Ctx) error {
	acRequirement := middlewares.ACRequirement{
		Roles: []string{"adminer"},
		Permissions: []string{
			ModelName + "::write::all",
			ModelName + "::create::self",
		},
	}
	if err := middlewares.TryAC(c, acRequirement); err != nil {
		log.Printf("🛑 TryAC failed: %v", err)
		return err
	}

	// Parse payload
	var payload models.PostFieldsCreate
	if err := c.Bind().Body(&payload); err != nil {
		log.Printf("CreateOne: failed to parse request body: %v", err)
		return c.Status(http.StatusBadRequest).JSON(handlers.GetHTTPMsg(http.StatusBadRequest))
	}

	// Append metadata fields
	var dataset models.Post
	copier.Copy(&dataset, &payload)
	dataset.UUID = uuid.NewString()
	// log.Printf("CreateOne: dataset: %#v\n", &dataset)

	// Start transaction
	tx := drivers.DBClient.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Do Create
	if err := tx.Create(&dataset).Error; err != nil {
		tx.Rollback()
		log.Printf("CreateOne: database operation failed: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(handlers.GetHTTPMsg(http.StatusInternalServerError))
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("CreateOne: failed to commit transaction: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(handlers.GetHTTPMsg(http.StatusInternalServerError))
	}

	return c.JSON(fiber.Map{
		"succeed": true,
		"id":      dataset.ID,
	})
}

func UpdateOne(c fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("UpdateOne: id=%s\n", id)

	acRequirement := middlewares.ACRequirement{
		Roles: []string{"adminer"},
		Permissions: []string{
			ModelName + "::write::all",
			ModelName + "::update::self",
		},
	}
	if err := middlewares.TryAC(c, acRequirement); err != nil {
		log.Printf("🛑 TryAC failed: %v", err)
		return err
	}

	// Fetch Original
	var data models.User
	if err := drivers.DBClient.Where("id = ?", id).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(http.StatusNotFound).JSON(handlers.GetHTTPMsg(http.StatusNotFound))
		}
		return c.Status(http.StatusInternalServerError).JSON(handlers.GetHTTPMsg(http.StatusInternalServerError))
	}

	// Parse Target
	var payload models.User
	if err := c.Bind().Body(&payload); err != nil {
		return c.Status(http.StatusBadRequest).JSON(handlers.GetHTTPMsg(http.StatusBadRequest))
	}

	// Map and Merge
	{
		pv := reflect.ValueOf(&payload).Elem()
		dv := reflect.ValueOf(&data).Elem()
		typ := dv.Type()

		for i := 0; i < typ.NumField(); i++ {
			f := typ.Field(i)
			// Keep only necessary fields
			if (f.Anonymous && f.Type.Name() == "Model") ||
				f.Name == "UUID" {
				continue
			}
			src := pv.Field(i)
			dst := dv.Field(i)
			if dst.CanSet() && src.Type() == dst.Type() {
				dst.Set(src)
			}
		}
	}

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
	id := c.Params("id")
	log.Printf("DeleteOne: id=%s\n", id)

	acRequirement := middlewares.ACRequirement{
		Roles: []string{"adminer"},
		Permissions: []string{
			ModelName + "::delete::self",
		},
	}
	if err := middlewares.TryAC(c, acRequirement); err != nil {
		log.Printf("🛑 TryAC failed: %v", err)
		return err
	}

	filter := map[string]any{"id": id}

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
