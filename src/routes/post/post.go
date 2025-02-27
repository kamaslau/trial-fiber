package post

import (
	"app/src/handlers/post"

	"github.com/gofiber/fiber/v3"
)

func InitRoutes(router fiber.Router) error {
	if router == nil {
		return fiber.NewError(fiber.StatusInternalServerError, "router instance is nil")
	}

	route := router.Group("/post")

	route.Get("/count", post.Count)
	// curl "http://localhost:3000/api/v1/post/count?filter=name:xxx"

	route.Get("/", post.Find)
	// curl "http://localhost:3000/api/v1/post?limit=5&offset=0&sorter=name+desc&filter=name:xxx"

	route.Get("/:id", post.FindOne)
	// curl "http://localhost:3000/api/v1/post/1"

	route.Post("/", post.Create)
	// curl -X POST -H "Content-Type: application/json" --data "{\"name\":\"Test Post Request Method Route\",\"content\":\"This is a placeholder.\"}" http://localhost:3000/api/v1/post

	route.Put("/:id", post.UpdateOne)
	// curl -X PUT -H "Content-Type: application/json" --data "{\"name\":\"Test Put Request Method Route\",\"content\":\"This is a placeholder.\"}" http://localhost:3000/api/v1/post/1

	route.Delete("/:id", post.DeleteOne)
	// curl -X DELETE http://localhost:3000/api/v1/post/1

	return nil
}
