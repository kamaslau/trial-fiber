package post

import (
	handler "app/src/internal/handlers/post"

	"github.com/gofiber/fiber/v3"
)

const routeName = "post"

func InitRoutes(router fiber.Router) error {
	if router == nil {
		return fiber.NewError(fiber.StatusInternalServerError, "router instance is nil")
	}

	route := router.Group("/" + routeName)

	route.Get("/count", handler.Count)
	// curl "http://localhost:3000/routeName/count?filter=name:xxx"

	route.Get("/", handler.Find)
	// curl "http://localhost:3000/routeName?limit=5&offset=0&sorter=name+desc&filter=name:xxx"

	route.Get("/:id", handler.FindOne)
	// curl "http://localhost:3000/routeName/1"

	route.Post("/", handler.CreateOne)
	// curl -X POST -H "Content-Type: application/json" --data "{\"name\":\"Test Post Request Method Route\",\"content\":\"This is a placeholder.\"}" http://localhost:3000/routeName

	route.Put("/:id", handler.UpdateOne)
	// curl -X PUT -H "Content-Type: application/json" --data "{\"name\":\"Test Put Request Method Route\",\"content\":\"This is a placeholder.\"}" http://localhost:3000/routeName/1

	route.Delete("/:id", handler.DeleteOne)
	// curl -X DELETE http://localhost:3000/routeName/1

	return nil
}
